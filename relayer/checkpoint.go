package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/pkg/errors"
	"math/big"
)

var checkpointProofFields = mustNewArguments("bytes32", "bytes32")

type checkpointProof struct {
	newValidators  []common.Address
	checkpointHash []byte
}

func createCheckpointProof(ctx context.Context, client *ethclient.Client, epochBlock, epochLength uint64) (*checkpointProof, error) {
	if epochBlock%epochLength != 0 {
		return nil, fmt.Errorf("bad epoch block (%d)", epochBlock)
	}
	result := &checkpointProof{}
	checkpointTrie := trie.NewEmpty(trie.NewDatabase(memorydb.New()))
	for i := epochBlock; i < epochBlock+epochLength; i++ {
		block, err := client.BlockByNumber(ctx, big.NewInt(int64(i)))
		if err != nil {
			return nil, errors.Wrapf(err, "can't fetch fetch block (%d)", i)
		}
		if i == epochBlock {
			prevEpochValidators, err := extractParliaValidators(block.Header())
			if err != nil {
				return nil, errors.Wrapf(err, "can't extract validators from block (%d)", i)
			}
			result.newValidators = prevEpochValidators
		}
		blockProof, err := checkpointProofFields.Pack(block.Hash(), block.ReceiptHash())
		if err != nil {
			return nil, errors.Wrapf(err, "can't create proof for block (%d)", i)
		}
		proofPath, err := rlp.EncodeToBytes(i)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to encode block (%d) path", i)
		}
		checkpointTrie.Update(proofPath, blockProof)
	}
	result.checkpointHash = checkpointTrie.Hash().Bytes()
	return result, nil
}
