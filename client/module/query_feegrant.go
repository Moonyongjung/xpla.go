package module

import (
	mfeegrant "github.com/Moonyongjung/xpla.go/core/feegrant"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// Query client for fee-grant module.
func (i IXplaClient) QueryFeegrant() (string, error) {
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
			return "", err
		}

	// Feegrant grants by grantee
	case i.Ixplac.GetMsgType() == mfeegrant.FeegrantQueryGrantsByGranteeMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(feegrant.QueryAllowancesRequest)
		res, err = queryClient.Allowances(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Feegrant grants by granter
	case i.Ixplac.GetMsgType() == mfeegrant.FeegrantQueryGrantsByGranterMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(feegrant.QueryAllowancesByGranterRequest)
		res, err = queryClient.AllowancesByGranter(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
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
