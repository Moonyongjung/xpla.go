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
