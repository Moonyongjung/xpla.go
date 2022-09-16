package client

import (
	"context"

	mfeegrant "github.com/Moonyongjung/xpla.go/core/feegrant"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// Query client for fee-grant module.
func queryFeegrant(xplac *XplaClient) (string, error) {
	queryClient := feegrant.NewQueryClient(xplac.Grpc)

	switch {
	// Feegrant state
	case xplac.MsgType == mfeegrant.FeegrantQueryGrantMsgType:
		convertMsg, _ := xplac.Msg.(feegrant.QueryAllowanceRequest)
		res, err = queryClient.Allowance(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Feegrant grants by grantee
	case xplac.MsgType == mfeegrant.FeegrantQueryGrantsByGranteeMsgType:
		convertMsg, _ := xplac.Msg.(feegrant.QueryAllowancesRequest)
		res, err = queryClient.Allowances(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Feegrant grants by granter
	case xplac.MsgType == mfeegrant.FeegrantQueryGrantsByGranterMsgType:
		convertMsg, _ := xplac.Msg.(feegrant.QueryAllowancesByGranterRequest)
		res, err = queryClient.AllowancesByGranter(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
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
