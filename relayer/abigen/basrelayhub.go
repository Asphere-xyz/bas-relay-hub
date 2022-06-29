// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abigen

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// RelayHubMetaData contains all meta data concerning the RelayHub contract.
var RelayHubMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIProofVerificationFunction\",\"name\":\"defaultVerificationFunction\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"initValidatorSet\",\"type\":\"address[]\"}],\"name\":\"ChainRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"newValidatorSet\",\"type\":\"address[]\"}],\"name\":\"ValidatorSetUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"name\":\"getBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"genesisBlock\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"epochLength\",\"type\":\"uint32\"}],\"name\":\"registerCertifiedBAS\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"checkpointBlock\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"checkpointHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"checkpointSignature\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"epochLength\",\"type\":\"uint32\"}],\"name\":\"registerUsingCheckpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"contractIProofVerificationFunction\",\"name\":\"verificationFunction\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"genesisBlock\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"epochLength\",\"type\":\"uint32\"}],\"name\":\"registerBAS\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"name\":\"getActiveValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"validators\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"}],\"name\":\"checkValidators\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"uniqueValidators\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"name\":\"getLatestTransitionedEpoch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"blockProofs\",\"type\":\"bytes[]\"}],\"name\":\"updateValidatorSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"blockProofs\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"rawReceipt\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"siblings\",\"type\":\"bytes\"}],\"name\":\"checkReceiptProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// RelayHubABI is the input ABI used to generate the binding from.
// Deprecated: Use RelayHubMetaData.ABI instead.
var RelayHubABI = RelayHubMetaData.ABI

// RelayHub is an auto generated Go binding around an Ethereum contract.
type RelayHub struct {
	RelayHubCaller     // Read-only binding to the contract
	RelayHubTransactor // Write-only binding to the contract
	RelayHubFilterer   // Log filterer for contract events
}

// RelayHubCaller is an auto generated read-only Go binding around an Ethereum contract.
type RelayHubCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayHubTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RelayHubTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayHubFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RelayHubFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayHubSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RelayHubSession struct {
	Contract     *RelayHub         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RelayHubCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RelayHubCallerSession struct {
	Contract *RelayHubCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RelayHubTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RelayHubTransactorSession struct {
	Contract     *RelayHubTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RelayHubRaw is an auto generated low-level Go binding around an Ethereum contract.
type RelayHubRaw struct {
	Contract *RelayHub // Generic contract binding to access the raw methods on
}

// RelayHubCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RelayHubCallerRaw struct {
	Contract *RelayHubCaller // Generic read-only contract binding to access the raw methods on
}

// RelayHubTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RelayHubTransactorRaw struct {
	Contract *RelayHubTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRelayHub creates a new instance of RelayHub, bound to a specific deployed contract.
func NewRelayHub(address common.Address, backend bind.ContractBackend) (*RelayHub, error) {
	contract, err := bindRelayHub(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RelayHub{RelayHubCaller: RelayHubCaller{contract: contract}, RelayHubTransactor: RelayHubTransactor{contract: contract}, RelayHubFilterer: RelayHubFilterer{contract: contract}}, nil
}

// NewRelayHubCaller creates a new read-only instance of RelayHub, bound to a specific deployed contract.
func NewRelayHubCaller(address common.Address, caller bind.ContractCaller) (*RelayHubCaller, error) {
	contract, err := bindRelayHub(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RelayHubCaller{contract: contract}, nil
}

// NewRelayHubTransactor creates a new write-only instance of RelayHub, bound to a specific deployed contract.
func NewRelayHubTransactor(address common.Address, transactor bind.ContractTransactor) (*RelayHubTransactor, error) {
	contract, err := bindRelayHub(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RelayHubTransactor{contract: contract}, nil
}

// NewRelayHubFilterer creates a new log filterer instance of RelayHub, bound to a specific deployed contract.
func NewRelayHubFilterer(address common.Address, filterer bind.ContractFilterer) (*RelayHubFilterer, error) {
	contract, err := bindRelayHub(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RelayHubFilterer{contract: contract}, nil
}

// bindRelayHub binds a generic wrapper to an already deployed contract.
func bindRelayHub(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RelayHubABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RelayHub *RelayHubRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RelayHub.Contract.RelayHubCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RelayHub *RelayHubRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RelayHub.Contract.RelayHubTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RelayHub *RelayHubRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RelayHub.Contract.RelayHubTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RelayHub *RelayHubCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RelayHub.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RelayHub *RelayHubTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RelayHub.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RelayHub *RelayHubTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RelayHub.Contract.contract.Transact(opts, method, params...)
}

// CheckReceiptProof is a free data retrieval call binding the contract method 0x973ebcd8.
//
// Solidity: function checkReceiptProof(uint256 chainId, bytes[] blockProofs, bytes rawReceipt, bytes path, bytes siblings) view returns(bool)
func (_RelayHub *RelayHubCaller) CheckReceiptProof(opts *bind.CallOpts, chainId *big.Int, blockProofs [][]byte, rawReceipt []byte, path []byte, siblings []byte) (bool, error) {
	var out []interface{}
	err := _RelayHub.contract.Call(opts, &out, "checkReceiptProof", chainId, blockProofs, rawReceipt, path, siblings)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckReceiptProof is a free data retrieval call binding the contract method 0x973ebcd8.
//
// Solidity: function checkReceiptProof(uint256 chainId, bytes[] blockProofs, bytes rawReceipt, bytes path, bytes siblings) view returns(bool)
func (_RelayHub *RelayHubSession) CheckReceiptProof(chainId *big.Int, blockProofs [][]byte, rawReceipt []byte, path []byte, siblings []byte) (bool, error) {
	return _RelayHub.Contract.CheckReceiptProof(&_RelayHub.CallOpts, chainId, blockProofs, rawReceipt, path, siblings)
}

// CheckReceiptProof is a free data retrieval call binding the contract method 0x973ebcd8.
//
// Solidity: function checkReceiptProof(uint256 chainId, bytes[] blockProofs, bytes rawReceipt, bytes path, bytes siblings) view returns(bool)
func (_RelayHub *RelayHubCallerSession) CheckReceiptProof(chainId *big.Int, blockProofs [][]byte, rawReceipt []byte, path []byte, siblings []byte) (bool, error) {
	return _RelayHub.Contract.CheckReceiptProof(&_RelayHub.CallOpts, chainId, blockProofs, rawReceipt, path, siblings)
}

// CheckValidators is a free data retrieval call binding the contract method 0xab50f3d0.
//
// Solidity: function checkValidators(uint256 chainId, address[] validators, uint64 epoch) view returns(uint64 uniqueValidators)
func (_RelayHub *RelayHubCaller) CheckValidators(opts *bind.CallOpts, chainId *big.Int, validators []common.Address, epoch uint64) (uint64, error) {
	var out []interface{}
	err := _RelayHub.contract.Call(opts, &out, "checkValidators", chainId, validators, epoch)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// CheckValidators is a free data retrieval call binding the contract method 0xab50f3d0.
//
// Solidity: function checkValidators(uint256 chainId, address[] validators, uint64 epoch) view returns(uint64 uniqueValidators)
func (_RelayHub *RelayHubSession) CheckValidators(chainId *big.Int, validators []common.Address, epoch uint64) (uint64, error) {
	return _RelayHub.Contract.CheckValidators(&_RelayHub.CallOpts, chainId, validators, epoch)
}

// CheckValidators is a free data retrieval call binding the contract method 0xab50f3d0.
//
// Solidity: function checkValidators(uint256 chainId, address[] validators, uint64 epoch) view returns(uint64 uniqueValidators)
func (_RelayHub *RelayHubCallerSession) CheckValidators(chainId *big.Int, validators []common.Address, epoch uint64) (uint64, error) {
	return _RelayHub.Contract.CheckValidators(&_RelayHub.CallOpts, chainId, validators, epoch)
}

// GetActiveValidators is a free data retrieval call binding the contract method 0xd5fa9e49.
//
// Solidity: function getActiveValidators(uint256 chainId) view returns(address[])
func (_RelayHub *RelayHubCaller) GetActiveValidators(opts *bind.CallOpts, chainId *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _RelayHub.contract.Call(opts, &out, "getActiveValidators", chainId)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetActiveValidators is a free data retrieval call binding the contract method 0xd5fa9e49.
//
// Solidity: function getActiveValidators(uint256 chainId) view returns(address[])
func (_RelayHub *RelayHubSession) GetActiveValidators(chainId *big.Int) ([]common.Address, error) {
	return _RelayHub.Contract.GetActiveValidators(&_RelayHub.CallOpts, chainId)
}

// GetActiveValidators is a free data retrieval call binding the contract method 0xd5fa9e49.
//
// Solidity: function getActiveValidators(uint256 chainId) view returns(address[])
func (_RelayHub *RelayHubCallerSession) GetActiveValidators(chainId *big.Int) ([]common.Address, error) {
	return _RelayHub.Contract.GetActiveValidators(&_RelayHub.CallOpts, chainId)
}

// GetBridgeAddress is a free data retrieval call binding the contract method 0x0c46a0e1.
//
// Solidity: function getBridgeAddress(uint256 chainId) view returns(address)
func (_RelayHub *RelayHubCaller) GetBridgeAddress(opts *bind.CallOpts, chainId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _RelayHub.contract.Call(opts, &out, "getBridgeAddress", chainId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBridgeAddress is a free data retrieval call binding the contract method 0x0c46a0e1.
//
// Solidity: function getBridgeAddress(uint256 chainId) view returns(address)
func (_RelayHub *RelayHubSession) GetBridgeAddress(chainId *big.Int) (common.Address, error) {
	return _RelayHub.Contract.GetBridgeAddress(&_RelayHub.CallOpts, chainId)
}

// GetBridgeAddress is a free data retrieval call binding the contract method 0x0c46a0e1.
//
// Solidity: function getBridgeAddress(uint256 chainId) view returns(address)
func (_RelayHub *RelayHubCallerSession) GetBridgeAddress(chainId *big.Int) (common.Address, error) {
	return _RelayHub.Contract.GetBridgeAddress(&_RelayHub.CallOpts, chainId)
}

// GetLatestTransitionedEpoch is a free data retrieval call binding the contract method 0xed594dd1.
//
// Solidity: function getLatestTransitionedEpoch(uint256 chainId) view returns(uint64)
func (_RelayHub *RelayHubCaller) GetLatestTransitionedEpoch(opts *bind.CallOpts, chainId *big.Int) (uint64, error) {
	var out []interface{}
	err := _RelayHub.contract.Call(opts, &out, "getLatestTransitionedEpoch", chainId)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetLatestTransitionedEpoch is a free data retrieval call binding the contract method 0xed594dd1.
//
// Solidity: function getLatestTransitionedEpoch(uint256 chainId) view returns(uint64)
func (_RelayHub *RelayHubSession) GetLatestTransitionedEpoch(chainId *big.Int) (uint64, error) {
	return _RelayHub.Contract.GetLatestTransitionedEpoch(&_RelayHub.CallOpts, chainId)
}

// GetLatestTransitionedEpoch is a free data retrieval call binding the contract method 0xed594dd1.
//
// Solidity: function getLatestTransitionedEpoch(uint256 chainId) view returns(uint64)
func (_RelayHub *RelayHubCallerSession) GetLatestTransitionedEpoch(chainId *big.Int) (uint64, error) {
	return _RelayHub.Contract.GetLatestTransitionedEpoch(&_RelayHub.CallOpts, chainId)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_RelayHub *RelayHubTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _RelayHub.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_RelayHub *RelayHubSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _RelayHub.Contract.Multicall(&_RelayHub.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_RelayHub *RelayHubTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _RelayHub.Contract.Multicall(&_RelayHub.TransactOpts, data)
}

// RegisterBAS is a paid mutator transaction binding the contract method 0x45f87d3c.
//
// Solidity: function registerBAS(uint256 chainId, address verificationFunction, bytes genesisBlock, address bridgeAddress, uint32 epochLength) returns()
func (_RelayHub *RelayHubTransactor) RegisterBAS(opts *bind.TransactOpts, chainId *big.Int, verificationFunction common.Address, genesisBlock []byte, bridgeAddress common.Address, epochLength uint32) (*types.Transaction, error) {
	return _RelayHub.contract.Transact(opts, "registerBAS", chainId, verificationFunction, genesisBlock, bridgeAddress, epochLength)
}

// RegisterBAS is a paid mutator transaction binding the contract method 0x45f87d3c.
//
// Solidity: function registerBAS(uint256 chainId, address verificationFunction, bytes genesisBlock, address bridgeAddress, uint32 epochLength) returns()
func (_RelayHub *RelayHubSession) RegisterBAS(chainId *big.Int, verificationFunction common.Address, genesisBlock []byte, bridgeAddress common.Address, epochLength uint32) (*types.Transaction, error) {
	return _RelayHub.Contract.RegisterBAS(&_RelayHub.TransactOpts, chainId, verificationFunction, genesisBlock, bridgeAddress, epochLength)
}

// RegisterBAS is a paid mutator transaction binding the contract method 0x45f87d3c.
//
// Solidity: function registerBAS(uint256 chainId, address verificationFunction, bytes genesisBlock, address bridgeAddress, uint32 epochLength) returns()
func (_RelayHub *RelayHubTransactorSession) RegisterBAS(chainId *big.Int, verificationFunction common.Address, genesisBlock []byte, bridgeAddress common.Address, epochLength uint32) (*types.Transaction, error) {
	return _RelayHub.Contract.RegisterBAS(&_RelayHub.TransactOpts, chainId, verificationFunction, genesisBlock, bridgeAddress, epochLength)
}

// RegisterCertifiedBAS is a paid mutator transaction binding the contract method 0xfae2f0b5.
//
// Solidity: function registerCertifiedBAS(uint256 chainId, bytes genesisBlock, address bridgeAddress, uint32 epochLength) returns()
func (_RelayHub *RelayHubTransactor) RegisterCertifiedBAS(opts *bind.TransactOpts, chainId *big.Int, genesisBlock []byte, bridgeAddress common.Address, epochLength uint32) (*types.Transaction, error) {
	return _RelayHub.contract.Transact(opts, "registerCertifiedBAS", chainId, genesisBlock, bridgeAddress, epochLength)
}

// RegisterCertifiedBAS is a paid mutator transaction binding the contract method 0xfae2f0b5.
//
// Solidity: function registerCertifiedBAS(uint256 chainId, bytes genesisBlock, address bridgeAddress, uint32 epochLength) returns()
func (_RelayHub *RelayHubSession) RegisterCertifiedBAS(chainId *big.Int, genesisBlock []byte, bridgeAddress common.Address, epochLength uint32) (*types.Transaction, error) {
	return _RelayHub.Contract.RegisterCertifiedBAS(&_RelayHub.TransactOpts, chainId, genesisBlock, bridgeAddress, epochLength)
}

// RegisterCertifiedBAS is a paid mutator transaction binding the contract method 0xfae2f0b5.
//
// Solidity: function registerCertifiedBAS(uint256 chainId, bytes genesisBlock, address bridgeAddress, uint32 epochLength) returns()
func (_RelayHub *RelayHubTransactorSession) RegisterCertifiedBAS(chainId *big.Int, genesisBlock []byte, bridgeAddress common.Address, epochLength uint32) (*types.Transaction, error) {
	return _RelayHub.Contract.RegisterCertifiedBAS(&_RelayHub.TransactOpts, chainId, genesisBlock, bridgeAddress, epochLength)
}

// RegisterUsingCheckpoint is a paid mutator transaction binding the contract method 0x7b69eff0.
//
// Solidity: function registerUsingCheckpoint(uint256 chainId, bytes checkpointBlock, bytes32 checkpointHash, bytes checkpointSignature, address bridgeAddress, uint32 epochLength) returns()
func (_RelayHub *RelayHubTransactor) RegisterUsingCheckpoint(opts *bind.TransactOpts, chainId *big.Int, checkpointBlock []byte, checkpointHash [32]byte, checkpointSignature []byte, bridgeAddress common.Address, epochLength uint32) (*types.Transaction, error) {
	return _RelayHub.contract.Transact(opts, "registerUsingCheckpoint", chainId, checkpointBlock, checkpointHash, checkpointSignature, bridgeAddress, epochLength)
}

// RegisterUsingCheckpoint is a paid mutator transaction binding the contract method 0x7b69eff0.
//
// Solidity: function registerUsingCheckpoint(uint256 chainId, bytes checkpointBlock, bytes32 checkpointHash, bytes checkpointSignature, address bridgeAddress, uint32 epochLength) returns()
func (_RelayHub *RelayHubSession) RegisterUsingCheckpoint(chainId *big.Int, checkpointBlock []byte, checkpointHash [32]byte, checkpointSignature []byte, bridgeAddress common.Address, epochLength uint32) (*types.Transaction, error) {
	return _RelayHub.Contract.RegisterUsingCheckpoint(&_RelayHub.TransactOpts, chainId, checkpointBlock, checkpointHash, checkpointSignature, bridgeAddress, epochLength)
}

// RegisterUsingCheckpoint is a paid mutator transaction binding the contract method 0x7b69eff0.
//
// Solidity: function registerUsingCheckpoint(uint256 chainId, bytes checkpointBlock, bytes32 checkpointHash, bytes checkpointSignature, address bridgeAddress, uint32 epochLength) returns()
func (_RelayHub *RelayHubTransactorSession) RegisterUsingCheckpoint(chainId *big.Int, checkpointBlock []byte, checkpointHash [32]byte, checkpointSignature []byte, bridgeAddress common.Address, epochLength uint32) (*types.Transaction, error) {
	return _RelayHub.Contract.RegisterUsingCheckpoint(&_RelayHub.TransactOpts, chainId, checkpointBlock, checkpointHash, checkpointSignature, bridgeAddress, epochLength)
}

// UpdateValidatorSet is a paid mutator transaction binding the contract method 0xf874421e.
//
// Solidity: function updateValidatorSet(uint256 chainId, bytes[] blockProofs) returns()
func (_RelayHub *RelayHubTransactor) UpdateValidatorSet(opts *bind.TransactOpts, chainId *big.Int, blockProofs [][]byte) (*types.Transaction, error) {
	return _RelayHub.contract.Transact(opts, "updateValidatorSet", chainId, blockProofs)
}

// UpdateValidatorSet is a paid mutator transaction binding the contract method 0xf874421e.
//
// Solidity: function updateValidatorSet(uint256 chainId, bytes[] blockProofs) returns()
func (_RelayHub *RelayHubSession) UpdateValidatorSet(chainId *big.Int, blockProofs [][]byte) (*types.Transaction, error) {
	return _RelayHub.Contract.UpdateValidatorSet(&_RelayHub.TransactOpts, chainId, blockProofs)
}

// UpdateValidatorSet is a paid mutator transaction binding the contract method 0xf874421e.
//
// Solidity: function updateValidatorSet(uint256 chainId, bytes[] blockProofs) returns()
func (_RelayHub *RelayHubTransactorSession) UpdateValidatorSet(chainId *big.Int, blockProofs [][]byte) (*types.Transaction, error) {
	return _RelayHub.Contract.UpdateValidatorSet(&_RelayHub.TransactOpts, chainId, blockProofs)
}

// RelayHubChainRegisteredIterator is returned from FilterChainRegistered and is used to iterate over the raw logs and unpacked data for ChainRegistered events raised by the RelayHub contract.
type RelayHubChainRegisteredIterator struct {
	Event *RelayHubChainRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayHubChainRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayHubChainRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayHubChainRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayHubChainRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayHubChainRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayHubChainRegistered represents a ChainRegistered event raised by the RelayHub contract.
type RelayHubChainRegistered struct {
	ChainId          *big.Int
	InitValidatorSet []common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterChainRegistered is a free log retrieval operation binding the contract event 0x81c059162914deb2e3f6468c0aed7e2e09b5f62ae648456a1735f32775bb316c.
//
// Solidity: event ChainRegistered(uint256 indexed chainId, address[] initValidatorSet)
func (_RelayHub *RelayHubFilterer) FilterChainRegistered(opts *bind.FilterOpts, chainId []*big.Int) (*RelayHubChainRegisteredIterator, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _RelayHub.contract.FilterLogs(opts, "ChainRegistered", chainIdRule)
	if err != nil {
		return nil, err
	}
	return &RelayHubChainRegisteredIterator{contract: _RelayHub.contract, event: "ChainRegistered", logs: logs, sub: sub}, nil
}

// WatchChainRegistered is a free log subscription operation binding the contract event 0x81c059162914deb2e3f6468c0aed7e2e09b5f62ae648456a1735f32775bb316c.
//
// Solidity: event ChainRegistered(uint256 indexed chainId, address[] initValidatorSet)
func (_RelayHub *RelayHubFilterer) WatchChainRegistered(opts *bind.WatchOpts, sink chan<- *RelayHubChainRegistered, chainId []*big.Int) (event.Subscription, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _RelayHub.contract.WatchLogs(opts, "ChainRegistered", chainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayHubChainRegistered)
				if err := _RelayHub.contract.UnpackLog(event, "ChainRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChainRegistered is a log parse operation binding the contract event 0x81c059162914deb2e3f6468c0aed7e2e09b5f62ae648456a1735f32775bb316c.
//
// Solidity: event ChainRegistered(uint256 indexed chainId, address[] initValidatorSet)
func (_RelayHub *RelayHubFilterer) ParseChainRegistered(log types.Log) (*RelayHubChainRegistered, error) {
	event := new(RelayHubChainRegistered)
	if err := _RelayHub.contract.UnpackLog(event, "ChainRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayHubValidatorSetUpdatedIterator is returned from FilterValidatorSetUpdated and is used to iterate over the raw logs and unpacked data for ValidatorSetUpdated events raised by the RelayHub contract.
type RelayHubValidatorSetUpdatedIterator struct {
	Event *RelayHubValidatorSetUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayHubValidatorSetUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayHubValidatorSetUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayHubValidatorSetUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayHubValidatorSetUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayHubValidatorSetUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayHubValidatorSetUpdated represents a ValidatorSetUpdated event raised by the RelayHub contract.
type RelayHubValidatorSetUpdated struct {
	ChainId         *big.Int
	NewValidatorSet []common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorSetUpdated is a free log retrieval operation binding the contract event 0x3d0eea40644a206ec25781dd5bb3b60eb4fa1264b993c3bddf3c73b14f29ef5e.
//
// Solidity: event ValidatorSetUpdated(uint256 indexed chainId, address[] newValidatorSet)
func (_RelayHub *RelayHubFilterer) FilterValidatorSetUpdated(opts *bind.FilterOpts, chainId []*big.Int) (*RelayHubValidatorSetUpdatedIterator, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _RelayHub.contract.FilterLogs(opts, "ValidatorSetUpdated", chainIdRule)
	if err != nil {
		return nil, err
	}
	return &RelayHubValidatorSetUpdatedIterator{contract: _RelayHub.contract, event: "ValidatorSetUpdated", logs: logs, sub: sub}, nil
}

// WatchValidatorSetUpdated is a free log subscription operation binding the contract event 0x3d0eea40644a206ec25781dd5bb3b60eb4fa1264b993c3bddf3c73b14f29ef5e.
//
// Solidity: event ValidatorSetUpdated(uint256 indexed chainId, address[] newValidatorSet)
func (_RelayHub *RelayHubFilterer) WatchValidatorSetUpdated(opts *bind.WatchOpts, sink chan<- *RelayHubValidatorSetUpdated, chainId []*big.Int) (event.Subscription, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _RelayHub.contract.WatchLogs(opts, "ValidatorSetUpdated", chainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayHubValidatorSetUpdated)
				if err := _RelayHub.contract.UnpackLog(event, "ValidatorSetUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorSetUpdated is a log parse operation binding the contract event 0x3d0eea40644a206ec25781dd5bb3b60eb4fa1264b993c3bddf3c73b14f29ef5e.
//
// Solidity: event ValidatorSetUpdated(uint256 indexed chainId, address[] newValidatorSet)
func (_RelayHub *RelayHubFilterer) ParseValidatorSetUpdated(log types.Log) (*RelayHubValidatorSetUpdated, error) {
	event := new(RelayHubValidatorSetUpdated)
	if err := _RelayHub.contract.UnpackLog(event, "ValidatorSetUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
