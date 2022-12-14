package module

import (
	mfeegrant "github.com/Moonyongjung/xpla.go/core/feegrant"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	feegrantv1beta1 "cosmossdk.io/api/cosmos/feegrant/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// Query client for fee-grant module.
func (i IXplaClient) QueryFeegrant() (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcFeegrant(i)
	} else {
		return queryByLcdFeegrant(i)
	}

}

func queryByGrpcFeegrant(i IXplaClient) (string, error) {
	queryClient := feegrant.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Feegrant state
	case i.Ixplac.GetMsgType() == mfeegrant.FeegrantQueryGrantMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(feegrant.QueryAllowanceRequest)
		res, err = queryClient.Allowance(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Feegrant grants by grantee
	case i.Ixplac.GetMsgType() == mfeegrant.FeegrantQueryGrantsByGranteeMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(feegrant.QueryAllowancesRequest)
		res, err = queryClient.Allowances(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Feegrant grants by granter
	case i.Ixplac.GetMsgType() == mfeegrant.FeegrantQueryGrantsByGranterMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(feegrant.QueryAllowancesByGranterRequest)
		res, err = queryClient.AllowancesByGranter(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
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
	feegrantAllowanceLabel  = "allowance"
	feegrantAllowancesLabel = "allowances"
)

func queryByLcdFeegrant(i IXplaClient) (string, error) {
	url := util.MakeQueryLcdUrl(feegrantv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Feegrant state
	case i.Ixplac.GetMsgType() == mfeegrant.FeegrantQueryGrantMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(feegrant.QueryAllowanceRequest)

		url = url + util.MakeQueryLabels(feegrantAllowanceLabel, convertMsg.Granter, convertMsg.Grantee)

	// Feegrant grants by grantee
	case i.Ixplac.GetMsgType() == mfeegrant.FeegrantQueryGrantsByGranteeMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(feegrant.QueryAllowancesRequest)

		url = url + util.MakeQueryLabels(feegrantAllowancesLabel, convertMsg.Grantee)

	// Feegrant grants by granter
	case i.Ixplac.GetMsgType() == mfeegrant.FeegrantQueryGrantsByGranterMsgType:
		return "", util.LogErr(errors.ErrNotSupport, "unsupported querying feegrant state(grants by granter) by using LCD")

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil

}
