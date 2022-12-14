package feegrant

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// (Tx) make msg - fee grant
func MakeFeeGrantMsg(feeGrantMsg types.FeeGrantMsg, privKey key.PrivateKey) (feegrant.MsgGrantAllowance, error) {
	msg, err := parseFeeGrantArgs(feeGrantMsg, privKey)
	if err != nil {
		return feegrant.MsgGrantAllowance{}, err
	}

	return msg, nil
}

// (Tx) make msg - fee grant revoke
func MakeRevokeFeeGrantMsg(revokeFeeGrantMsg types.RevokeFeeGrantMsg, privKey key.PrivateKey) (feegrant.MsgRevokeAllowance, error) {
	msg, err := parseRevokeFeeGrantArgs(revokeFeeGrantMsg, privKey)
	if err != nil {
		return feegrant.MsgRevokeAllowance{}, err
	}

	return msg, nil
}

// (Query) make msg - query fee grants
func MakeQueryFeeGrantMsg(queryFeeGrantMsg types.QueryFeeGrantMsg) (feegrant.QueryAllowanceRequest, error) {
	msg, err := parseQueryFeeGrantArgs(queryFeeGrantMsg)
	if err != nil {
		return feegrant.QueryAllowanceRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - fee grants by grantee
func MakeQueryFeeGrantsByGranteeMsg(queryFeeGrantMsg types.QueryFeeGrantMsg) (feegrant.QueryAllowancesRequest, error) {
	msg, err := parseQueryFeeGrantsByGranteeArgs(queryFeeGrantMsg)
	if err != nil {
		return feegrant.QueryAllowancesRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - fee grants by granter
func MakeQueryFeeGrantsByGranterMsg(queryFeeGrantMsg types.QueryFeeGrantMsg) (feegrant.QueryAllowancesByGranterRequest, error) {
	msg, err := parseQueryFeeGrantsByGranterArgs(queryFeeGrantMsg)
	if err != nil {
		return feegrant.QueryAllowancesByGranterRequest{}, err
	}

	return msg, nil
}
