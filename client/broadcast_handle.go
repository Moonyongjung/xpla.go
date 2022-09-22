package client

import (
	"context"
	"encoding/json"
	"time"

	mevm "github.com/Moonyongjung/xpla.go/core/evm"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
)

// Broadcast generated transactions.
// Broadcast responses, excluding evm, are delivered as "TxResponse" of the entire response structure of the xpla client.
// Support broadcast by using LCD and gRPC at the same time. Default method is gRPC.
func broadcastTx(xplac *XplaClient, txBytes []byte, mode txtypes.BroadcastMode) (*types.TxRes, error) {

	broadcastReq := txtypes.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    mode,
	}

	if xplac.Opts.GrpcURL == "" {
		reqBytes, err := json.Marshal(broadcastReq)
		if err != nil {
			return nil, util.LogErr(err, "failed to marshal")
		}

		out, err := ctxHttpClient("POST", xplac.Opts.LcdURL+broadcastUrl, reqBytes)
		if err != nil {
			return nil, err
		}

		var broadcastTxResponse txtypes.BroadcastTxResponse
		err = xplac.EncodingConfig.Marshaler.UnmarshalJSON(out, &broadcastTxResponse)
		if err != nil {
			return nil, util.LogErr(err, "failed to unmarshal response")
		}

		txResponse := broadcastTxResponse.TxResponse
		if txResponse.Code != 0 {
			return &xplaTxRes, util.LogErr("tx failed with code", txResponse.Code, ":", txResponse.RawLog)
		}

		xplaTxRes.Response = txResponse
	} else {
		txClient := txtypes.NewServiceClient(xplac.Grpc)
		txResponse, err := txClient.BroadcastTx(context.Background(), &broadcastReq)
		if err != nil {
			return nil, err
		}
		xplaTxRes.Response = txResponse.TxResponse
	}

	return &xplaTxRes, nil
}

// Broadcast generated transactions of ethereum type.
// Broadcast responses, including evm, are delivered as "TxResponse".
func broadcastTxEvm(xplac *XplaClient, txBytes []byte, broadcastMode string, evmClient *EvmClient) (*types.TxRes, error) {
	switch {
	case xplac.MsgType == mevm.EvmSendCoinMsgType ||
		xplac.MsgType == mevm.EvmInvokeSolContractMsgType:
		var signedTx evmtypes.Transaction
		err := signedTx.UnmarshalJSON(txBytes)
		if err != nil {
			return nil, err
		}

		err = evmClient.Client.SendTransaction(evmClient.Ctx, &signedTx)
		if err != nil {
			return nil, err
		}

		return checkEvmBroadcastMode(broadcastMode, evmClient, &signedTx)

	case xplac.MsgType == mevm.EvmDeploySolContractMsgType:
		var deployTx deploySolTx

		err := json.Unmarshal(txBytes, &deployTx)
		if err != nil {
			return nil, err
		}

		ethPrivKey, err := toECDSA(xplac.PrivateKey)
		if err != nil {
			return nil, err
		}

		contractAuth, err := bind.NewKeyedTransactorWithChainID(ethPrivKey, deployTx.ChainId)
		if err != nil {
			return nil, err
		}
		contractAuth.Nonce = deployTx.Nonce
		contractAuth.Value = deployTx.Value
		contractAuth.GasLimit = deployTx.GasLimit
		contractAuth.GasPrice = deployTx.GasPrice

		_, transaction, _, err := mevm.DeployXplaSolContract(contractAuth, evmClient.Client)
		if err != nil {
			return nil, err
		}

		return checkEvmBroadcastMode(broadcastMode, evmClient, transaction)

	default:
		return nil, util.LogErr("invalid evm msg type")
	}
}

// Handle evm broadcast mode.
// Similarly, determine broadcast mode included in the options of xpla client.
func checkEvmBroadcastMode(broadcastMode string, evmClient *EvmClient, tx *evmtypes.Transaction) (*types.TxRes, error) {
	// Wait tx receipt (Broadcast Block)
	if broadcastMode == "block" {
		receipt, err := waitTxReceipt(evmClient, tx)
		if err != nil {
			return nil, err
		}
		xplaTxRes.EvmReceipt = receipt
		return &xplaTxRes, nil
	} else {
		return nil, nil
	}
}

// If broadcast mode is "block", client waits transaction receipt of evm.
// The wait time is in seconds and is set inside the library as timeout count.
// When the timeout occurs, no longer wait for the transaction receipt.
func waitTxReceipt(evmClient *EvmClient, signedTx *evmtypes.Transaction) (*evmtypes.Receipt, error) {
	count := defaultEvmTxReceiptTimeout
	for {
		receipt, err := evmClient.Client.TransactionReceipt(evmClient.Ctx, signedTx.Hash())
		if err != nil {
			count = count - 1
			if count < 0 {
				return nil, err
			}
			time.Sleep(time.Second * 1)
		} else {
			return receipt, nil
		}
	}
}
