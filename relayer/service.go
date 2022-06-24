package main

import (
	"context"
	"github.com/Ankr-network/bas-relay-hub/relayer/abigen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
	nc.chainId = big.NewInt(int64(chainConfig.ChainId))
	return nc
}

func (s *RelayService) epochWorker(source, target *nodeConfig) {
	log := log.WithField("source-name", source.chainConfig.ChainName)
	log.Infof("subscribing to the chain head events")
	blockNumberChannel := s.latestBlockFetcher(source.config, source.client)
	log.Infof("listening for incomming events")
	if source.chainConfig.EpochBlocks == 0 {
		log.Fatalf("zero epoch blocks is not possible")
	}
	latestTransitionedEpoch, err := target.relayHub.GetLatestTransitionedEpoch(&bind.CallOpts{}, source.chainId)
	if err != nil {
		log.Fatalf("failed to fetch latest transitioned epoch: %+v", err)
	}
	log.Infof("found latest transitioned epoch: %d", latestTransitionedEpoch)
	for {
		select {
		case number := <-blockNumberChannel:
			waitForBlock := (latestTransitionedEpoch + 1) * source.chainConfig.EpochBlocks
			log.WithField("block", number).WithField("nextEpochBlock", waitForBlock).WithField("diff", int64(waitForBlock)-int64(number)).Infof("latest block changed for chain")
			if number < waitForBlock {
				continue
			}
			log.Infof("epoch %d is reached, doing transition", latestTransitionedEpoch+1)
			if err := s.createEpochTransition(source, target, latestTransitionedEpoch*source.chainConfig.EpochBlocks); err != nil {
				log.Fatalf("failed to create epoch transition: %+v", err)
			}
			latestTransitionedEpoch++
		}
	}
}

func (s *RelayService) createEpochTransition(source, target *nodeConfig, epochBlock uint64) error {
	var blockProofs [][]byte
	for i := epochBlock; i <= epochBlock+source.config.ConfirmationBlocks; i++ {
		block, err := source.client.BlockByNumber(context.TODO(), big.NewInt(int64(i)))
		if err != nil {
			return err
		}
		blockRlp := parliaRlp(block.Header(), source.chainId)
		blockProofs = append(blockProofs, blockRlp)
	}
	_, _ = target.relayHub.UpdateValidatorSet(&bind.TransactOpts{}, source.chainId, blockProofs)
	return nil
}

func (s *RelayService) Stop() error {
	return nil
}
