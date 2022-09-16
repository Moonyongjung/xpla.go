package client

import (
	"context"

	mstaking "github.com/Moonyongjung/xpla.go/core/staking"
	"github.com/Moonyongjung/xpla.go/util"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Query client for staking module.
func queryStaking(xplac *XplaClient) (string, error) {
	queryClient := stakingtypes.NewQueryClient(xplac.Grpc)

	switch {
	// Skating validator
	case xplac.MsgType == mstaking.StakingQueryValidatorMsgType:
		convertMsg, _ := xplac.Msg.(stakingtypes.QueryValidatorRequest)
		res, err = queryClient.Validator(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking validators
	case xplac.MsgType == mstaking.StakingQueryValidatorsMsgType:
		convertMsg, _ := xplac.Msg.(stakingtypes.QueryValidatorsRequest)
		res, err = queryClient.Validators(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking delegation
	case xplac.MsgType == mstaking.StakingQueryDelegationMsgType:
		convertMsg, _ := xplac.Msg.(stakingtypes.QueryDelegationRequest)
		res, err = queryClient.Delegation(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking delegations
	case xplac.MsgType == mstaking.StakingQueryDelegationsMsgType:
		convertMsg, _ := xplac.Msg.(stakingtypes.QueryDelegatorDelegationsRequest)
		res, err = queryClient.DelegatorDelegations(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking delegations to
	case xplac.MsgType == mstaking.StakingQueryDelegationsToMsgType:
		convertMsg, _ := xplac.Msg.(stakingtypes.QueryValidatorDelegationsRequest)
		res, err = queryClient.ValidatorDelegations(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking unbonding delegation
	case xplac.MsgType == mstaking.StakingQueryUnbondingDelegationMsgType:
		convertMsg, _ := xplac.Msg.(stakingtypes.QueryUnbondingDelegationRequest)
		res, err = queryClient.UnbondingDelegation(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking unbonding delegations
	case xplac.MsgType == mstaking.StakingQueryUnbondingDelegationsMsgType:
		convertMsg, _ := xplac.Msg.(stakingtypes.QueryDelegatorUnbondingDelegationsRequest)
		res, err = queryClient.DelegatorUnbondingDelegations(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking unbonding delegations from
	case xplac.MsgType == mstaking.StakingQueryUnbondingDelegationsFromMsgType:
		convertMsg, _ := xplac.Msg.(stakingtypes.QueryValidatorUnbondingDelegationsRequest)
		res, err = queryClient.ValidatorUnbondingDelegations(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking redelegations
	case xplac.MsgType == mstaking.StakingQueryRedelegationMsgType ||
		xplac.MsgType == mstaking.StakingQueryRedelegationsMsgType ||
		xplac.MsgType == mstaking.StakingQueryRedelegationsFromMsgType:
		convertMsg, _ := xplac.Msg.(stakingtypes.QueryRedelegationsRequest)
		res, err = queryClient.Redelegations(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking historical information
	case xplac.MsgType == mstaking.StakingHistoricalInfoMsgType:
		convertMsg, _ := xplac.Msg.(stakingtypes.QueryHistoricalInfoRequest)
		res, err = queryClient.HistoricalInfo(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking pool
	case xplac.MsgType == mstaking.StakingQueryStakingPoolMsgType:
		convertMsg, _ := xplac.Msg.(stakingtypes.QueryPoolRequest)
		res, err = queryClient.Pool(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking params
	case xplac.MsgType == mstaking.StakingQueryStakingParamsMsgType:
		convertMsg, _ := xplac.Msg.(stakingtypes.QueryParamsRequest)
		res, err = queryClient.Params(
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
