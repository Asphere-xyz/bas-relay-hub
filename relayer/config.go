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
	Confirmations   uint64
	GasLimit        uint64
	EpochLength     uint64
}

type RelayerConfig struct {
	PrivateKey  *ecdsa.PrivateKey
	GrpcAddress string
	RelayerUrls []string
}

type Config struct {
	Root, Child ChainConfig

	Relayer RelayerConfig
}

func configFromViper(viper *viper.Viper) *Config {
	config := &Config{
		Root: ChainConfig{
			ChainName:       viper.GetString("root.chain-name"),
			RpcUrl:          viper.GetString("root.rpc-url"),
			RelayHubAddress: viper.GetString("root.relay-hub-address"),
			ChainId:         viper.GetUint64("root.chain-id"),
			Confirmations:   viper.GetUint64("root.confirmations"),
			GasLimit:        viper.GetUint64("root.gas-limit"),
			EpochLength:     viper.GetUint64("root.epoch-blocks"),
		},
		Child: ChainConfig{
			ChainName:       viper.GetString("child.chain-name"),
			RpcUrl:          viper.GetString("child.rpc-url"),
			RelayHubAddress: viper.GetString("child.relay-hub-address"),
			ChainId:         viper.GetUint64("child.chain-id"),
			Confirmations:   viper.GetUint64("child.confirmations"),
			GasLimit:        viper.GetUint64("child.gas-limit"),
			EpochLength:     viper.GetUint64("child.epoch-blocks"),
		},
		Relayer: RelayerConfig{
			PrivateKey:  must(crypto.HexToECDSA(viper.GetString("relayer.private-key"))),
			GrpcAddress: viper.GetString("relayer.grpc-address"),
			RelayerUrls: viper.GetStringSlice("relayer.relayer-urls"),
		},
	}
	return config
}
