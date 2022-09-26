package module

import (
	mstaking "github.com/Moonyongjung/xpla.go/core/staking"
	"github.com/Moonyongjung/xpla.go/util"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Query client for staking module.
func (i IXplaClient) QueryStaking() (string, error) {
	queryClient := stakingtypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Skating validator
	case i.Ixplac.GetMsgType() == mstaking.StakingQueryValidatorMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(stakingtypes.QueryValidatorRequest)
		res, err = queryClient.Validator(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking validators
	case i.Ixplac.GetMsgType() == mstaking.StakingQueryValidatorsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(stakingtypes.QueryValidatorsRequest)
		res, err = queryClient.Validators(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking delegation
	case i.Ixplac.GetMsgType() == mstaking.StakingQueryDelegationMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(stakingtypes.QueryDelegationRequest)
		res, err = queryClient.Delegation(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking delegations
	case i.Ixplac.GetMsgType() == mstaking.StakingQueryDelegationsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(stakingtypes.QueryDelegatorDelegationsRequest)
		res, err = queryClient.DelegatorDelegations(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking delegations to
	case i.Ixplac.GetMsgType() == mstaking.StakingQueryDelegationsToMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(stakingtypes.QueryValidatorDelegationsRequest)
		res, err = queryClient.ValidatorDelegations(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking unbonding delegation
	case i.Ixplac.GetMsgType() == mstaking.StakingQueryUnbondingDelegationMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(stakingtypes.QueryUnbondingDelegationRequest)
		res, err = queryClient.UnbondingDelegation(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking unbonding delegations
	case i.Ixplac.GetMsgType() == mstaking.StakingQueryUnbondingDelegationsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(stakingtypes.QueryDelegatorUnbondingDelegationsRequest)
		res, err = queryClient.DelegatorUnbondingDelegations(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking unbonding delegations from
	case i.Ixplac.GetMsgType() == mstaking.StakingQueryUnbondingDelegationsFromMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(stakingtypes.QueryValidatorUnbondingDelegationsRequest)
		res, err = queryClient.ValidatorUnbondingDelegations(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking redelegations
	case i.Ixplac.GetMsgType() == mstaking.StakingQueryRedelegationMsgType ||
		i.Ixplac.GetMsgType() == mstaking.StakingQueryRedelegationsMsgType ||
		i.Ixplac.GetMsgType() == mstaking.StakingQueryRedelegationsFromMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(stakingtypes.QueryRedelegationsRequest)
		res, err = queryClient.Redelegations(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking historical information
	case i.Ixplac.GetMsgType() == mstaking.StakingHistoricalInfoMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(stakingtypes.QueryHistoricalInfoRequest)
		res, err = queryClient.HistoricalInfo(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking pool
	case i.Ixplac.GetMsgType() == mstaking.StakingQueryStakingPoolMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(stakingtypes.QueryPoolRequest)
		res, err = queryClient.Pool(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Staking params
	case i.Ixplac.GetMsgType() == mstaking.StakingQueryStakingParamsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(stakingtypes.QueryParamsRequest)
		res, err = queryClient.Params(
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
