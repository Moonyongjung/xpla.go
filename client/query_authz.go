package client

import (
	mauthz "github.com/Moonyongjung/xpla.go/core/authz"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/x/authz"
)

// Query client for authz module.
func queryAuthz(xplac *XplaClient) (string, error) {
	queryClient := authz.NewQueryClient(xplac.Grpc)

	switch {
	// Authz grant
	case xplac.MsgType == mauthz.AuthzQueryGrantMsgType:
		convertMsg, _ := xplac.Msg.(authz.QueryGrantsRequest)
		res, err = queryClient.Grants(
			xplac.Context,
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Authz grant by grantee
	case xplac.MsgType == mauthz.AuthzQueryGrantsByGranteeMsgType:
		convertMsg, _ := xplac.Msg.(authz.QueryGranteeGrantsRequest)
		res, err = queryClient.GranteeGrants(
			xplac.Context,
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Authz grant by granter
	case xplac.MsgType == mauthz.AuthzQueryGrantsByGranterMsgType:
		convertMsg, _ := xplac.Msg.(authz.QueryGranterGrantsRequest)
		res, err = queryClient.GranterGrants(
			xplac.Context,
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
