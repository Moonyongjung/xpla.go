// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package evm

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

// XplaSolContractMetaData contains all meta data concerning the XplaSolContract contract.
var XplaSolContractMetaData *bind.MetaData

// DeployXplaSolContract deploys a new Ethereum contract, binding an instance of XplaSolContract to it.
func DeployXplaSolContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *XplaSolContract, error) {
	parsed, err := XplaSolContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(XplaSolContractMetaData.Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &XplaSolContract{XplaSolContractCaller: XplaSolContractCaller{contract: contract}, XplaSolContractTransactor: XplaSolContractTransactor{contract: contract}, XplaSolContractFilterer: XplaSolContractFilterer{contract: contract}}, nil
}

// XplaSolContract is an auto generated Go binding around an Ethereum contract.
type XplaSolContract struct {
	XplaSolContractCaller     // Read-only binding to the contract
	XplaSolContractTransactor // Write-only binding to the contract
	XplaSolContractFilterer   // Log filterer for contract events
}

// XplaSolContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type XplaSolContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XplaSolContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type XplaSolContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XplaSolContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type XplaSolContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XplaSolContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type XplaSolContractSession struct {
	Contract     *XplaSolContract  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// XplaSolContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type XplaSolContractCallerSession struct {
	Contract *XplaSolContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// XplaSolContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type XplaSolContractTransactorSession struct {
	Contract     *XplaSolContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// XplaSolContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type XplaSolContractRaw struct {
	Contract *XplaSolContract // Generic contract binding to access the raw methods on
}

// XplaSolContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type XplaSolContractCallerRaw struct {
	Contract *XplaSolContractCaller // Generic read-only contract binding to access the raw methods on
}

// XplaSolContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type XplaSolContractTransactorRaw struct {
	Contract *XplaSolContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewXplaSolContract creates a new instance of XplaSolContract, bound to a specific deployed contract.
func NewXplaSolContract(address common.Address, backend bind.ContractBackend) (*XplaSolContract, error) {
	contract, err := bindXplaSolContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &XplaSolContract{XplaSolContractCaller: XplaSolContractCaller{contract: contract}, XplaSolContractTransactor: XplaSolContractTransactor{contract: contract}, XplaSolContractFilterer: XplaSolContractFilterer{contract: contract}}, nil
}

// NewXplaSolContractCaller creates a new read-only instance of XplaSolContract, bound to a specific deployed contract.
func NewXplaSolContractCaller(address common.Address, caller bind.ContractCaller) (*XplaSolContractCaller, error) {
	contract, err := bindXplaSolContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &XplaSolContractCaller{contract: contract}, nil
}

// NewXplaSolContractTransactor creates a new write-only instance of XplaSolContract, bound to a specific deployed contract.
func NewXplaSolContractTransactor(address common.Address, transactor bind.ContractTransactor) (*XplaSolContractTransactor, error) {
	contract, err := bindXplaSolContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &XplaSolContractTransactor{contract: contract}, nil
}

// NewXplaSolContractFilterer creates a new log filterer instance of XplaSolContract, bound to a specific deployed contract.
func NewXplaSolContractFilterer(address common.Address, filterer bind.ContractFilterer) (*XplaSolContractFilterer, error) {
	contract, err := bindXplaSolContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &XplaSolContractFilterer{contract: contract}, nil
}

// bindXplaSolContract binds a generic wrapper to an already deployed contract.
func bindXplaSolContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(XplaSolContractMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_XplaSolContract *XplaSolContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _XplaSolContract.Contract.XplaSolContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_XplaSolContract *XplaSolContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _XplaSolContract.Contract.XplaSolContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_XplaSolContract *XplaSolContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _XplaSolContract.Contract.XplaSolContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_XplaSolContract *XplaSolContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _XplaSolContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_XplaSolContract *XplaSolContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _XplaSolContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_XplaSolContract *XplaSolContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _XplaSolContract.Contract.contract.Transact(opts, method, params...)
}
