package evm

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
)

// (Tx) make msg - send coin
func MakeSendCoinMsg(sendCoinMsg types.SendCoinMsg, privKey key.PrivateKey) (types.SendCoinMsg, error) {
	msg, err := parseSendCoinArgs(sendCoinMsg, privKey)
	if err != nil {
		return types.SendCoinMsg{}, err
	}

	return msg, nil
}

// (Tx) make msg - deploy solidity contract
func MakeDeploySolContractMsg(deploySolContractMsg types.DeploySolContractMsg) error {
	err := parseDeploySolContractArgs(deploySolContractMsg)
	if err != nil {
		return err
	}

	return nil
}

// (Tx) make msg - invoke solidity contract
func MakeInvokeSolContractMsg(InvokeSolContractMsg types.InvokeSolContractMsg) (types.InvokeSolContractMsg, error) {
	msg, err := parseInvokeSolContractArgs(InvokeSolContractMsg)
	if err != nil {
		return types.InvokeSolContractMsg{}, err
	}

	return msg, nil
}

// (Query) make msg - call solidity contract
func MakeCallSolContractMsg(callSolContractMsg types.CallSolContractMsg) (types.CallSolContractMsg, error) {
	msg, err := parseCallSolContractArgs(callSolContractMsg)
	if err != nil {
		return types.CallSolContractMsg{}, err
	}

	return msg, nil
}

// (Query) make msg - transaction by hash
func MakeGetTransactionByHashMsg(getTransactionByHashMsg types.GetTransactionByHashMsg) (types.GetTransactionByHashMsg, error) {
	msg, err := parseGetTransactionByHashArgs(getTransactionByHashMsg)
	if err != nil {
		return types.GetTransactionByHashMsg{}, err
	}

	return msg, nil
}

// (Query) make msg - block by hash or height
func MakeGetBlockByHashHeightMsg(getBlockByHashHeightMsg types.GetBlockByHashHeightMsg) (types.GetBlockByHashHeightMsg, error) {
	msg, err := parseGetBlockByHashHeightArgs(getBlockByHashHeightMsg)
	if err != nil {
		return types.GetBlockByHashHeightMsg{}, err
	}

	return msg, nil
}

// (Query) make msg - account info
func MakeQueryAccountInfoMsg(accountInfoMsg types.AccountInfoMsg) (types.AccountInfoMsg, error) {
	msg, err := parseQueryAccountInfoArgs(accountInfoMsg)
	if err != nil {
		return types.AccountInfoMsg{}, err
	}

	return msg, nil
}

// (Query) make msg - web3 sha3
func MakeWeb3Sha3Msg(web3Sha3Msg types.Web3Sha3Msg) (types.Web3Sha3Msg, error) {
	msg, err := parseWeb3Sha3Args(web3Sha3Msg)
	if err != nil {
		return types.Web3Sha3Msg{}, err
	}

	return msg, nil
}

// (Query) make msg - get transaction count of the block number
func MakeEthGetBlockTransactionCountMsg(ethGetBlockTransactionCountMsg types.EthGetBlockTransactionCountMsg) (types.EthGetBlockTransactionCountMsg, error) {
	msg, err := parseEthGetBlockTransactionCountArgs(ethGetBlockTransactionCountMsg)
	if err != nil {
		return types.EthGetBlockTransactionCountMsg{}, err
	}

	return msg, nil
}

// (Query) make msg - sol contract estimate gas
func MakeEstimateGasSolMsg(invokeSolContractMsg types.InvokeSolContractMsg) (types.InvokeSolContractMsg, error) {
	msg, err := parseEstimateGasSolArgs(invokeSolContractMsg)
	if err != nil {
		return types.InvokeSolContractMsg{}, err
	}

	return msg, nil
}

// (Query) make msg - get transaction by block hash and index
func MakeGetTransactionByBlockHashAndIndexMsg(getTransactionByBlockHashAndIndexMsg types.GetTransactionByBlockHashAndIndexMsg) (types.GetTransactionByBlockHashAndIndexMsg, error) {
	msg, err := parseGetTransactionByBlockHashAndIndexArgs(getTransactionByBlockHashAndIndexMsg)
	if err != nil {
		return types.GetTransactionByBlockHashAndIndexMsg{}, err
	}

	return msg, nil
}

// (Query) make msg - get transaction receipt
func MakeGetTransactionReceiptMsg(getTransactionReceiptMsg types.GetTransactionReceiptMsg) (types.GetTransactionReceiptMsg, error) {
	msg, err := parseGetTransactionReceiptArgs(getTransactionReceiptMsg)
	if err != nil {
		return types.GetTransactionReceiptMsg{}, err
	}

	return msg, nil
}
