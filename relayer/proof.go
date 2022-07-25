package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"math/big"
)

func calcRequiredQuorumForNextEpoch(ctx context.Context, nc *nodeConfig, epochBlock uint64) (int, error) {
	prevEpochBlock, err := nc.client.BlockByNumber(ctx, big.NewInt(int64(epochBlock-nc.chainConfig.EpochLength)))
	if err != nil {
		return 0, errors.Wrapf(err, "can't fetch prev epoch block (%d)", epochBlock-nc.chainConfig.EpochLength)
	}
	prevEpochValidators, _, err := extractParliaValidators(prevEpochBlock.Header())
	if err != nil {
		return 0, errors.Wrapf(err, "failed to extract parlia block validators")
	}
	return len(prevEpochValidators) * 2 / 3, nil
}

func createBlockProofs(ctx context.Context, client *ethclient.Client, epochBlock, epochLength uint64) ([][]byte, error) {
	chainId, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	prevEpochBlock, err := client.BlockByNumber(ctx, big.NewInt(int64(epochBlock-epochLength)))
	if err != nil {
		return nil, errors.Wrapf(err, "can't fetch prev epoch block (%d)", epochBlock-epochLength)
	}
	prevEpochValidators, validatorsMap, err := extractParliaValidators(prevEpochBlock.Header())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to extract parlia block validators")
	}
	uniqueSigners := make(map[common.Address]bool)
	quorumRequired := len(prevEpochValidators) * 2 / 3
	var blockProofs [][]byte
	for i := epochBlock; i < epochBlock+epochLength; i++ {
		block, err := client.BlockByNumber(ctx, big.NewInt(int64(i)))
		if err != nil {
			return nil, errors.Wrapf(err, "can't fetch fetch block (%d)", i)
		}
		// check block signer
		coinbase := block.Coinbase()
		signer, err := recoverParliaBlockSigner(block.Header(), chainId)
		if err != nil {
			return nil, err
		} else if coinbase != signer {
			return nil, fmt.Errorf("recovered bad block signer (coinbase != signer): %s != %s", coinbase.Hex(), signer.Hex())
		}
		if !validatorsMap[signer] {
			log.Warnf("block's (%d) coinbase not in validator set (%s)", block.NumberU64(), signer.Hex())
		}
		// append new proof
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
		return nil, fmt.Errorf("failed to reach quorum for epoch block %d, reached only %d/%d", epochBlock, len(uniqueSigners), quorumRequired)
	}
	return blockProofs, nil
}
