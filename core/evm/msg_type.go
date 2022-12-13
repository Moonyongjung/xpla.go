package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	EvmModule                                   = "evm"
	EvmSendCoinMsgType                          = "evm-send-coin"
	EvmDeploySolContractMsgType                 = "deploy-sol-contract"
	EvmInvokeSolContractMsgType                 = "invoke-sol-contract"
	EvmCallSolContractMsgType                   = "call-sol-contract"
	EvmGetTransactionByHashMsgType              = "evm-get-transaction-by-hash"
	EvmGetBlockByHashHeightMsgType              = "evm-get-block"
	EvmQueryAccountInfoMsgType                  = "evm-query-account-info"
	EvmSuggestGasPriceMsgType                   = "suggest-gas-price"
	EvmQueryChainIdMsgType                      = "evm-chain-id"
	EvmQueryCurrentBlockNumberMsgType           = "current-block-number"
	EvmWeb3ClientVersionMsgType                 = "web3-client-version"
	EvmWeb3Sha3MsgType                          = "web3-sha"
	EvmNetVersionMsgType                        = "net-version"
	EvmNetPeerCountMsgType                      = "net-peer-count"
	EvmNetListeningMsgType                      = "net-listening"
	EvmEthProtocolVersionMsgType                = "eth-protocol-version"
	EvmEthSyncingMsgType                        = "eth-syncing"
	EvmEthAccountsMsgType                       = "eth-accounts"
	EvmEthGetBlockTransactionCountMsgType       = "eth-get-block-transaction-count"
	EvmEthEstimateGasMsgType                    = "eth-estimate-gas"
	EvmGetTransactionByBlockHashAndIndexMsgType = "eth-get-transaction-by-block-hash-and-index"
	EvmGetTransactionReceiptMsgType             = "eth-get-transaction-receipt"
	EvmEthNewFilterMsgType                      = "eth-new-filter"
	EvmEthNewBlockFilterMsgType                 = "eth-new-block-filter"
	EvmEthNewPendingTransactionFilterMsgType    = "eth-new-pending-transaction-filter"
	EvmEthUninsatllFilterMsgType                = "eth-uninstall-filter"
	EvmEthGetFilterChangesMsgType               = "eth-get-filter-changes"
	EvmEthGetFilterLogsMsgType                  = "eth-get-filter-logs"
	EvmEthGetLogsMsgType                        = "eth-get-logs"
	EvmEthCoinbaseMsgType                       = "eth-coinbase"
)

type EthNewFilterParseMsg struct {
	BlockHash *common.Hash     `json:"blockHash,omitempty"`
	FromBlock *rpc.BlockNumber `json:"fromBlock"`
	ToBlock   *rpc.BlockNumber `json:"toBlock"`
	Addresses interface{}      `json:"address"`
	Topics    []interface{}    `json:"topics"`
}
