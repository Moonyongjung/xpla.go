package module

import (
	"io/ioutil"
	"strings"

	mwasm "github.com/Moonyongjung/xpla.go/core/wasm"
	"github.com/Moonyongjung/xpla.go/util"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// Query client for wasm module.
func (i IXplaClient) QueryWasm() (string, error) {
	queryClient := wasmtypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Wasm query contract
	case i.Ixplac.GetMsgType() == mwasm.WasmQueryContractMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QuerySmartContractStateRequest)
		res, err = queryClient.SmartContractState(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm list code
	case i.Ixplac.GetMsgType() == mwasm.WasmListCodeMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryCodesRequest)
		res, err = queryClient.Codes(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm list contract by code
	case i.Ixplac.GetMsgType() == mwasm.WasmListContractByCodeMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryContractsByCodeRequest)
		res, err = queryClient.ContractsByCode(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm download
	case i.Ixplac.GetMsgType() == mwasm.WasmDownloadMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().([]interface{})[0].(wasmtypes.QueryCodeRequest)
		downloadFileName, _ := i.Ixplac.GetMsg().([]interface{})[1].(string)
		if !strings.Contains(downloadFileName, ".wasm") {
			downloadFileName = downloadFileName + ".wasm"
		}
		res, err := queryClient.Code(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}
		ioutil.WriteFile(downloadFileName, res.Data, 0o600)

	// Wasm code info
	case i.Ixplac.GetMsgType() == mwasm.WasmCodeInfoMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryCodeRequest)
		res, err = queryClient.Code(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm contract info
	case i.Ixplac.GetMsgType() == mwasm.WasmContractInfoMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryContractInfoRequest)
		res, err = queryClient.ContractInfo(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm contract state all
	case i.Ixplac.GetMsgType() == mwasm.WasmContractStateAllMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryAllContractStateRequest)
		res, err = queryClient.AllContractState(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm contract history
	case i.Ixplac.GetMsgType() == mwasm.WasmContractHistoryMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryContractHistoryRequest)
		res, err = queryClient.ContractHistory(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm pinned
	case i.Ixplac.GetMsgType() == mwasm.WasmPinnedMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryPinnedCodesRequest)
		res, err = queryClient.PinnedCodes(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Wasm libwasmvm version
	case i.Ixplac.GetMsgType() == mwasm.WasmLibwasmvmVersionMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(string)
		return convertMsg, nil

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err = printProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
