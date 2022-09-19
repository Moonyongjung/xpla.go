package evm

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Parsing - send coin
func parseSendCoinArgs(sendCoinMsg types.SendCoinMsg, privKey key.PrivateKey) (types.SendCoinMsg, error) {
	addr := util.GetAddrByPrivKey(privKey)
	if sendCoinMsg.FromAddress != addr.String() {
		return types.SendCoinMsg{}, util.LogErr("Account address generated by private key is not equal")
	}

	sendCoinMsg.Amount = util.DenomRemove(sendCoinMsg.Amount)
	return sendCoinMsg, nil
}

// Parsing - deploy solidity contract
func parseDeploySolContractArgs(deploySolContractMsg types.DeploySolContractMsg) error {
	var err error
	bytecode := deploySolContractMsg.Bytecode
	if deploySolContractMsg.BytecodeJsonFilePath != "" {
		bytecode = util.BytecodeParsing(deploySolContractMsg.BytecodeJsonFilePath)
	}

	abi := deploySolContractMsg.ABI
	if deploySolContractMsg.ABIJsonFilePath != "" {
		abi, err = util.AbiParsing(deploySolContractMsg.ABIJsonFilePath)
		if err != nil {
			return err
		}
	}

	if abi == "" || bytecode == "" {
		return util.LogErr("empty parameters, need ABI and bytecode")
	}

	XplaSolContractMetaData = &bind.MetaData{
		ABI: abi,
		Bin: bytecode,
	}

	return nil
}

// Parsing - invoke solidity contract
func parseInvokeSolContractArgs(invokeSolContractMsg types.InvokeSolContractMsg) (types.InvokeSolContractMsg, error) {
	var err error
	bytecode := invokeSolContractMsg.Bytecode
	if invokeSolContractMsg.BytecodeJsonFilePath != "" {
		bytecode = util.BytecodeParsing(invokeSolContractMsg.BytecodeJsonFilePath)
	}

	abi := invokeSolContractMsg.ABI
	if invokeSolContractMsg.ABIJsonFilePath != "" {
		abi, err = util.AbiParsing(invokeSolContractMsg.ABIJsonFilePath)
		if err != nil {
			return types.InvokeSolContractMsg{}, err
		}
	}
	invokeSolContractMsg.ContractAddress = util.ToTypeHexString(invokeSolContractMsg.ContractAddress)
	XplaSolContractMetaData = &bind.MetaData{
		ABI: abi,
		Bin: bytecode,
	}

	return invokeSolContractMsg, nil
}

// Parsing - call solidity contract
func parseCallSolContractArgs(callSolContractMsg types.CallSolContractMsg) (types.CallSolContractMsg, error) {
	var err error
	bytecode := callSolContractMsg.Bytecode
	if callSolContractMsg.BytecodeJsonFilePath != "" {
		bytecode = util.BytecodeParsing(callSolContractMsg.BytecodeJsonFilePath)
	}

	abi := callSolContractMsg.ABI
	if callSolContractMsg.ABIJsonFilePath != "" {
		abi, err = util.AbiParsing(callSolContractMsg.ABIJsonFilePath)
		if err != nil {
			return types.CallSolContractMsg{}, err
		}
	}
	callSolContractMsg.ContractAddress = util.ToTypeHexString(callSolContractMsg.ContractAddress)
	XplaSolContractMetaData = &bind.MetaData{
		ABI: abi,
		Bin: bytecode,
	}

	return callSolContractMsg, nil
}

// Parsing - transaction by hash
func parseGetTransactionByHashArgs(getTransactionByHashMsg types.GetTransactionByHashMsg) (types.GetTransactionByHashMsg, error) {
	return getTransactionByHashMsg, nil
}

// Parsing - block by hash or height
func parseGetBlockByHashHeightArgs(getBlockByHashHeightMsg types.GetBlockByHashHeightMsg) (types.GetBlockByHashHeightMsg, error) {
	if getBlockByHashHeightMsg.BlockHash != "" && getBlockByHashHeightMsg.BlockHeight != "" {
		return types.GetBlockByHashHeightMsg{}, util.LogErr("need only one parameter, hash or height")
	}

	return getBlockByHashHeightMsg, nil
}

// Parsing - account info
func parseQueryAccountInfoArgs(accountInfoMsg types.AccountInfoMsg) (types.AccountInfoMsg, error) {
	return accountInfoMsg, nil
}
