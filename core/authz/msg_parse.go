package authz

import (
	"time"

	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/xpladev/xpla/app/params"
)

// Parsing - authz grant
func parseAuthzGrantArgs(authzGrantMsg types.AuthzGrantMsg, privKey key.PrivateKey) (authz.MsgGrant, error) {
	granter, err := util.GetAddrByPrivKey(privKey)
	if err != nil {
		return authz.MsgGrant{}, util.LogErr(errors.ErrParse, err)
	}
	if authzGrantMsg.Granter != granter.String() {
		return authz.MsgGrant{}, util.LogErr(errors.ErrAccountNotMatch, "Account address generated by private key is not equal input granter of msg")
	}
	grantee, err := sdk.AccAddressFromBech32(authzGrantMsg.Grantee)
	if err != nil {
		return authz.MsgGrant{}, util.LogErr(errors.ErrParse)
	}
	exp := util.FromStringToInt64(authzGrantMsg.Expiration)

	var authorization authz.Authorization

	switch authzGrantMsg.AuthorizationType {
	case "send":
		limit := authzGrantMsg.SpendLimit
		if limit == "" {
			return authz.MsgGrant{}, util.LogErr(errors.ErrInsufficientParams, "require bank spend limit")
		}
		spendLimit, err := sdk.ParseCoinsNormalized(util.DenomAdd(limit))
		if err != nil {
			return authz.MsgGrant{}, util.LogErr(errors.ErrParse, err)
		}
		if !spendLimit.IsAllPositive() {
			return authz.MsgGrant{}, util.LogErr(errors.ErrInvalidRequest, "spend-limit should be greater than zero")
		}
		authorization = banktypes.NewSendAuthorization(spendLimit)

	case "generic":
		msgType := authzGrantMsg.MsgType
		authorization = authz.NewGenericAuthorization(msgType)

	case "delegate", "unbond", "redelegate":
		limit := authzGrantMsg.SpendLimit
		if limit == "" {
			return authz.MsgGrant{}, util.LogErr(errors.ErrInsufficientParams, "require spend limit")
		}
		allowValidators := authzGrantMsg.AllowValidators
		denyValidators := authzGrantMsg.DenyValidators

		var delegateLimit *sdk.Coin

		spendLimit, err := sdk.ParseCoinsNormalized(util.DenomAdd(limit))
		if err != nil {
			return authz.MsgGrant{}, util.LogErr(errors.ErrParse, err)
		}

		if !spendLimit.IsAllPositive() {
			return authz.MsgGrant{}, util.LogErr(errors.ErrInvalidRequest, "spend-limit should be greater than zero")
		}
		delegateLimit = &spendLimit[0]

		allowed, err := util.Bech32toValidatorAddress(allowValidators)
		if err != nil {
			return authz.MsgGrant{}, util.LogErr(errors.ErrParse, err)
		}

		denied, err := util.Bech32toValidatorAddress(denyValidators)
		if err != nil {
			return authz.MsgGrant{}, util.LogErr(errors.ErrParse, err)
		}

		switch authzGrantMsg.AuthorizationType {
		case "delegate":
			authorization, err = stakingtypes.NewStakeAuthorization(allowed, denied, stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_DELEGATE, delegateLimit)
		case "unbond":
			authorization, err = stakingtypes.NewStakeAuthorization(allowed, denied, stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_UNDELEGATE, delegateLimit)
		default:
			authorization, err = stakingtypes.NewStakeAuthorization(allowed, denied, stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_REDELEGATE, delegateLimit)
		}
		if err != nil {
			return authz.MsgGrant{}, util.LogErr(errors.ErrParse, err)
		}
	default:
		return authz.MsgGrant{}, util.LogErr(errors.ErrInvalidMsgType, "invalid authorization type, ", authzGrantMsg.AuthorizationType)
	}

	msg, err := authz.NewMsgGrant(granter, grantee, authorization, time.Unix(exp, 0))
	if err != nil {
		return authz.MsgGrant{}, util.LogErr(errors.ErrParse, err)
	}
	return *msg, nil
}

// Parsing - revoke
func parseAuthzRevokeArgs(authzRevokeMsg types.AuthzRevokeMsg, privKey key.PrivateKey) (authz.MsgRevoke, error) {
	granter, err := util.GetAddrByPrivKey(privKey)
	if err != nil {
		return authz.MsgRevoke{}, util.LogErr(errors.ErrParse, err)
	}
	if authzRevokeMsg.Granter != granter.String() {
		return authz.MsgRevoke{}, util.LogErr(errors.ErrAccountNotMatch, "Account address generated by private key is not equal input granter of msg")
	}
	grantee, err := sdk.AccAddressFromBech32(authzRevokeMsg.Grantee)
	if err != nil {
		return authz.MsgRevoke{}, util.LogErr(errors.ErrParse, err)
	}

	msg := authz.NewMsgRevoke(granter, grantee, authzRevokeMsg.MsgType)
	return msg, nil
}

// Parsing - authz execute
func parseAuthzExecArgs(authzExecMsg types.AuthzExecMsg, encodingConfig params.EncodingConfig) (authz.MsgExec, error) {
	var readTx sdk.Tx
	grantee, err := sdk.AccAddressFromBech32(authzExecMsg.Grantee)
	if err != nil {
		return authz.MsgExec{}, util.LogErr(errors.ErrParse, err)
	}

	clientCtx, err := util.NewClient()
	if err != nil {
		return authz.MsgExec{}, err
	}

	if authzExecMsg.ExecFile != "" {
		readTx, err = authclient.ReadTxFromFile(clientCtx, authzExecMsg.ExecFile)
		if err != nil {
			return authz.MsgExec{}, util.LogErr(errors.ErrParse, err)
		}
	} else if authzExecMsg.ExecTxString != "" {
		readTx, err = encodingConfig.TxConfig.TxJSONDecoder()([]byte(authzExecMsg.ExecTxString))
		if err != nil {
			return authz.MsgExec{}, util.LogErr(errors.ErrParse, err)
		}
	} else {
		return authz.MsgExec{}, util.LogErr(errors.ErrInsufficientParams, "no authz exec info")
	}

	msg := authz.NewMsgExec(grantee, readTx.GetMsgs())

	return msg, nil
}

// Parsing - authz grants
func parseQueryAuthzGrantsArgs(queryAuthzGrantMsg types.QueryAuthzGrantMsg) (authz.QueryGrantsRequest, error) {
	granter, err := sdk.AccAddressFromBech32(queryAuthzGrantMsg.Granter)
	if err != nil {
		return authz.QueryGrantsRequest{}, util.LogErr(errors.ErrParse, err)
	}
	if queryAuthzGrantMsg.Granter != granter.String() {
		return authz.QueryGrantsRequest{}, util.LogErr(errors.ErrAccountNotMatch, "Account address generated by private key is not equal input granter of msg")
	}
	grantee, err := sdk.AccAddressFromBech32(queryAuthzGrantMsg.Grantee)
	if err != nil {
		return authz.QueryGrantsRequest{}, util.LogErr(errors.ErrParse, err)
	}

	msgAuthorized := ""
	if queryAuthzGrantMsg.MsgType != "" {
		msgAuthorized = queryAuthzGrantMsg.MsgType
	}

	return authz.QueryGrantsRequest{
		Granter:    granter.String(),
		Grantee:    grantee.String(),
		MsgTypeUrl: msgAuthorized,
		Pagination: core.PageRequest,
	}, nil
}

// Parsing - authz grants by grantee
func parseQueryAuthzGrantsByGranteeArgs(queryAuthzGrantMsg types.QueryAuthzGrantMsg) (authz.QueryGranteeGrantsRequest, error) {
	grantee, err := sdk.AccAddressFromBech32(queryAuthzGrantMsg.Grantee)
	if err != nil {
		return authz.QueryGranteeGrantsRequest{}, util.LogErr(errors.ErrParse, err)
	}

	return authz.QueryGranteeGrantsRequest{
		Grantee:    grantee.String(),
		Pagination: core.PageRequest,
	}, nil
}

// Parsing - authz grants by granter
func parseQueryAuthzGrantsByGranterArgs(queryAuthzGrantMsg types.QueryAuthzGrantMsg) (authz.QueryGranterGrantsRequest, error) {
	granter, err := sdk.AccAddressFromBech32(queryAuthzGrantMsg.Granter)
	if err != nil {
		return authz.QueryGranterGrantsRequest{}, util.LogErr(errors.ErrParse, err)
	}

	return authz.QueryGranterGrantsRequest{
		Granter:    granter.String(),
		Pagination: core.PageRequest,
	}, nil
}
