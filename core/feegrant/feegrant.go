package feegrant

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// (Tx) make msg - fee grant
func MakeGrantMsg(grantMsg types.GrantMsg, privKey key.PrivateKey) (*feegrant.MsgGrantAllowance, error) {
	msg, err := parseGrantArgs(grantMsg, privKey)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// (Tx) make msg - fee grant revoke
func MakeRevokeGrantMsg(revokeGrantMsg types.RevokeGrantMsg, privKey key.PrivateKey) (feegrant.MsgRevokeAllowance, error) {
	msg, err := parseRevokeGrantArgs(revokeGrantMsg, privKey)
	if err != nil {
		return feegrant.MsgRevokeAllowance{}, err
	}

	return msg, nil
}

// (Query) make msg - query fee grants
func MakeQueryGrantMsg(queryGrantMsg types.QueryGrantMsg) (feegrant.QueryAllowanceRequest, error) {
	msg, err := parseQueryGrantArgs(queryGrantMsg)
	if err != nil {
		return feegrant.QueryAllowanceRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - fee grants by grantee
func MakeQueryGrantsByGranteeMsg(queryGrantMsg types.QueryGrantMsg) (feegrant.QueryAllowancesRequest, error) {
	msg, err := parseQueryGrantsByGranteeArgs(queryGrantMsg)
	if err != nil {
		return feegrant.QueryAllowancesRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - fee grants by granter
func MakeQueryGrantsByGranterMsg(queryGrantMsg types.QueryGrantMsg) (feegrant.QueryAllowancesByGranterRequest, error) {
	msg, err := parseQueryGrantsByGranterArgs(queryGrantMsg)
	if err != nil {
		return feegrant.QueryAllowancesByGranterRequest{}, err
	}

	return msg, nil
}
