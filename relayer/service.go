package main

import (
	"context"
	"fmt"
	"github.com/Ankr-network/bas-relay-hub/relayer/abigen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/big"
	"time"
)

type RelayService struct {
	rootRpc, childRpc *ethclient.Client
}

func NewRelayService() *RelayService {
	return &RelayService{}
}

func (s *RelayService) Start(config *Config) (err error) {
	// establish connection with root and child RPCs
	log.Infof("connecting to the root RPC: %s", config.Root.RpcUrl)
	s.rootRpc, err = ethclient.Dial(config.Root.RpcUrl)
	if err != nil {
		return err
	}
	log.Infof("connecting to the child RPC: %s", config.Child.RpcUrl)
	s.childRpc, err = ethclient.Dial(config.Child.RpcUrl)
	if err != nil {
		return err
	}
	// start background worker
	log.Infof("starting background worker")
	go func() {
		if err := s.backgroundWorker(config); err != nil {
			log.Fatalf("background worker failed: %+v", err)
		}
	}()
	return nil
}

func (s *RelayService) latestBlockFetcher(config *Config, client *ethclient.Client) chan uint64 {
	log.Infof("starting latest block fetcher")
	blockNumberChannel := make(chan uint64)
	refreshRate := time.Tick(3 * time.Second)
	var latestBlockNumber uint64
	go func() {
		for {
			<-refreshRate
			blockNumber, err := client.BlockNumber(context.TODO())
			if err != nil {
				log.Error("failed to fetch latest block number: %+v", err)
				continue
			}
			if blockNumber <= latestBlockNumber {
				continue
			}
			latestBlockNumber = blockNumber
			confirmedBlock := blockNumber - config.ConfirmationBlocks
			if confirmedBlock < 0 {
				continue
			}
			blockNumberChannel <- confirmedBlock
		}
	}()
	return blockNumberChannel
}

func (s *RelayService) backgroundWorker(config *Config) error {
	// create stream from root to child
	go s.epochWorker(s.toNodeConfig(config, true), s.toNodeConfig(config, false))
	go s.epochWorker(s.toNodeConfig(config, false), s.toNodeConfig(config, true))
	return nil
}

type nodeConfig struct {
	client      *ethclient.Client
	relayHub    *abigen.RelayHub
	config      *Config
	chainConfig ChainConfig
	chainId     *big.Int
}

func (s *RelayService) toNodeConfig(config *Config, isRoot bool) *nodeConfig {
	var chainConfig ChainConfig
	nc := &nodeConfig{}
	if isRoot {
		nc.client = s.rootRpc
		chainConfig = config.Root
	} else {
		nc.client = s.childRpc
		chainConfig = config.Child
	}
	nc.relayHub, _ = abigen.NewRelayHub(common.HexToAddress(chainConfig.RelayHubAddress), nc.client)
	nc.config = config
	nc.chainConfig = chainConfig
	// make sure chain id is correct
	chainId, err := nc.client.ChainID(context.TODO())
	if err != nil {
		log.WithField("err", err).Fatalf("failed to fetch chain id from node")
	} else if big.NewInt(int64(chainConfig.ChainId)).Cmp(chainId) != 0 {
		log.Fatalf("chain id mismatched: %d != %d", chainConfig.ChainId, chainId.Uint64())
	}
	nc.chainId = chainId
	return nc
}

func (s *RelayService) epochWorker(source, target *nodeConfig) {
	log := log.WithField("chain", source.chainConfig.ChainName)
	log.Infof("subscribing to the chain head events")
	blockNumberChannel := s.latestBlockFetcher(source.config, source.client)
	log.Infof("listening for incomming events")
	if source.chainConfig.EpochLength == 0 {
		log.Fatalf("zero epoch blocks is not possible")
	}
	latestTransitionedEpoch, err := target.relayHub.GetLatestTransitionedEpoch(&bind.CallOpts{}, source.chainId)
	if err != nil {
		log.Fatalf("failed to fetch latest transitioned epoch: %+v", err)
	}
	lastProcessTime := time.Now().Unix()
	var totalProcessedBlocks uint64
	log.WithField("epoch", latestTransitionedEpoch).Infof("found latest transitioned epoch")
	for {
		select {
		case latestKnownBlock := <-blockNumberChannel:
			waitForBlock := (latestTransitionedEpoch + 1) * source.chainConfig.EpochLength
			log.WithField("block", latestKnownBlock).WithField("nextEpochBlock", waitForBlock).WithField("diff", int64(waitForBlock)-int64(latestKnownBlock)).Debug("latest block changed for chain")
			if latestKnownBlock < waitForBlock || waitForBlock > latestKnownBlock {
				continue
			}
			log.WithField("epoch", latestTransitionedEpoch+1).Infof("epoch is reached, doing transition")
			if err := s.createEpochTransition(context.TODO(), source, target, waitForBlock, latestKnownBlock, int64(source.config.GasLimit)); err != nil {
				log.WithField("err", err).Errorf("failed to create epoch transition")
				time.Sleep(30 * time.Second)
				break
			}
			prevLatestTransitionedEpoch := latestTransitionedEpoch
			latestTransitionedEpoch, err = target.relayHub.GetLatestTransitionedEpoch(&bind.CallOpts{}, source.chainId)
			if err != nil {
				log.Fatalf("failed to fetch latest transitioned epoch: %+v", err)
			}
			processedBlocks := (latestTransitionedEpoch - prevLatestTransitionedEpoch) * source.chainConfig.EpochLength
			totalProcessedBlocks += processedBlocks
			blocksPerSecond := int64(totalProcessedBlocks) / (time.Now().Unix() - lastProcessTime)
			if blocksPerSecond == 0 {
				continue
			}
			estimateTime := int64(latestKnownBlock-latestTransitionedEpoch*source.chainConfig.EpochLength) / blocksPerSecond
			log.WithField("bps", blocksPerSecond).WithField("eta", prettyFormatTime(estimateTime)).Info("chain epochs synchronization stats")
		}
	}
}

func (s *RelayService) createBlockProofs(ctx context.Context, client *ethclient.Client, atBlock, epochLength uint64) ([][]byte, error) {
	prevEpochBlock, err := client.BlockByNumber(ctx, big.NewInt(int64(atBlock-epochLength)))
	if err != nil {
		return nil, errors.Wrapf(err, "can't fetch prev epoch block (%d)", atBlock-epochLength)
	}
	prevEpochValidators, err := extractParliaValidators(prevEpochBlock.Header())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to extract parlia block validators")
	}
	uniqueSigners := make(map[common.Address]bool)
	quorumRequired := len(prevEpochValidators) * 2 / 3
	var blockProofs [][]byte
	for i := atBlock; i < atBlock+epochLength; i++ {
		block, err := client.BlockByNumber(ctx, big.NewInt(int64(i)))
		if err != nil {
			return nil, errors.Wrapf(err, "can't fetch fetch block (%d)", i)
		}
		blockRlp, err := rlp.EncodeToBytes(block.Header())
		if err != nil {
			return nil, errors.Wrapf(err, "rlp encode failed")
		}
		blockProofs = append(blockProofs, blockRlp)
		// make sure quorum is reached
		uniqueSigners[block.Header().Coinbase] = true
		if len(uniqueSigners) >= quorumRequired {
			break
		}
	}
	if len(uniqueSigners) < quorumRequired {
		return nil, fmt.Errorf("failed to reach quorum for epoch block %d, reached only %d/%d", atBlock, len(uniqueSigners), quorumRequired)
	}
	return blockProofs, nil
}

func (s *RelayService) createEpochTransition(ctx context.Context, source, target *nodeConfig, epochBlock, latestKnownBlock uint64, gasLimit int64) (err error) {
	opts := injectSigner(target.chainId, source.config.Relayer.PrivateKey, nil)
	opts.Context = ctx
	log.WithFields(logrus.Fields{
		"source":     source.chainConfig.ChainName,
		"target":     target.chainConfig.ChainName,
		"epochBlock": epochBlock,
		"from":       opts.From.Hex(),
	}).Infof("executing epoch transition")
	var tx *types.Transaction
	batchSize := 1
	if gasLimit > 0 {
		var inputs [][]byte
		// estimate gas consumption
		blockProofs, err := s.createBlockProofs(ctx, source.client, epochBlock, source.chainConfig.EpochLength)
		if err != nil {
			return err
		}
		opts.NoSend = true
		tx, err = target.relayHub.UpdateValidatorSet(opts, source.chainId, blockProofs)
		if err != nil {
			return errors.Wrapf(err, "failed to update validator set")
		}
		opts.NoSend = false
		estimateGasUsage := tx.Gas()
		const txGasReserve = 100_000
		const txMaxInput = 100 * 1024
		batchSize = int(gasLimit-txGasReserve) / int(estimateGasUsage)
		inputSize := 0
		// create batch
		for bs := 0; bs < batchSize; bs++ {
			batchEpochBlock := epochBlock + uint64(bs)*source.chainConfig.EpochLength
			if batchEpochBlock > latestKnownBlock {
				break
			}
			blockProofs, err := s.createBlockProofs(ctx, source.client, batchEpochBlock, source.chainConfig.EpochLength)
			if err != nil {
				return err
			}
			input := encodeFunctionCall("updateValidatorSet(uint256,bytes[])", source.chainId, blockProofs)
			// transaction has max size
			if len(input)+inputSize > txMaxInput {
				batchSize = bs
				break
			}
			inputSize += len(input)
			inputs = append(inputs, input)
		}
		tx, err = target.relayHub.Multicall(opts, inputs)
		if err != nil {
			return errors.Wrapf(err, "failed to do multicall")
		}
	} else {
		blockProofs, err := s.createBlockProofs(ctx, source.client, epochBlock, source.chainConfig.EpochLength)
		if err != nil {
			return err
		}
		tx, err = target.relayHub.UpdateValidatorSet(opts, source.chainId, blockProofs)
		if err != nil {
			return errors.Wrapf(err, "failed to update validator set")
		}
	}
	log.WithFields(logrus.Fields{
		"chain":     target.chainConfig.ChainName,
		"hash":      tx.Hash().Hex(),
		"batchSize": batchSize,
		"gasUsed":   tx.Gas(),
	}).Infof("updated validator set")
	return waitTxToBeMined(ctx, target.client, tx.Hash())
}

func (s *RelayService) Stop() error {
	return nil
}
