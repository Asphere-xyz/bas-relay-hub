package main

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

type ChainConfig struct {
	ChainName       string
	RpcUrl          string
	RelayHubAddress string
	ChainId         uint64
	EpochBlocks     uint64
}

type RelayerConfig struct {
	PrivateKey *ecdsa.PrivateKey
}

type Config struct {
	Root, Child ChainConfig

	Relayer RelayerConfig

	ConfirmationBlocks uint64
}

func ConfigFromViper(viper *viper.Viper) *Config {
	config := &Config{
		Root: ChainConfig{
			ChainName:       viper.GetString("root.chain-name"),
			RpcUrl:          viper.GetString("root.rpc-url"),
			RelayHubAddress: viper.GetString("root.relay-hub-address"),
			ChainId:         viper.GetUint64("root.chain-id"),
			EpochBlocks:     viper.GetUint64("root.epoch-blocks"),
		},
		Child: ChainConfig{
			ChainName:       viper.GetString("child.chain-name"),
			RpcUrl:          viper.GetString("child.rpc-url"),
			RelayHubAddress: viper.GetString("child.relay-hub-address"),
			ChainId:         viper.GetUint64("child.chain-id"),
			EpochBlocks:     viper.GetUint64("child.epoch-blocks"),
		},
		Relayer: RelayerConfig{
			PrivateKey: must(crypto.HexToECDSA(viper.GetString("relayer.private-key"))),
		},
		ConfirmationBlocks: viper.GetUint64("confirmation-blocks"),
	}
	return config
}
