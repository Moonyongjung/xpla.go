package client

import (
	"context"
	"io/ioutil"
	"strings"

	mwasm "github.com/Moonyongjung/xpla.go/core/wasm"
	"github.com/Moonyongjung/xpla.go/util"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// Query client for wasm module.
func queryWasm(xplac *XplaClient) (string, error) {
	queryClient := wasmtypes.NewQueryClient(xplac.Grpc)

	switch {
	// Wasm query contract
	case xplac.MsgType == mwasm.WasmQueryContractMsgType:
		convertMsg, _ := xplac.Msg.(wasmtypes.QuerySmartContractStateRequest)
		res, err = queryClient.SmartContractState(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm list code
	case xplac.MsgType == mwasm.WasmListCodeMsgType:
		convertMsg, _ := xplac.Msg.(wasmtypes.QueryCodesRequest)
		res, err = queryClient.Codes(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm list contract by code
	case xplac.MsgType == mwasm.WasmListContractByCodeMsgType:
		convertMsg, _ := xplac.Msg.(wasmtypes.QueryContractsByCodeRequest)
		res, err = queryClient.ContractsByCode(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm download
	case xplac.MsgType == mwasm.WasmDownloadMsgType:
		convertMsg, _ := xplac.Msg.([]interface{})[0].(wasmtypes.QueryCodeRequest)
		downloadFileName, _ := xplac.Msg.([]interface{})[1].(string)
		if !strings.Contains(downloadFileName, ".wasm") {
			downloadFileName = downloadFileName + ".wasm"
		}
		res, err := queryClient.Code(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}
		ioutil.WriteFile(downloadFileName, res.Data, 0o600)

	// Wasm code info
	case xplac.MsgType == mwasm.WasmCodeInfoMsgType:
		convertMsg, _ := xplac.Msg.(wasmtypes.QueryCodeRequest)
		res, err = queryClient.Code(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm contract info
	case xplac.MsgType == mwasm.WasmContractInfoMsgType:
		convertMsg, _ := xplac.Msg.(wasmtypes.QueryContractInfoRequest)
		res, err = queryClient.ContractInfo(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm contract state all
	case xplac.MsgType == mwasm.WasmContractStateAllMsgType:
		convertMsg, _ := xplac.Msg.(wasmtypes.QueryAllContractStateRequest)
		res, err = queryClient.AllContractState(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm contract history
	case xplac.MsgType == mwasm.WasmContractHistoryMsgType:
		convertMsg, _ := xplac.Msg.(wasmtypes.QueryContractHistoryRequest)
		res, err = queryClient.ContractHistory(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm pinned
	case xplac.MsgType == mwasm.WasmPinnedMsgType:
		convertMsg, _ := xplac.Msg.(wasmtypes.QueryPinnedCodesRequest)
		res, err = queryClient.PinnedCodes(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm libwasmvm version
	case xplac.MsgType == mwasm.WasmLibwasmvmVersionMsgType:
		convertMsg, _ := xplac.Msg.(string)
		return convertMsg, nil

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err = printProto(xplac, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
