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
	prevEpochValidators, err := extractParliaValidators(prevEpochBlock.Header())
	if err != nil {
		return 0, errors.Wrapf(err, "failed to extract parlia block validators")
	}
	return len(prevEpochValidators) * 2 / 3, nil
}

func createBlockProofs(ctx context.Context, client *ethclient.Client, epochBlock, epochLength uint64) ([][]byte, error) {
	prevEpochBlock, err := client.BlockByNumber(ctx, big.NewInt(int64(epochBlock-epochLength)))
	if err != nil {
		return nil, errors.Wrapf(err, "can't fetch prev epoch block (%d)", epochBlock-epochLength)
	}
	prevEpochValidators, err := extractParliaValidators(prevEpochBlock.Header())
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
