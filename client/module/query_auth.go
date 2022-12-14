package module

import (
	mauth "github.com/Moonyongjung/xpla.go/core/auth"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
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
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Auth account
	case i.Ixplac.GetMsgType() == mauth.AuthQueryAccAddressMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(authtypes.QueryAccountRequest)
		res, err = queryClient.Account(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Auth accounts
	case i.Ixplac.GetMsgType() == mauth.AuthQueryAccountsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(authtypes.QueryAccountsRequest)
		res, err = queryClient.Accounts(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Auth tx by event
	case i.Ixplac.GetMsgType() == mauth.AuthQueryTxsByEventsMsgType:
		if i.Ixplac.GetRpc() == "" {
			return "", util.LogErr(errors.ErrNotSatisfiedOptions, "query txs by events, need RPC URL when txs methods")
		}
		convertMsg, _ := i.Ixplac.GetMsg().(mauth.QueryTxsByEventParseMsg)
		clientCtx, err := clientForQuery(i)
		if err != nil {
			return "", err
		}

		res, err = authtx.QueryTxsByEvents(clientCtx, convertMsg.TmEvents, convertMsg.Page, convertMsg.Limit, "")
		if err != nil {
			return "", util.LogErr(errors.ErrRpcRequest, err)
		}

	// Auth tx
	case i.Ixplac.GetMsgType() == mauth.AuthQueryTxMsgType:
		if i.Ixplac.GetRpc() == "" {
			return "", util.LogErr(errors.ErrNotSatisfiedOptions, "auth query tx msg, need RPC URL when txs methods")
		}
		convertMsg, _ := i.Ixplac.GetMsg().(mauth.QueryTxParseMsg)

		clientCtx, err := clientForQuery(i)
		if err != nil {
			return "", err
		}

		if convertMsg.TxType == "hash" {
			res, err = authtx.QueryTx(clientCtx, convertMsg.TmEvents[0])
			if err != nil {
				return "", util.LogErr(errors.ErrRpcRequest, err)
			}
		} else {
			res, err = authtx.QueryTxsByEvents(clientCtx, convertMsg.TmEvents, rest.DefaultPage, rest.DefaultLimit, "")
			if err != nil {
				return "", util.LogErr(errors.ErrRpcRequest, err)
			}
		}

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
		convertMsg, _ := i.Ixplac.GetMsg().(mauth.QueryTxsByEventParseMsg)

		if len(convertMsg.TmEvents) > 1 {
			return "", util.LogErr(errors.ErrNotSupport, "support only one event on the LCD")
		}

		parsedEvent := convertMsg.TmEvents[0]
		parsedPage := convertMsg.Page
		parsedLimit := convertMsg.Limit

		events := "?events=" + parsedEvent
		page := "&pagination.page=" + util.FromIntToString(parsedPage)
		limit := "&pagination.limit=" + util.FromIntToString(parsedLimit)

		url = "/cosmos/tx/v1beta1/"
		url = url + authTxsLabel + events + page + limit

	// Auth tx
	case i.Ixplac.GetMsgType() == mauth.AuthQueryTxMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(mauth.QueryTxParseMsg)

		if len(convertMsg.TmEvents) > 1 {
			return "", util.LogErr(errors.ErrNotSupport, "support only one event on the LCD")
		}

		parsedValue := convertMsg.TmEvents
		parsedTxType := convertMsg.TxType

		url = "/cosmos/tx/v1beta1/"
		if parsedTxType == "hash" {
			url = url + util.MakeQueryLabels(authTxsLabel, parsedValue[0])

		} else if parsedTxType == "signature" {
			// inactivate
			return "", util.LogErr(errors.ErrNotSupport, "inactivate GetTxEvent('signature') when using LCD because of sometimes generating parsing error that based64 encoded signature has '='")
			// events := "?events=" + parsedValue
			// page := "&pagination.page=" + util.FromIntToString(rest.DefaultPage)
			// limit := "&pagination.limit=" + util.FromIntToString(rest.DefaultLimit)

			// url = url + authTxsLabel + events + page + limit
		} else {
			events := "?events=" + parsedValue[0]
			page := "&pagination.page=" + util.FromIntToString(rest.DefaultPage)
			limit := "&pagination.limit=" + util.FromIntToString(rest.DefaultLimit)

			url = url + authTxsLabel + events + page + limit
		}

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil
}
