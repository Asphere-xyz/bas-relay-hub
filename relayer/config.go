package main

import (
	"github.com/spf13/viper"
)

type ChainConfig struct {
	ChainName       string
	RpcUrl          string
	RelayHubAddress string
	ChainId         uint64
	EpochBlocks     uint64
}

type Config struct {
	Root, Child ChainConfig

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
		ConfirmationBlocks: viper.GetUint64("confirmation-blocks"),
	}
	return config
}
