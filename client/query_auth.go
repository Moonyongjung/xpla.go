package client

import (
	"context"

	mauth "github.com/Moonyongjung/xpla.go/core/auth"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/types/rest"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Query client for auth module.
func queryAuth(xplac *XplaClient) (string, error) {
	queryClient := authtypes.NewQueryClient(xplac.Grpc)

	switch {
	// Auth params
	case xplac.MsgType == mauth.AuthQueryParamsMsgType:
		convertMsg, _ := xplac.Msg.(authtypes.QueryParamsRequest)
		res, err = queryClient.Params(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Auth account
	case xplac.MsgType == mauth.AuthQueryAccAddressMsgType:
		convertMsg, _ := xplac.Msg.(authtypes.QueryAccountRequest)
		res, err = queryClient.Account(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Auth accounts
	case xplac.MsgType == mauth.AuthQueryAccountsMsgType:
		convertMsg, _ := xplac.Msg.(authtypes.QueryAccountsRequest)
		res, err = queryClient.Accounts(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Auth tx by event
	case xplac.MsgType == mauth.AuthQueryTxsByEventsMsgType:
		if xplac.Opts.RpcURL == "" {
			return "", util.LogErr("Error: Need RPC URL when txs methods")
		}
		convertMsg, _ := xplac.Msg.([]string)
		msgLength := len(convertMsg)
		tmEvents := convertMsg[:msgLength-2]
		page := util.FromStringToInt(convertMsg[msgLength-2])
		limit := util.FromStringToInt(convertMsg[msgLength-1])
		clientCtx, err := clientForQuery(xplac)
		if err != nil {
			return "", err
		}

		res, err = authtx.QueryTxsByEvents(clientCtx, tmEvents, page, limit, "")
		if err != nil {
			return "", err
		}

	// Auth tx
	case xplac.MsgType == mauth.AuthQueryTxMsgType:
		if xplac.Opts.RpcURL == "" {
			return "", util.LogErr("Error: Need RPC URL when txs methods")
		}
		convertMsg, _ := xplac.Msg.([]string)
		msgLength := len(convertMsg)
		tmEvents := convertMsg[:msgLength-1]
		txType := convertMsg[msgLength-1]

		clientCtx, err := clientForQuery(xplac)
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

	out, err = printProto(xplac, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
