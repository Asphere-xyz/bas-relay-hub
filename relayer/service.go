package main

import (
	"context"
	"github.com/Ankr-network/bas-relay-hub/relayer/abigen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	lru "github.com/hashicorp/golang-lru"
	"github.com/olebedev/emitter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/big"
	"sync"
	"time"
)

type RelayService struct {
	emitter *emitter.Emitter

	rootRpc, childRpc *ethclient.Client

	proofsLock sync.RWMutex
	proofsMap  *lru.Cache
}

func NewRelayService(emitter *emitter.Emitter) *RelayService {
	return &RelayService{
		emitter:   emitter,
		proofsMap: must(lru.New(1_000)),
	}
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
	// collect checkpoint proofs
	go s.collectCheckpointProofsWorker()
	return nil
}

func (s *RelayService) collectCheckpointProofsWorker() {
	event := s.emitter.On(checkpointProofReceived)
	for {
		select {
		case e := <-event:
			proof := e.Args[0].(*checkpointProof)
			s.mergeCheckpointProof(proof)
		}
	}
}

func (s *RelayService) findCheckpointProof(checkpointHash common.Hash) *checkpointProof {
	s.proofsLock.RLock()
	defer s.proofsLock.RUnlock()
	found, ok := s.proofsMap.Get(checkpointHash)
	if ok {
		return found.(*checkpointProof)
	}
	return nil
}

func (s *RelayService) mergeCheckpointProof(proof *checkpointProof) *checkpointProof {
	s.proofsLock.Lock()
	defer s.proofsLock.Unlock()
	existing, ok := s.proofsMap.Get(proof.checkpointHash)
	if !ok {
		s.proofsMap.Add(proof.epochNumber, proof)
		return proof
	}
	existingProof := existing.(*checkpointProof)
	for signer, signature := range proof.signatures {
		if existingProof.addSignature(signer, signature) {
			log.WithField("epoch", proof.epochNumber).WithField("signer", signer.Hex()).Infof("received signature for epoch")
		}
	}
	return existingProof
}

func (s *RelayService) latestBlockFetcher(nc *nodeConfig, client *ethclient.Client) chan uint64 {
	log.Infof("starting latest block fetcher")
	blockNumberChannel := make(chan uint64)
	refreshRate := time.Tick(3 * time.Second)
	var latestBlockNumber uint64
	go func() {
		for {
			<-refreshRate
			blockNumber, err := client.BlockNumber(context.TODO())
			if err != nil {
				log.WithError(err).Errorf("failed to fetch latest block number")
				continue
			}
			if blockNumber <= latestBlockNumber {
				continue
			}
			latestBlockNumber = blockNumber
			confirmedBlock := blockNumber - nc.chainConfig.Confirmations
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
		log.WithError(err).Fatalf("failed to fetch chain id from node")
	} else if big.NewInt(int64(chainConfig.ChainId)).Cmp(chainId) != 0 {
		log.Fatalf("chain id mismatched: %d != %d", chainConfig.ChainId, chainId.Uint64())
	}
	nc.chainId = chainId
	return nc
}

func (s *RelayService) epochWorker(source, target *nodeConfig) {
	log := log.WithField("source", source.chainConfig.ChainName).WithField("target", target.chainConfig.ChainName)
	log.Infof("subscribing to the chain head events")
	blockNumberChannel := s.latestBlockFetcher(source, source.client)
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
			// run checkpoint checks
			var success bool
			if len(source.config.Relayer.RelayerUrls) > 0 {
				const checkpointTransitionTimeout = 2 * time.Minute
				success = func() bool {
					ctx, cancel := context.WithTimeout(context.TODO(), checkpointTransitionTimeout)
					defer cancel()
					log.WithField("epoch", latestTransitionedEpoch+1).Infof("epoch is reached, doing checkpoint transition")
					err := s.createCheckpointTransition(ctx, source, target, waitForBlock, latestTransitionedEpoch)
					if err != nil {
						log.WithError(err).Warnf("checkpoint transition failed, trying confirmation transition")
						return false
					}
					return true
				}()
			}
			if !success {
				log.WithField("epoch", latestTransitionedEpoch+1).Infof("epoch is reached, doing confirmation transition")
				if err := s.createEpochTransition(context.TODO(), source, target, waitForBlock, latestKnownBlock); err != nil {
					log.WithError(err).WithField("epoch", latestTransitionedEpoch+1).Errorf("failed to create epoch transition")
					time.Sleep(30 * time.Second)
					break
				}
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

func (s *RelayService) createCheckpointTransition(ctx context.Context, source, target *nodeConfig, epochBlock, latestKnownBlock uint64) (err error) {
	// create and sign checkpoint proof
	proof, err := createCheckpointProof(ctx, source.client, epochBlock, source.chainConfig.EpochLength)
	if err != nil {
		return err
	}
	err = proof.signCheckpointProof(source.config.Relayer.PrivateKey)
	if err != nil {
		return err
	}
	proof = s.mergeCheckpointProof(proof)
	<-s.emitter.Emit(checkpointProofSigned, proof)
	// wait for quorum to be reached
	quorumRequired, err := calcRequiredQuorumForNextEpoch(ctx, source, epochBlock)
	if err != nil {
		return err
	}
	ticker := time.Tick(1 * time.Second)
	for len(proof.signatures) < quorumRequired {
		select {
		case <-ticker:
			break
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	// send transaction
	opts := injectSigner(target.chainId, source.config.Relayer.PrivateKey, nil)
	opts.Context = ctx
	tx, err := target.relayHub.CheckpointTransition(opts, source.chainId, proof.rawEpochBlock, proof.checkpointHash, mappingValues(proof.signatures))
	if err != nil {
		return err
	}
	log.WithFields(logrus.Fields{
		"chain":     target.chainConfig.ChainName,
		"hash":      tx.Hash().Hex(),
		"batchSize": 1,
		"gasUsed":   tx.Gas(),
	}).Infof("updated validator set")
	return waitTxToBeMined(ctx, target.client, tx.Hash())
}

func (s *RelayService) createEpochTransition(ctx context.Context, source, target *nodeConfig, epochBlock, latestKnownBlock uint64) (err error) {
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
	if target.chainConfig.GasLimit > 0 {
		var inputs [][]byte
		// estimate gas consumption
		blockProofs, err := createBlockProofs(ctx, source.client, epochBlock, source.chainConfig.EpochLength)
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
		batchSize = int(target.chainConfig.GasLimit-txGasReserve) / int(estimateGasUsage)
		inputSize := 0
		// create batch
		for bs := 0; bs < batchSize; bs++ {
			batchEpochBlock := epochBlock + uint64(bs)*source.chainConfig.EpochLength
			if batchEpochBlock > latestKnownBlock {
				break
			}
			blockProofs, err := createBlockProofs(ctx, source.client, batchEpochBlock, source.chainConfig.EpochLength)
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
		opts.GasLimit = 1_000_000
		blockProofs, err := createBlockProofs(ctx, source.client, epochBlock, source.chainConfig.EpochLength)
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
