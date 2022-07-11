package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/pkg/errors"
	"math/big"
	"sync"
)

type checkpointProof struct {
	epochNumber    uint64
	rawEpochBlock  []byte
	blockHash      common.Hash
	checkpointHash common.Hash
	signatures     map[common.Address][]byte
	// state
	rwLock sync.RWMutex
}

func createCheckpointProof(ctx context.Context, client *ethclient.Client, epochBlock, epochLength uint64) (*checkpointProof, error) {
	if epochBlock%epochLength != 0 {
		return nil, fmt.Errorf("bad epoch block (%d)", epochBlock)
	}
	result := &checkpointProof{
		epochNumber: epochBlock / epochLength,
	}
	checkpointTrie := trie.NewEmpty(trie.NewDatabase(memorydb.New()))
	for i := epochBlock; i < epochBlock+epochLength; i++ {
		block, err := client.BlockByNumber(ctx, big.NewInt(int64(i)))
		if err != nil {
			return nil, errors.Wrapf(err, "can't fetch fetch block (%d)", i)
		}
		if i == epochBlock {
			result.rawEpochBlock, err = rlp.EncodeToBytes(block.Header())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to encode epoch block (%d)", block.NumberU64())
			}
			result.blockHash = block.Hash()
		}
		proofPath, err := rlp.EncodeToBytes(i)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to encode block (%d) path", i)
		}
		checkpointTrie.Update(proofPath, block.Hash().Bytes())
	}
	result.checkpointHash = checkpointTrie.Hash()
	return result, nil
}

func (proof *checkpointProof) getSignatureOrNil(signer common.Address) []byte {
	proof.rwLock.RLock()
	defer proof.rwLock.RUnlock()
	if proof.signatures == nil {
		return nil
	}
	return proof.signatures[signer]
}

func (proof *checkpointProof) signCheckpointProof(privateKey *ecdsa.PrivateKey) error {
	proof.rwLock.Lock()
	defer proof.rwLock.Unlock()
	signature, err := crypto.Sign(proof.checkpointHash.Bytes(), privateKey)
	if err != nil {
		return err
	} else if proof.signatures == nil {
		proof.signatures = make(map[common.Address][]byte)
	}
	signer := crypto.PubkeyToAddress(privateKey.PublicKey)
	if _, ok := proof.signatures[signer]; ok {
		log.Warnf("proof for epoch block (%s) is already signed by (%s)", proof.blockHash.Hex(), signer.Hex())
		return nil
	}
	proof.signatures[signer] = signature
	return nil
}

func (proof *checkpointProof) hasSignature(signer common.Address) bool {
	proof.rwLock.RLock()
	defer proof.rwLock.RUnlock()
	if proof.signatures == nil {
		return false
	}
	_, ok := proof.signatures[signer]
	return ok
}

func (proof *checkpointProof) addSignature(signer common.Address, signature []byte) bool {
	proof.rwLock.Lock()
	defer proof.rwLock.Unlock()
	if proof.signatures == nil {
		proof.signatures = make(map[common.Address][]byte)
	}
	_, exists := proof.signatures[signer]
	proof.signatures[signer] = signature
	return !exists
}
