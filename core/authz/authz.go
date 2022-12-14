package authz

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/xpladev/xpla/app/params"
)

// (Tx) make msg - authz grant
func MakeAuthzGrantMsg(authzGrantMsg types.AuthzGrantMsg, privKey key.PrivateKey) (authz.MsgGrant, error) {
	msg, err := parseAuthzGrantArgs(authzGrantMsg, privKey)
	if err != nil {
		return authz.MsgGrant{}, err
	}

	return msg, nil
}

// (Tx) make msg - revoke
func MakeAuthzRevokeMsg(authzRevokeMsg types.AuthzRevokeMsg, privKey key.PrivateKey) (authz.MsgRevoke, error) {
	msg, err := parseAuthzRevokeArgs(authzRevokeMsg, privKey)
	if err != nil {
		return authz.MsgRevoke{}, err
	}

	return msg, nil
}

// (Tx) make msg - authz execute
func MakeAuthzExecMsg(authzExecMsg types.AuthzExecMsg, encodingConfig params.EncodingConfig) (authz.MsgExec, error) {
	msg, err := parseAuthzExecArgs(authzExecMsg, encodingConfig)
	if err != nil {
		return authz.MsgExec{}, err
	}

	return msg, nil
}

// (Query) make msg - authz grants
func MakeQueryAuthzGrantsMsg(queryAuthzGrantMsg types.QueryAuthzGrantMsg) (authz.QueryGrantsRequest, error) {
	msg, err := parseQueryAuthzGrantsArgs(queryAuthzGrantMsg)
	if err != nil {
		return authz.QueryGrantsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - authz grants by grantee
func MakeQueryAuthzGrantsByGranteeMsg(queryAuthzGrantMsg types.QueryAuthzGrantMsg) (authz.QueryGranteeGrantsRequest, error) {
	msg, err := parseQueryAuthzGrantsByGranteeArgs(queryAuthzGrantMsg)
	if err != nil {
		return authz.QueryGranteeGrantsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - authz grants by granter
func MakeQueryAuthzGrantsByGranterMsg(queryAuthzGrantMsg types.QueryAuthzGrantMsg) (authz.QueryGranterGrantsRequest, error) {
	msg, err := parseQueryAuthzGrantsByGranterArgs(queryAuthzGrantMsg)
	if err != nil {
		return authz.QueryGranterGrantsRequest{}, err
	}

	return msg, nil
}
