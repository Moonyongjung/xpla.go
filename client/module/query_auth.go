package module

import (
	mauth "github.com/Moonyongjung/xpla.go/core/auth"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Query client for auth module.
func (i IXplaClient) QueryAuth() (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcAuth(i)
	} else {
		return queryByLcdAuth(i)
	}
}

func queryByGrpcAuth(i IXplaClient) (string, error) {
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
			return "", util.LogErr("error: need RPC URL when txs methods")
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

const (
	authParamsLabel   = "params"
	authAccountsLabel = "accounts"
	authTxsLabel      = "txs"
)

func queryByLcdAuth(i IXplaClient) (string, error) {

	url := util.MakeQueryLcdUrl(authv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Auth params
	case i.Ixplac.GetMsgType() == mauth.AuthQueryParamsMsgType:
		url = url + authParamsLabel

	// Auth account
	case i.Ixplac.GetMsgType() == mauth.AuthQueryAccAddressMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(authtypes.QueryAccountRequest)
		url = url + util.MakeQueryLabels(authAccountsLabel, convertMsg.Address)

	// Auth accounts
	case i.Ixplac.GetMsgType() == mauth.AuthQueryAccountsMsgType:
		url = url + authAccountsLabel

	// Auth tx by event
	case i.Ixplac.GetMsgType() == mauth.AuthQueryTxsByEventsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().([]string)
		parsedEvent := convertMsg[0]
		parsedPage := convertMsg[1]
		parsedLimit := convertMsg[2]

		events := "?events=" + parsedEvent
		page := "&pagination.page=" + parsedPage
		limit := "&pagination.limit=" + parsedLimit

		url = "/cosmos/tx/v1beta1/"
		url = url + authTxsLabel + events + page + limit

	// Auth tx
	case i.Ixplac.GetMsgType() == mauth.AuthQueryTxMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().([]string)
		parsedValue := convertMsg[0]
		parsedTxType := convertMsg[1]

		url = "/cosmos/tx/v1beta1/"
		if parsedTxType == "hash" {
			url = url + util.MakeQueryLabels(authTxsLabel, parsedValue)

		} else if parsedTxType == "signature" {
			// inactivate
			return "", util.LogErr("inactivate GetTxEvent('signature') when using LCD because of sometimes generating parsing error that based64 encoded signature has '='")
			// events := "?events=" + parsedValue
			// page := "&pagination.page=" + util.FromIntToString(rest.DefaultPage)
			// limit := "&pagination.limit=" + util.FromIntToString(rest.DefaultLimit)

			// url = url + authTxsLabel + events + page + limit
		} else {
			events := "?events=" + parsedValue
			page := "&pagination.page=" + util.FromIntToString(rest.DefaultPage)
			limit := "&pagination.limit=" + util.FromIntToString(rest.DefaultLimit)

			url = url + authTxsLabel + events + page + limit
		}

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil
}
