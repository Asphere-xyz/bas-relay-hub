package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"strings"
	"time"
)

func must[T interface{}](value *T, err error) *T {
	if err != nil {
		panic(err)
	}
	return value
}

func mappingKeys[K comparable, V any](value map[K]V) (result []K) {
	for k := range value {
		result = append(result, k)
	}
	return result
}

func mappingValues[K comparable, V any](value map[K]V) (result []V) {
	for _, v := range value {
		result = append(result, v)
	}
	return result
}

func waitTxToBeMined(ctx context.Context, eth *ethclient.Client, txHash common.Hash) error {
	tries := 30
	for tries > 0 {
		receipt, _ := eth.TransactionReceipt(ctx, txHash)
		if receipt != nil {
			return nil
		}
		tries--
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("transaction (%s) can't be mined in time", txHash.Hex())
}

func mustNewArguments(types ...string) (result abi.Arguments) {
	var err error
	for _, t := range types {
		var typ abi.Type
		items := strings.Split(t, " ")
		var name string
		if len(items) == 2 {
			name = items[1]
		} else {
			name = items[0]
		}
		typ, err = abi.NewType(items[0], items[0], nil)
		if err != nil {
			panic(err)
		}
		result = append(result, abi.Argument{Type: typ, Name: name})
	}
	return result
}

func encodeFunctionCall(functionSignature string, args ...interface{}) []byte {
	functionName := functionSignature[0:strings.Index(functionSignature, "(")]
	functionTypes := functionSignature[strings.Index(functionSignature, "(")+1 : len(functionSignature)-1]
	input := mustNewArguments(strings.Split(functionTypes, ",")...)
	method := abi.NewMethod(functionName, functionName, abi.Function, "", false, false, input, nil)
	packed, err := input.Pack(args...)
	if err != nil {
		log.Fatalf("failed to pack input params: %+v", err)
	}
	return append(method.ID, packed...)
}

func prettyFormatTime(input int64) (result string) {
	years := math.Floor(float64(input) / 60 / 60 / 24 / 7 / 30 / 12)
	if years > 0 {
		result += fmt.Sprintf("%dy", int(years))
	}
	seconds := input % (60 * 60 * 24 * 7 * 30 * 12)
	months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
	if months > 0 {
		result += fmt.Sprintf("%dm", int(months))
	}
	seconds = input % (60 * 60 * 24 * 7 * 30)
	weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
	if weeks > 0 {
		result += fmt.Sprintf("%dw", int(weeks))
	}
	seconds = input % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	if days > 0 {
		result += fmt.Sprintf("%dd", int(days))
	}
	seconds = input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	if hours > 0 {
		result += fmt.Sprintf("%dh", int(hours))
	}
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	if minutes > 0 {
		result += fmt.Sprintf("%dm", int(minutes))
	}
	seconds = input % 60
	if seconds > 0 {
		result += fmt.Sprintf("%ds", int(seconds))
	}
	return
}
