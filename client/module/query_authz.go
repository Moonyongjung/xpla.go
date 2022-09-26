package module

import (
	mauthz "github.com/Moonyongjung/xpla.go/core/authz"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/x/authz"
)

// Query client for authz module.
func (i IXplaClient) QueryAuthz() (string, error) {
	queryClient := authz.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Authz grant
	case i.Ixplac.GetMsgType() == mauthz.AuthzQueryGrantMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(authz.QueryGrantsRequest)
		res, err = queryClient.Grants(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Authz grant by grantee
	case i.Ixplac.GetMsgType() == mauthz.AuthzQueryGrantsByGranteeMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(authz.QueryGranteeGrantsRequest)
		res, err = queryClient.GranteeGrants(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Authz grant by granter
	case i.Ixplac.GetMsgType() == mauthz.AuthzQueryGrantsByGranterMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(authz.QueryGranterGrantsRequest)
		res, err = queryClient.GranterGrants(
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
