package client

import (
	mevm "github.com/Moonyongjung/xpla.go/core/evm"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
)

var xplaTxRes types.TxRes

// Broadcast the transaction.
// Default broadcast mode is "sync" if not xpla client has broadcast mode option.
// The broadcast method is determined according to the broadcast mode option of the xpla client.
// For evm transaction broadcast, use a separate method in this function.
func (xplac *XplaClient) Broadcast(txBytes []byte) (*types.TxRes, error) {

	if xplac.Module == mevm.EvmModule {
		return xplac.broadcastEvm(txBytes)

	} else {
		broadcastMode := xplac.Opts.BroadcastMode
		switch {
		case broadcastMode == "block":
			return xplac.BroadcastBlock(txBytes)
		case broadcastMode == "async":
			return xplac.BroadcastAsync(txBytes)
		case broadcastMode == "sync":
			return broadcastTx(xplac, txBytes, txtypes.BroadcastMode_BROADCAST_MODE_SYNC)
		default:
			return broadcastTx(xplac, txBytes, txtypes.BroadcastMode_BROADCAST_MODE_SYNC)
		}
	}
}

// Broadcast the transaction with mode "block".
// It takes precedence over the option of the xpla client.
func (xplac *XplaClient) BroadcastBlock(txBytes []byte) (*types.TxRes, error) {
	return broadcastTx(xplac, txBytes, txtypes.BroadcastMode_BROADCAST_MODE_BLOCK)
}

// Broadcast the transaction with mode "Async".
// It takes precedence over the option of the xpla client.
func (xplac *XplaClient) BroadcastAsync(txBytes []byte) (*types.TxRes, error) {
	return broadcastTx(xplac, txBytes, txtypes.BroadcastMode_BROADCAST_MODE_ASYNC)
}

// Broadcast the transaction which is evm transaction by using ethclient of go-ethereum.
func (xplac *XplaClient) broadcastEvm(txBytes []byte) (*types.TxRes, error) {
	evmClient, err := util.NewEvmClient(xplac.Opts.EvmRpcURL, xplac.Context)
	if err != nil {
		return nil, err
	}
	broadcastMode := xplac.Opts.BroadcastMode
	return broadcastTxEvm(xplac, txBytes, broadcastMode, evmClient)
}
