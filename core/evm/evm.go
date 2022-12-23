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
func MakeCallSolContractMsg(callSolContractMsg types.CallSolContractMsg, byteAddress string) (CallSolContractParseMsg, error) {
	msg, err := parseCallSolContractArgs(callSolContractMsg, byteAddress)
	if err != nil {
		return CallSolContractParseMsg{}, err
	}

	return msg, nil
}

// (Query) make msg - transaction by hash
func MakeGetTransactionByHashMsg(getTransactionByHashMsg types.GetTransactionByHashMsg) (types.GetTransactionByHashMsg, error) {
	return getTransactionByHashMsg, nil
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
	return accountInfoMsg, nil
}

// (Query) make msg - web3 sha3
func MakeWeb3Sha3Msg(web3Sha3Msg types.Web3Sha3Msg) (types.Web3Sha3Msg, error) {
	return web3Sha3Msg, nil
}

// (Query) make msg - get transaction count of the block number
func MakeEthGetBlockTransactionCountMsg(ethGetBlockTransactionCountMsg types.EthGetBlockTransactionCountMsg) (types.EthGetBlockTransactionCountMsg, error) {
	return ethGetBlockTransactionCountMsg, nil
}

// (Query) make msg - sol contract estimate gas
func MakeEstimateGasSolMsg(invokeSolContractMsg types.InvokeSolContractMsg, byteAddress string) (CallSolContractParseMsg, error) {
	msg, err := parseEstimateGasSolArgs(invokeSolContractMsg, byteAddress)
	if err != nil {
		return CallSolContractParseMsg{}, err
	}

	return msg, nil
}

// (Query) make msg - get transaction by block hash and index
func MakeGetTransactionByBlockHashAndIndexMsg(getTransactionByBlockHashAndIndexMsg types.GetTransactionByBlockHashAndIndexMsg) (types.GetTransactionByBlockHashAndIndexMsg, error) {
	return getTransactionByBlockHashAndIndexMsg, nil
}

// (Query) make msg - get transaction receipt
func MakeGetTransactionReceiptMsg(getTransactionReceiptMsg types.GetTransactionReceiptMsg) (types.GetTransactionReceiptMsg, error) {
	return getTransactionReceiptMsg, nil
}

// (Query) make msg - eth new filter
func MakeEthNewFilterMsg(ethNewFilterMsg types.EthNewFilterMsg) (EthNewFilterParseMsg, error) {
	msg, err := parseEthNewFilterArgs(ethNewFilterMsg)
	if err != nil {
		return EthNewFilterParseMsg{}, err
	}

	return msg, nil
}

// (Query) make msg - eth uninstall filter
func MakeEthUninstallFilterMsg(ethUninsatllFilter types.EthUninsatllFilterMsg) (types.EthUninsatllFilterMsg, error) {
	return ethUninsatllFilter, nil
}

// (Query) make msg - eth get filter changes
func MakeEthGetFilterChangesMsg(ethGetFilterChangesMsg types.EthGetFilterChangesMsg) (types.EthGetFilterChangesMsg, error) {
	return ethGetFilterChangesMsg, nil
}

// (Query) make msg - eth get filter logs
func MakeEthGetFilterLogsMsg(ethGetFilterLogsMsg types.EthGetFilterLogsMsg) (types.EthGetFilterLogsMsg, error) {
	return ethGetFilterLogsMsg, nil
}

// (Query) make msg - eth get logs
func MakeEthGetLogsMsg(ethGetLogsMsg types.EthGetLogsMsg) (EthNewFilterParseMsg, error) {
	msg, err := parseEthGetLogsArgs(ethGetLogsMsg)
	if err != nil {
		return EthNewFilterParseMsg{}, err
	}

	return msg, nil
}
