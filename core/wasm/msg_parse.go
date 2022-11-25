package wasm

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/CosmWasm/wasmd/x/wasm/ioutils"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvm "github.com/CosmWasm/wasmvm"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	instantiateByEverybody = "instantiate-everybody"
	instantiateNobody      = "instantiate-nobody"
	instantiateBySender    = "instantiate-only-sender"
	instantiateByAddress   = "instantiate-only-address"
)

// Parsing - store code
func parseStoreCodeArgs(storeMsg types.StoreMsg, sender sdk.AccAddress) (wasmtypes.MsgStoreCode, error) {
	if storeMsg.FilePath == "" {
		return wasmtypes.MsgStoreCode{}, util.LogErr("filepath is empty")
	}

	wasm, err := ioutil.ReadFile(storeMsg.FilePath)
	if err != nil {
		return wasmtypes.MsgStoreCode{}, err
	}

	// gzip the wasm file
	if ioutils.IsWasm(wasm) {
		wasm, err = ioutils.GzipIt(wasm)

		if err != nil {
			return wasmtypes.MsgStoreCode{}, err
		}
	} else if !ioutils.IsGzip(wasm) {
		return wasmtypes.MsgStoreCode{}, util.LogErr("invalid input file. Use wasm binary or gzip")
	}

	permission, err := instantiatePermission(storeMsg.InstantiatePermission, sender)
	if err != nil {
		return wasmtypes.MsgStoreCode{}, err
	}

	msg := wasmtypes.MsgStoreCode{
		Sender:                sender.String(),
		WASMByteCode:          wasm,
		InstantiatePermission: permission,
	}
	return msg, nil
}

func instantiatePermission(permission string, sender sdk.AccAddress) (*wasmtypes.AccessConfig, error) {
	var permMethod string
	var onlyAddr string

	if strings.Contains(permission, ".") {
		perm := strings.Split(permission, ".")
		permMethod = perm[0]
		onlyAddr = perm[1]
	} else {
		permMethod = permission
		onlyAddr = ""
	}

	switch {
	case permMethod == "" || permMethod == instantiateByEverybody:
		return &wasmtypes.AllowEverybody, nil

	case permMethod == instantiateBySender:
		x := wasmtypes.AccessTypeOnlyAddress.With(sender)
		return &x, nil

	case permMethod == instantiateByAddress:
		if onlyAddr == "" {
			return nil, util.LogErr("invalid permission, empty address")
		}
		addr, err := sdk.AccAddressFromBech32(onlyAddr)
		if err != nil {
			return nil, err
		}
		x := wasmtypes.AccessTypeOnlyAddress.With(addr)
		return &x, nil

	case permMethod == instantiateNobody:
		return &wasmtypes.AllowNobody, nil

	default:
		return nil, util.LogErr("invalid permission type")
	}
}

// Parsing - instantiate
func parseInstantiateArgs(
	instantiateMsgData types.InstantiateMsg,
	sender sdk.AccAddress) (wasmtypes.MsgInstantiateContract, error) {

	rawCodeID := instantiateMsgData.CodeId
	if rawCodeID == "" {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr("No code ID")
	}

	// get the id of the code to instantiate
	codeID, err := strconv.ParseUint(rawCodeID, 10, 64)
	if err != nil {
		return wasmtypes.MsgInstantiateContract{}, err
	}

	amountStr := instantiateMsgData.Amount
	if amountStr == "" {
		amountStr = "0"
	}
	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(amountStr))
	if err != nil {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr("amount:", err)
	}

	label := instantiateMsgData.Label
	if label == "" {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr("label is required on all contracts")
	}

	initMsg := instantiateMsgData.InitMsg
	if initMsg == "" {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr("No Init Message")
	}

	adminStr := instantiateMsgData.Admin

	noAdminBool := true
	noAdminStr := instantiateMsgData.NoAdmin
	if noAdminStr == "true" {
		noAdminBool = true
	} else if noAdminStr == "" || noAdminStr == "false" {
		noAdminBool = false
	} else {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr("noAdmin parameter must set \"true\" or \"false\"")
	}

	// ensure sensible admin is set (or explicitly immutable)
	if adminStr == "" && !noAdminBool {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr("you must set an admin or explicitly pass --no-admin to make it immutible (wasmd issue #719)")
	}
	if adminStr != "" && noAdminBool {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr("you set an admin and passed --no-admin, those cannot both be true")
	}

	// build and sign the transaction, then broadcast to Tendermint
	msg := wasmtypes.MsgInstantiateContract{
		Sender: sender.String(),
		CodeID: codeID,
		Label:  label,
		Funds:  amount,
		Msg:    []byte(initMsg),
		Admin:  adminStr,
	}
	return msg, nil
}

// Parsing - execute
func parseExecuteArgs(executeMsgData types.ExecuteMsg,
	sender sdk.AccAddress) (wasmtypes.MsgExecuteContract, error) {
	amountStr := executeMsgData.Amount
	if amountStr == "" {
		amountStr = "0"
	}
	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(amountStr))
	if err != nil {
		return wasmtypes.MsgExecuteContract{}, util.LogErr("amount:", err)
	}

	return wasmtypes.MsgExecuteContract{
		Sender:   sender.String(),
		Contract: executeMsgData.ContractAddress,
		Funds:    amount,
		Msg:      []byte(executeMsgData.ExecMsg),
	}, nil
}

// Parsing - clear contract admin
func parseClearContractAdminArgs(clearContractAdminMsg types.ClearContractAdminMsg, privKey key.PrivateKey) (wasmtypes.MsgClearAdmin, error) {
	return wasmtypes.MsgClearAdmin{
		Sender:   util.GetAddrByPrivKey(privKey).String(),
		Contract: clearContractAdminMsg.ContractAddress,
	}, nil
}

// Parsing - set contract admin
func parseSetContractAdminArgs(setContractAdminMsg types.SetContractAdminMsg, privKey key.PrivateKey) (wasmtypes.MsgUpdateAdmin, error) {
	msg := wasmtypes.MsgUpdateAdmin{
		Sender:   util.GetAddrByPrivKey(privKey).String(),
		Contract: setContractAdminMsg.ContractAddress,
		NewAdmin: setContractAdminMsg.NewAdmin,
	}

	if err := msg.ValidateBasic(); err != nil {
		return wasmtypes.MsgUpdateAdmin{}, err
	}

	return msg, nil
}

// Parsing - migrate
func parseMigrateArgs(migrateMsg types.MigrateMsg, privKey key.PrivateKey) (wasmtypes.MsgMigrateContract, error) {
	return wasmtypes.MsgMigrateContract{
		Sender:   util.GetAddrByPrivKey(privKey).String(),
		Contract: migrateMsg.ContractAddress,
		CodeID:   util.FromStringToUint64(migrateMsg.CodeId),
		Msg:      []byte(migrateMsg.MigrateMsg),
	}, nil
}

// Parsing - query contract
func parseQueryArgs(queryMsgData types.QueryMsg,
	sender sdk.AccAddress) (wasmtypes.QuerySmartContractStateRequest, error) {
	decoder := NewArgDecoder(AsciiDecodeString)

	queryData, err := decoder.DecodeString(queryMsgData.QueryMsg)
	if err != nil {
		return wasmtypes.QuerySmartContractStateRequest{}, util.LogErr(err)
	}

	return wasmtypes.QuerySmartContractStateRequest{
		Address:   queryMsgData.ContractAddress,
		QueryData: queryData,
	}, nil
}

// Parsing - list code
func parseListcodeArgs() wasmtypes.QueryCodesRequest {
	return wasmtypes.QueryCodesRequest{
		Pagination: core.PageRequest,
	}
}

// Parsing - list contract by code
func parseListContractByCodeArgs(listContractByCodeMsgData types.ListContractByCodeMsg) wasmtypes.QueryContractsByCodeRequest {
	return wasmtypes.QueryContractsByCodeRequest{
		CodeId:     util.FromStringToUint64(listContractByCodeMsgData.CodeId),
		Pagination: core.PageRequest,
	}
}

// Parsing - download
func parseDownloadArgs(downloadMsgData types.DownloadMsg) wasmtypes.QueryCodeRequest {
	return wasmtypes.QueryCodeRequest{
		CodeId: util.FromStringToUint64(downloadMsgData.CodeId),
	}
}

// Parsing - code info
func parseCodeInfoArgs(codeInfoMsgData types.CodeInfoMsg) wasmtypes.QueryCodeRequest {
	return wasmtypes.QueryCodeRequest{
		CodeId: util.FromStringToUint64(codeInfoMsgData.CodeId),
	}
}

// Parsing - contract info
func parseContractInfoArgs(contractInfoMsgData types.ContractInfoMsg) wasmtypes.QueryContractInfoRequest {
	return wasmtypes.QueryContractInfoRequest{
		Address: contractInfoMsgData.ContractAddress,
	}
}

// Parsing - contract state all
func parseContractStateAllArgs(contractStateAllMsgData types.ContractStateAllMsg) wasmtypes.QueryAllContractStateRequest {
	return wasmtypes.QueryAllContractStateRequest{
		Address:    contractStateAllMsgData.ContractAddress,
		Pagination: core.PageRequest,
	}
}

// Parsing - history
func parseContractHistoryArgs(contractHistoryMsgData types.ContractHistoryMsg) wasmtypes.QueryContractHistoryRequest {
	return wasmtypes.QueryContractHistoryRequest{
		Address:    contractHistoryMsgData.ContractAddress,
		Pagination: core.PageRequest,
	}
}

// Parsing - pinned
func parsePinnedArgs() wasmtypes.QueryPinnedCodesRequest {
	return wasmtypes.QueryPinnedCodesRequest{
		Pagination: core.PageRequest,
	}
}

// Parsing - libwasmvm version
func parseLibwasmvmVersionArgs() (string, error) {
	version, err := wasmvm.LibwasmvmVersion()
	if err != nil {
		return "", err
	}
	return version, nil
}
