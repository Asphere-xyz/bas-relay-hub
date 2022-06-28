package main

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func injectSigner(chainId *big.Int, privateKey *ecdsa.PrivateKey, opts *bind.TransactOpts) *bind.TransactOpts {
	if opts == nil {
		opts = &bind.TransactOpts{}
	}
	opts.Signer = func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return types.SignTx(tx, types.NewLondonSigner(chainId), privateKey)
	}
	opts.From = crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey))
	return opts
}
