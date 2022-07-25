package main

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

const (
	extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
)

var errBadParliaBlock = fmt.Errorf("bad parlia block")

func parliaRlp(header *types.Header, chainId *big.Int) []byte {
	b := new(bytes.Buffer)
	encodeSigHeader(b, header, chainId)
	return b.Bytes()
}

func encodeSigHeader(w io.Writer, header *types.Header, chainId *big.Int) {
	err := rlp.Encode(w, []interface{}{
		chainId,
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // this will panic if extra is too short, should check before calling encodeSigHeader
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}

func recoverParliaBlockSigner(header *types.Header, chainId *big.Int) (signer common.Address, err error) {
	if len(header.Extra) < extraSeal {
		return signer, errBadParliaBlock
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]
	b := new(bytes.Buffer)
	err = rlp.Encode(b, []interface{}{
		chainId,
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // this will panic if extra is too short, should check before calling encodeSigHeader
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
	signingData := b.Bytes()
	publicKey, err := crypto.Ecrecover(crypto.Keccak256(signingData), signature)
	if err != nil {
		return signer, err
	}
	copy(signer[:], crypto.Keccak256(publicKey[1:])[12:])
	return signer, nil
}

func extractParliaValidators(header *types.Header) ([]common.Address, map[common.Address]bool, error) {
	validatorBytes := header.Extra[extraVanity : len(header.Extra)-extraSeal]
	if len(validatorBytes)%common.AddressLength != 0 {
		return nil, nil, errBadParliaBlock
	}
	n := len(validatorBytes) / common.AddressLength
	result := make([]common.Address, n)
	mapping := make(map[common.Address]bool)
	for i := 0; i < n; i++ {
		address := make([]byte, common.AddressLength)
		copy(address, validatorBytes[i*common.AddressLength:(i+1)*common.AddressLength])
		result[i] = common.BytesToAddress(address)
		mapping[common.BytesToAddress(address)] = true
	}
	return result, mapping, nil
}
