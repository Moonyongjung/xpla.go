package auth

import (
	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// (Query) make msg - auth param
func MakeAuthParamMsg() (authtypes.QueryParamsRequest, error) {
	return authtypes.QueryParamsRequest{}, nil
}

// (Query) make msg - auth account
func MakeQueryAccAddressMsg(queryAccAddressMsg types.QueryAccAddressMsg) (authtypes.QueryAccountRequest, error) {
	if (types.QueryAccAddressMsg{}) == queryAccAddressMsg {
		return authtypes.QueryAccountRequest{}, util.LogErr("Empty request or type of parameter is not correct")
	}

	return authtypes.QueryAccountRequest{
		Address: queryAccAddressMsg.Address,
	}, nil
}

// (Query) make msg - auth accounts
func MakeQueryAccountsMsg() (authtypes.QueryAccountsRequest, error) {
	return authtypes.QueryAccountsRequest{
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - transactions by events
func MakeTxsByEventsMsg(txsByEventsMsg types.QueryTxsByEventsMsg) (QueryTxsByEventParseMsg, error) {
	if (types.QueryTxsByEventsMsg{}) == txsByEventsMsg {
		return QueryTxsByEventParseMsg{}, util.LogErr("Empty request or type of parameter is not correct")
	}

	msg, err := parseTxsByEventsArgs(txsByEventsMsg)
	if err != nil {
		return QueryTxsByEventParseMsg{}, nil
	}

	return msg, nil
}

// (Query) make msg - transaction
func MakeQueryTxMsg(queryTxMsg types.QueryTxMsg) (QueryTxParseMsg, error) {
	msg, err := parseQueryTxArgs(queryTxMsg)
	if err != nil {
		return QueryTxParseMsg{}, err
	}

	return msg, nil
}
