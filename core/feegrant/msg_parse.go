package feegrant

import (
	"time"

	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// Parsing - fee grant
func parseFeeGrantArgs(feeGrantMsg types.FeeGrantMsg, privKey key.PrivateKey) (feegrant.MsgGrantAllowance, error) {
	granter, err := util.GetAddrByPrivKey(privKey)
	if err != nil {
		return feegrant.MsgGrantAllowance{}, err
	}

	if feeGrantMsg.Granter != granter.String() {
		return feegrant.MsgGrantAllowance{}, util.LogErr(errors.ErrAccountNotMatch, "Account address generated by private key is not equal input granter of msg")
	}

	grantee, err := sdk.AccAddressFromBech32(feeGrantMsg.Grantee)
	if err != nil {
		return feegrant.MsgGrantAllowance{}, err
	}

	limit := util.DenomAdd(feeGrantMsg.SpendLimit)
	if feeGrantMsg.SpendLimit == "" {
		limit = ""
	}
	spendLimit, err := sdk.ParseCoinsNormalized(limit)
	if err != nil {
		return feegrant.MsgGrantAllowance{}, err
	}

	basic := feegrant.BasicAllowance{
		SpendLimit: spendLimit,
	}

	var expireTime time.Time
	if feeGrantMsg.Expiration != "" {
		expireTime, err := time.Parse(time.RFC3339, feeGrantMsg.Expiration)
		if err != nil {
			return feegrant.MsgGrantAllowance{}, err
		}
		basic.Expiration = &expireTime
	}

	var grant feegrant.FeeAllowanceI
	grant = &basic

	periodClock := util.FromStringToInt64(feeGrantMsg.Period)

	if periodClock > 0 || feeGrantMsg.PeriodLimit != "" {
		periodLimit, err := sdk.ParseCoinsNormalized(util.DenomAdd(feeGrantMsg.PeriodLimit))
		if err != nil {
			return feegrant.MsgGrantAllowance{}, err
		}

		if periodClock <= 0 {
			return feegrant.MsgGrantAllowance{}, util.LogErr(errors.ErrInsufficientParams, "period clock was not set")
		}

		if periodLimit == nil {
			return feegrant.MsgGrantAllowance{}, util.LogErr(errors.ErrInsufficientParams, "period limit was not set")
		}

		periodReset := getPeriodReset(periodClock)
		if feeGrantMsg.Expiration != "" && periodReset.Sub(expireTime) > 0 {
			return feegrant.MsgGrantAllowance{}, util.LogErr(errors.ErrInvalidRequest, "period (", periodClock, ") cannot reset after expiration (", feeGrantMsg.Expiration, ")")
		}

		periodic := feegrant.PeriodicAllowance{
			Basic:            basic,
			Period:           getPeriod(periodClock),
			PeriodReset:      getPeriodReset(periodClock),
			PeriodSpendLimit: periodLimit,
			PeriodCanSpend:   periodLimit,
		}

		grant = &periodic
	}

	if len(feeGrantMsg.AllowedMsg) > 0 {
		grant, err = feegrant.NewAllowedMsgAllowance(grant, feeGrantMsg.AllowedMsg)
		if err != nil {
			return feegrant.MsgGrantAllowance{}, err
		}
	}

	msg, err := feegrant.NewMsgGrantAllowance(grant, granter, grantee)
	if err != nil {
		return feegrant.MsgGrantAllowance{}, err
	}

	return *msg, nil
}

// Parsing - fee grant revoke
func parseRevokeFeeGrantArgs(revokeFeeGrantMsg types.RevokeFeeGrantMsg, privKey key.PrivateKey) (feegrant.MsgRevokeAllowance, error) {
	granter, err := util.GetAddrByPrivKey(privKey)
	if err != nil {
		return feegrant.MsgRevokeAllowance{}, err
	}

	if revokeFeeGrantMsg.Granter != granter.String() {
		return feegrant.MsgRevokeAllowance{}, util.LogErr(errors.ErrAccountNotMatch, "Account address generated by private key is not equal input granter of msg")
	}
	grantee, err := sdk.AccAddressFromBech32(revokeFeeGrantMsg.Grantee)
	if err != nil {
		return feegrant.MsgRevokeAllowance{}, err
	}

	msg := feegrant.NewMsgRevokeAllowance(granter, grantee)

	return msg, nil
}

// Parsing - query grants
func parseQueryFeeGrantArgs(queryGrantMsg types.QueryFeeGrantMsg) (feegrant.QueryAllowanceRequest, error) {
	granter, err := sdk.AccAddressFromBech32(queryGrantMsg.Granter)
	if err != nil {
		return feegrant.QueryAllowanceRequest{}, err
	}
	grantee, err := sdk.AccAddressFromBech32(queryGrantMsg.Grantee)
	if err != nil {
		return feegrant.QueryAllowanceRequest{}, err
	}

	return feegrant.QueryAllowanceRequest{
		Granter: granter.String(),
		Grantee: grantee.String(),
	}, nil
}

// Parsing - grants by grantee
func parseQueryFeeGrantsByGranteeArgs(queryGrantMsg types.QueryFeeGrantMsg) (feegrant.QueryAllowancesRequest, error) {
	grantee, err := sdk.AccAddressFromBech32(queryGrantMsg.Grantee)
	if err != nil {
		return feegrant.QueryAllowancesRequest{}, err
	}

	return feegrant.QueryAllowancesRequest{
		Grantee:    grantee.String(),
		Pagination: core.PageRequest,
	}, nil
}

// Parsing - grants by granter
func parseQueryFeeGrantsByGranterArgs(queryGrantMsg types.QueryFeeGrantMsg) (feegrant.QueryAllowancesByGranterRequest, error) {
	granter, err := sdk.AccAddressFromBech32(queryGrantMsg.Granter)
	if err != nil {
		return feegrant.QueryAllowancesByGranterRequest{}, err
	}

	return feegrant.QueryAllowancesByGranterRequest{
		Granter:    granter.String(),
		Pagination: core.PageRequest,
	}, nil
}

func getPeriodReset(duration int64) time.Time {
	return time.Now().Add(getPeriod(duration))
}

func getPeriod(duration int64) time.Duration {
	return time.Duration(duration) * time.Second
}
