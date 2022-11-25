package wasm

import (
	"encoding/base64"
	"encoding/hex"
	"errors"

	wasm "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// (Tx) make msg - store code
func MakeStoreCodeMsg(storeMsg types.StoreMsg, addr sdk.AccAddress) (wasm.MsgStoreCode, error) {
	msg, err := parseStoreCodeArgs(storeMsg, addr)
	if err != nil {
		return wasm.MsgStoreCode{}, util.LogErr(err)
	}

	if err = msg.ValidateBasic(); err != nil {
		return wasm.MsgStoreCode{}, util.LogErr(err)
	}

	return msg, nil
}

// (Tx) make msg - instantiate
func MakeInstantiateMsg(instantiateMsg types.InstantiateMsg, addr sdk.AccAddress) (wasm.MsgInstantiateContract, error) {
	if (types.InstantiateMsg{}) == instantiateMsg {
		return wasm.MsgInstantiateContract{}, util.LogErr("Empty request or type of parameter is not correct")
	}

	if instantiateMsg.CodeId == "" ||
		instantiateMsg.Amount == "" ||
		instantiateMsg.Label == "" ||
		instantiateMsg.InitMsg == "" {
		return wasm.MsgInstantiateContract{}, util.LogErr("Empty mandatory parameters")
	}

	msg, err := parseInstantiateArgs(instantiateMsg, addr)
	if err != nil {
		return wasm.MsgInstantiateContract{}, util.LogErr(err)
	}

	if err = msg.ValidateBasic(); err != nil {
		return wasm.MsgInstantiateContract{}, util.LogErr(err)
	}

	return msg, nil
}

// (Tx) make msg - execute
func MakeExecuteMsg(executeMsg types.ExecuteMsg, addr sdk.AccAddress) (wasm.MsgExecuteContract, error) {
	if (types.ExecuteMsg{}) == executeMsg {
		return wasm.MsgExecuteContract{}, util.LogErr("Empty request or type of parameter is not correct")
	}

	msg, err := parseExecuteArgs(executeMsg, addr)
	if err != nil {
		util.LogErr(err)
		return wasm.MsgExecuteContract{}, err
	}

	if err = msg.ValidateBasic(); err != nil {
		util.LogErr(err)
		return wasm.MsgExecuteContract{}, err
	}

	return msg, nil
}

// (Tx) make msg - clear contract admin
func MakeClearContractAdminMsg(clearContractAdminMsg types.ClearContractAdminMsg, privKey key.PrivateKey) (wasm.MsgClearAdmin, error) {
	msg, err := parseClearContractAdminArgs(clearContractAdminMsg, privKey)
	if err != nil {
		return wasm.MsgClearAdmin{}, err
	}

	return msg, nil
}

// (Tx) make msg - set contract admin
func MakeSetContractAdmintMsg(setContractAdminMsg types.SetContractAdminMsg, privKey key.PrivateKey) (wasm.MsgUpdateAdmin, error) {
	msg, err := parseSetContractAdminArgs(setContractAdminMsg, privKey)
	if err != nil {
		return wasm.MsgUpdateAdmin{}, err
	}

	return msg, nil
}

// (Tx) make msg - migrate
func MakeMigrateMsg(migrateMsg types.MigrateMsg, privKey key.PrivateKey) (wasm.MsgMigrateContract, error) {
	msg, err := parseMigrateArgs(migrateMsg, privKey)
	if err != nil {
		return wasm.MsgMigrateContract{}, err
	}

	return msg, nil
}

// (Query) make msg - query contract
func MakeQueryMsg(queryMsg types.QueryMsg, addr sdk.AccAddress) (wasm.QuerySmartContractStateRequest, error) {
	if (types.QueryMsg{}) == queryMsg {
		return wasm.QuerySmartContractStateRequest{}, util.LogErr("Empty request or type of parameter is not correct")
	}

	msg, err := parseQueryArgs(queryMsg, addr)
	if err != nil {
		util.LogErr(err)
		return wasm.QuerySmartContractStateRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - list code
func MakeListcodeMsg() (wasm.QueryCodesRequest, error) {
	msg := parseListcodeArgs()
	return msg, nil
}

// (Query) make msg - list contract by code
func MakeListContractByCodeMsg(listContractByCodeMsg types.ListContractByCodeMsg) (wasm.QueryContractsByCodeRequest, error) {
	if (types.ListContractByCodeMsg{}) == listContractByCodeMsg {
		return wasm.QueryContractsByCodeRequest{}, util.LogErr("Empty request or type of parameter is not correct")
	}
	msg := parseListContractByCodeArgs(listContractByCodeMsg)
	return msg, nil
}

// (Query) make msg - download
func MakeDownloadMsg(downloadMsg types.DownloadMsg) ([]interface{}, error) {
	var msgInterfaceSlice []interface{}
	if (types.DownloadMsg{}) == downloadMsg {
		return nil, util.LogErr("Empty request or type of parameter is not correct")
	}
	msg := parseDownloadArgs(downloadMsg)
	msgInterfaceSlice = append(msgInterfaceSlice, msg)
	msgInterfaceSlice = append(msgInterfaceSlice, downloadMsg.DownloadFileName)
	return msgInterfaceSlice, nil
}

// (Query) make msg - code info
func MakeCodeInfoMsg(codeInfoMsg types.CodeInfoMsg) (wasm.QueryCodeRequest, error) {
	if (types.CodeInfoMsg{}) == codeInfoMsg {
		return wasm.QueryCodeRequest{}, util.LogErr("Empty request or type of parameter is not correct")
	}
	msg := parseCodeInfoArgs(codeInfoMsg)
	return msg, nil
}

// (Query) make msg - contract info
func MakeContractInfoMsg(contractInfoMsg types.ContractInfoMsg) (wasm.QueryContractInfoRequest, error) {
	if (types.ContractInfoMsg{}) == contractInfoMsg {
		return wasm.QueryContractInfoRequest{}, util.LogErr("Empty request or type of parameter is not correct")
	}
	msg := parseContractInfoArgs(contractInfoMsg)
	return msg, nil
}

// (Query) make msg - contract state all
func MakeContractStateAllMsg(contractStateAllMsg types.ContractStateAllMsg) (wasm.QueryAllContractStateRequest, error) {
	if (types.ContractStateAllMsg{}) == contractStateAllMsg {
		return wasm.QueryAllContractStateRequest{}, util.LogErr("Empty request or type of parameter is not correct")
	}
	msg := parseContractStateAllArgs(contractStateAllMsg)
	return msg, nil
}

// (Query) make msg - history
func MakeContractHistoryMsg(contractHistoryMsg types.ContractHistoryMsg) (wasm.QueryContractHistoryRequest, error) {
	if (types.ContractHistoryMsg{}) == contractHistoryMsg {
		return wasm.QueryContractHistoryRequest{}, util.LogErr("Empty request or type of parameter is not correct")
	}
	msg := parseContractHistoryArgs(contractHistoryMsg)
	return msg, nil
}

// (Query) make msg - pinned
func MakePinnedMsg() (wasm.QueryPinnedCodesRequest, error) {
	msg := parsePinnedArgs()
	return msg, nil
}

// (Query) make msg - libwasmvm version
func MakeLibwasmvmVersionMsg() (string, error) {
	msg, err := parseLibwasmvmVersionArgs()
	if err != nil {
		return "", err
	}

	return msg, nil
}

type ArgumentDecoder struct {
	// dec is the default decoder
	dec                func(string) ([]byte, error)
	asciiF, hexF, b64F bool
}

// Make new query decoder.
func NewArgDecoder(def func(string) ([]byte, error)) *ArgumentDecoder {
	return &ArgumentDecoder{dec: def}
}

func (a *ArgumentDecoder) DecodeString(s string) ([]byte, error) {
	found := -1
	for i, v := range []*bool{&a.asciiF, &a.hexF, &a.b64F} {
		if !*v {
			continue
		}
		if found != -1 {

			return nil, errors.New("multiple decoding flags used")
		}
		found = i
	}
	switch found {
	case 0:
		return AsciiDecodeString(s)
	case 1:
		return hex.DecodeString(s)
	case 2:
		return base64.StdEncoding.DecodeString(s)
	default:
		return a.dec(s)
	}
}

func AsciiDecodeString(s string) ([]byte, error) {
	return []byte(s), nil
}
