package queries

import (
	"encoding/base64"
	"os"
	"strings"

	mwasm "github.com/Moonyongjung/xpla.go/core/wasm"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// Query client for wasm module.
func (i IXplaClient) QueryWasm() (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcWasm(i)
	} else {
		return queryByLcdWasm(i)
	}
}

func queryByGrpcWasm(i IXplaClient) (string, error) {
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
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Wasm list code
	case i.Ixplac.GetMsgType() == mwasm.WasmListCodeMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryCodesRequest)
		res, err = queryClient.Codes(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Wasm list contract by code
	case i.Ixplac.GetMsgType() == mwasm.WasmListContractByCodeMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryContractsByCodeRequest)
		res, err = queryClient.ContractsByCode(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
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
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}
		os.WriteFile(downloadFileName, res.Data, 0o600)
		return "download complete", nil

	// Wasm code info
	case i.Ixplac.GetMsgType() == mwasm.WasmCodeInfoMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryCodeRequest)
		res, err = queryClient.Code(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Wasm contract info
	case i.Ixplac.GetMsgType() == mwasm.WasmContractInfoMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryContractInfoRequest)
		res, err = queryClient.ContractInfo(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Wasm contract state all
	case i.Ixplac.GetMsgType() == mwasm.WasmContractStateAllMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryAllContractStateRequest)
		res, err = queryClient.AllContractState(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Wasm contract history
	case i.Ixplac.GetMsgType() == mwasm.WasmContractHistoryMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryContractHistoryRequest)
		res, err = queryClient.ContractHistory(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Wasm pinned
	case i.Ixplac.GetMsgType() == mwasm.WasmPinnedMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryPinnedCodesRequest)
		res, err = queryClient.PinnedCodes(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Wasm libwasmvm version
	case i.Ixplac.GetMsgType() == mwasm.WasmLibwasmvmVersionMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(string)
		return convertMsg, nil

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err = printProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

const (
	wasmContractLabel = "contract"
	wasmSmartLabel    = "smart"
	wasmCodeLabel     = "code"
	wasmCodesLabel    = "codes"
	wasmStateLabel    = "state"
	wasmHistoryLabel  = "history"
	wasmPinnedLabel   = "pinned"
)

func queryByLcdWasm(i IXplaClient) (string, error) {
	url := "/cosmwasm/wasm/v1/"

	switch {
	// Wasm query contract
	case i.Ixplac.GetMsgType() == mwasm.WasmQueryContractMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QuerySmartContractStateRequest)
		based64EncodedData := base64.StdEncoding.EncodeToString([]byte(convertMsg.QueryData))

		url = url + util.MakeQueryLabels(wasmContractLabel, convertMsg.Address, wasmSmartLabel, based64EncodedData)

	// Wasm list code
	case i.Ixplac.GetMsgType() == mwasm.WasmListCodeMsgType:
		url = url + wasmCodeLabel

	// Wasm list contract by code
	case i.Ixplac.GetMsgType() == mwasm.WasmListContractByCodeMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryContractsByCodeRequest)

		url = url + util.MakeQueryLabels(wasmCodeLabel, util.FromUint64ToString(convertMsg.CodeId))

	// Wasm download
	case i.Ixplac.GetMsgType() == mwasm.WasmDownloadMsgType:
		return "", util.LogErr(errors.ErrNotSupport, "unsupported download wasm file by using LCD. query delegations of a delegator")

	// Wasm code info
	case i.Ixplac.GetMsgType() == mwasm.WasmCodeInfoMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryCodeRequest)

		url = url + util.MakeQueryLabels(wasmCodeLabel, util.FromUint64ToString(convertMsg.CodeId))

	// Wasm contract info
	case i.Ixplac.GetMsgType() == mwasm.WasmContractInfoMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryContractInfoRequest)

		url = url + util.MakeQueryLabels(wasmContractLabel, convertMsg.Address)

	// Wasm contract state all
	case i.Ixplac.GetMsgType() == mwasm.WasmContractStateAllMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryAllContractStateRequest)

		url = url + util.MakeQueryLabels(wasmContractLabel, convertMsg.Address, wasmStateLabel)

	// Wasm contract history
	case i.Ixplac.GetMsgType() == mwasm.WasmContractHistoryMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(wasmtypes.QueryContractHistoryRequest)

		url = url + util.MakeQueryLabels(wasmContractLabel, convertMsg.Address, wasmHistoryLabel)

	// Wasm pinned
	case i.Ixplac.GetMsgType() == mwasm.WasmPinnedMsgType:

		url = url + util.MakeQueryLabels(wasmCodesLabel, wasmPinnedLabel)

	// Wasm libwasmvm version
	case i.Ixplac.GetMsgType() == mwasm.WasmLibwasmvmVersionMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(string)
		return convertMsg, nil

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil

}
