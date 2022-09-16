package auth

import (
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// (Query) make msg - auth param
func MakeAuthParamMsg() (auth.QueryParamsRequest, error) {
	msg, err := ParseAuthParamArgs()
	if err != nil {
		return auth.QueryParamsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - auth account
func MakeQueryAccAddressMsg(queryAccAddressMsg types.QueryAccAddressMsg) (auth.QueryAccountRequest, error) {
	if (types.QueryAccAddressMsg{}) == queryAccAddressMsg {
		return auth.QueryAccountRequest{}, util.LogErr("Empty request or type of parameter is not correct")
	}

	msg, err := parseQueryAccAddressArgs(queryAccAddressMsg)
	if err != nil {
		return auth.QueryAccountRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - auth accounts
func MakeQueryAccountsMsg() (auth.QueryAccountsRequest, error) {
	msg, err := parseQueryAccountsArgs()
	if err != nil {
		return auth.QueryAccountsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - transactions by events
func MakeTxsByEventsMsg(txsByEventsMsg types.QueryTxsByEventsMsg) ([]string, error) {
	if (types.QueryTxsByEventsMsg{}) == txsByEventsMsg {
		return []string{}, util.LogErr("Empty request or type of parameter is not correct")
	}

	msg, err := parseTxsByEventsArgs(txsByEventsMsg)
	if err != nil {
		return []string{}, nil
	}

	return msg, nil
}

// (Query) make msg - transaction
func MakeQueryTxMsg(queryTxMsg types.QueryTxMsg) ([]string, error) {
	msg, err := parseQueryTxArgs(queryTxMsg)
	if err != nil {
		return []string{}, err
	}

	return msg, nil
}
