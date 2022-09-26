package module

import (
	mauth "github.com/Moonyongjung/xpla.go/core/auth"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/types/rest"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Query client for auth module.
func (i IXplaClient) QueryAuth() (string, error) {
	queryClient := authtypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Auth params
	case i.Ixplac.GetMsgType() == mauth.AuthQueryParamsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(authtypes.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Auth account
	case i.Ixplac.GetMsgType() == mauth.AuthQueryAccAddressMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(authtypes.QueryAccountRequest)
		res, err = queryClient.Account(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Auth accounts
	case i.Ixplac.GetMsgType() == mauth.AuthQueryAccountsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(authtypes.QueryAccountsRequest)
		res, err = queryClient.Accounts(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Auth tx by event
	case i.Ixplac.GetMsgType() == mauth.AuthQueryTxsByEventsMsgType:
		if i.Ixplac.GetRpc() == "" {
			return "", util.LogErr("error: need RPC URL when txs methods")
		}
		convertMsg, _ := i.Ixplac.GetMsg().([]string)
		msgLength := len(convertMsg)
		tmEvents := convertMsg[:msgLength-2]
		page := util.FromStringToInt(convertMsg[msgLength-2])
		limit := util.FromStringToInt(convertMsg[msgLength-1])
		clientCtx, err := clientForQuery(i)
		if err != nil {
			return "", err
		}

		res, err = authtx.QueryTxsByEvents(clientCtx, tmEvents, page, limit, "")
		if err != nil {
			return "", err
		}

	// Auth tx
	case i.Ixplac.GetMsgType() == mauth.AuthQueryTxMsgType:
		if i.Ixplac.GetRpc() == "" {
			return "", util.LogErr("Error: Need RPC URL when txs methods")
		}
		convertMsg, _ := i.Ixplac.GetMsg().([]string)
		msgLength := len(convertMsg)
		tmEvents := convertMsg[:msgLength-1]
		txType := convertMsg[msgLength-1]

		clientCtx, err := clientForQuery(i)
		if err != nil {
			return "", err
		}

		if txType == "hash" {
			res, err = authtx.QueryTx(clientCtx, tmEvents[0])
			if err != nil {
				return "", err
			}
		} else {
			res, err = authtx.QueryTxsByEvents(clientCtx, tmEvents, rest.DefaultPage, rest.DefaultLimit, "")
			if err != nil {
				return "", err
			}
		}

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err = printProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
