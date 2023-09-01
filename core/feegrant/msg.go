package feegrant

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// (Tx) make msg - fee grant
func MakeFeeGrantMsg(feeGrantMsg types.FeeGrantMsg, privKey key.PrivateKey) (feegrant.MsgGrantAllowance, error) {
	return parseFeeGrantArgs(feeGrantMsg, privKey)
}

// (Tx) make msg - fee grant revoke
func MakeRevokeFeeGrantMsg(revokeFeeGrantMsg types.RevokeFeeGrantMsg, privKey key.PrivateKey) (feegrant.MsgRevokeAllowance, error) {
	return parseRevokeFeeGrantArgs(revokeFeeGrantMsg, privKey)
}

// (Query) make msg - query fee grants
func MakeQueryFeeGrantMsg(queryFeeGrantMsg types.QueryFeeGrantMsg) (feegrant.QueryAllowanceRequest, error) {
	return parseQueryFeeGrantArgs(queryFeeGrantMsg)
}

// (Query) make msg - fee grants by grantee
func MakeQueryFeeGrantsByGranteeMsg(queryFeeGrantMsg types.QueryFeeGrantMsg) (feegrant.QueryAllowancesRequest, error) {
	return parseQueryFeeGrantsByGranteeArgs(queryFeeGrantMsg)
}

// (Query) make msg - fee grants by granter
func MakeQueryFeeGrantsByGranterMsg(queryFeeGrantMsg types.QueryFeeGrantMsg) (feegrant.QueryAllowancesByGranterRequest, error) {
	return parseQueryFeeGrantsByGranterArgs(queryFeeGrantMsg)
}
