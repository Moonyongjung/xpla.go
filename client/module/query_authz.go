package module

import (
	mauthz "github.com/Moonyongjung/xpla.go/core/authz"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	authzv1beta1 "cosmossdk.io/api/cosmos/authz/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

// Query client for authz module.
func (i IXplaClient) QueryAuthz() (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcAuthz(i)
	} else {
		return queryByLcdAuthz(i)
	}
}

func queryByGrpcAuthz(i IXplaClient) (string, error) {
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

const (
	authzGrantsLabel = "grants"
)

func queryByLcdAuthz(i IXplaClient) (string, error) {

	url := util.MakeQueryLcdUrl(authzv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Authz grant
	case i.Ixplac.GetMsgType() == mauthz.AuthzQueryGrantMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(authz.QueryGrantsRequest)
		parsedGranter := convertMsg.Granter
		parsedGrantee := convertMsg.Grantee

		granter := "?granter=" + parsedGranter
		grantee := "&grantee=" + parsedGrantee

		url = url + authzGrantsLabel + granter + grantee

	// Authz grant by grantee
	case i.Ixplac.GetMsgType() == mauthz.AuthzQueryGrantsByGranteeMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(authz.QueryGranteeGrantsRequest)
		grantee := convertMsg.Grantee

		url = url + util.MakeQueryLabels(authzGrantsLabel, "grantee", grantee)

	// Authz grant by granter
	case i.Ixplac.GetMsgType() == mauthz.AuthzQueryGrantsByGranterMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(authz.QueryGranterGrantsRequest)
		granter := convertMsg.Granter

		url = url + util.MakeQueryLabels(authzGrantsLabel, "granter", granter)

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil
}
