package evm

const (
	EvmModule                         = "evm"
	EvmSendCoinMsgType                = "evm-send-coin"
	EvmDeploySolContractMsgType       = "deploy-sol-contract"
	EvmInvokeSolContractMsgType       = "invoke-sol-contract"
	EvmCallSolContractMsgType         = "call-sol-contract"
	EvmGetTransactionByHashMsgType    = "evm-get-transaction-by-hash"
	EvmGetBlockByHashHeightMsgType    = "evm-get-block"
	EvmQueryAccountInfoMsgType        = "evm-query-account-info"
	EvmSuggestGasPriceMsgType         = "suggest-gas-price"
	EvmQueryChainIdMsgType            = "evm-chain-id"
	EvmQueryCurrentBlockNumberMsgType = "current-block-number"
)
