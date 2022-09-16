package staking

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// (Tx) make msg - create validator
func MakeCreateValidatorMsg(createValidatorMsg types.CreateValidatorMsg, output string) (sdk.Msg, error) {
	msg, err := parseCreateValidatorArgs(createValidatorMsg, output)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// (Tx) make msg - edit validator
func MakeEditValidatorMsg(editValidatorMsg types.EditValidatorMsg, privKey key.PrivateKey) (*stakingtypes.MsgEditValidator, error) {
	msg, err := parseEditValidatorArgs(editValidatorMsg, privKey)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// (Tx) make msg - delegate
func MakeDelegateMsg(delegateMsg types.DelegateMsg, privKey key.PrivateKey) (*stakingtypes.MsgDelegate, error) {
	msg, err := parseDelegateArgs(delegateMsg, privKey)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// (Tx) make msg - unbond
func MakeUnbondMsg(unbondMsg types.UnbondMsg, privKey key.PrivateKey) (*stakingtypes.MsgUndelegate, error) {
	msg, err := parseUnbondArgs(unbondMsg, privKey)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// (Tx) make msg - redelegate
func MakeRedelegateMsg(redelegateMsg types.RedelegateMsg, privKey key.PrivateKey) (*stakingtypes.MsgBeginRedelegate, error) {
	msg, err := parseRedelegateArgs(redelegateMsg, privKey)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// (Query) make msg - validator
func MakeQueryValidatorMsg(queryValidatorMsg types.QueryValidatorMsg) (stakingtypes.QueryValidatorRequest, error) {
	msg, err := parseQueryValidatorArgs(queryValidatorMsg)
	if err != nil {
		return stakingtypes.QueryValidatorRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - validators
func MakeQueryValidatorsMsg() (stakingtypes.QueryValidatorsRequest, error) {
	msg, err := parseQueryValidatorsArgs()
	if err != nil {
		return stakingtypes.QueryValidatorsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - query delegation
func MakeQueryDelegationMsg(queryDelegationMsg types.QueryDelegationMsg) (stakingtypes.QueryDelegationRequest, error) {
	msg, err := parseQueryDelegationArgs(queryDelegationMsg)
	if err != nil {
		return stakingtypes.QueryDelegationRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - query delegations
func MakeQueryDelegationsMsg(queryDelegationMsg types.QueryDelegationMsg) (stakingtypes.QueryDelegatorDelegationsRequest, error) {
	msg, err := parseQueryDelegationsArgs(queryDelegationMsg)
	if err != nil {
		return stakingtypes.QueryDelegatorDelegationsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - query delegations to
func MakeQueryDelegationsToMsg(queryDelegationMsg types.QueryDelegationMsg) (stakingtypes.QueryValidatorDelegationsRequest, error) {
	msg, err := parseQueryDelegationsToArgs(queryDelegationMsg)
	if err != nil {
		return stakingtypes.QueryValidatorDelegationsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - query unbonding delegation
func MakeQueryUnbondingDelegationMsg(queryUnbondingDelegationMsg types.QueryUnbondingDelegationMsg) (stakingtypes.QueryUnbondingDelegationRequest, error) {
	msg, err := parseQueryUnbondingDelegationArgs(queryUnbondingDelegationMsg)
	if err != nil {
		return stakingtypes.QueryUnbondingDelegationRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - query unbonding delegations
func MakeQueryUnbondingDelegationsMsg(queryUnbondingDelegationMsg types.QueryUnbondingDelegationMsg) (stakingtypes.QueryDelegatorUnbondingDelegationsRequest, error) {
	msg, err := parseQueryUnbondingDelegationsArgs(queryUnbondingDelegationMsg)
	if err != nil {
		return stakingtypes.QueryDelegatorUnbondingDelegationsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - query unbonding delegations from
func MakeQueryUnbondingDelegationsFromMsg(queryUnbondingDelegationMsg types.QueryUnbondingDelegationMsg) (stakingtypes.QueryValidatorUnbondingDelegationsRequest, error) {
	msg, err := parseQueryUnbondingDelegationsFromArgs(queryUnbondingDelegationMsg)
	if err != nil {
		return stakingtypes.QueryValidatorUnbondingDelegationsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - query redelegation
func MakeQueryRedelegationMsg(queryRedelegationMsg types.QueryRedelegationMsg) (stakingtypes.QueryRedelegationsRequest, error) {
	msg, err := parseQueryRedelegationArgs(queryRedelegationMsg)
	if err != nil {
		return stakingtypes.QueryRedelegationsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - query redelegations
func MakeQueryRedelegationsMsg(queryRedelegationMsg types.QueryRedelegationMsg) (stakingtypes.QueryRedelegationsRequest, error) {
	msg, err := parseQueryRedelegationsArgs(queryRedelegationMsg)
	if err != nil {
		return stakingtypes.QueryRedelegationsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - query redelegations from
func MakeQueryRedelegationsFromMsg(queryRedelegationMsg types.QueryRedelegationMsg) (stakingtypes.QueryRedelegationsRequest, error) {
	msg, err := parseQueryRedelegationsFromArgs(queryRedelegationMsg)
	if err != nil {
		return stakingtypes.QueryRedelegationsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - historical
func MakeHistoricalInfoMsg(historicalInfoMsg types.HistoricalInfoMsg) (stakingtypes.QueryHistoricalInfoRequest, error) {
	msg, err := parseHistoricalInfoArgs(historicalInfoMsg)
	if err != nil {
		return stakingtypes.QueryHistoricalInfoRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - staking pool
func MakeQueryStakingPoolMsg() (stakingtypes.QueryPoolRequest, error) {
	msg, _ := parseQueryStakingPoolArgs()
	return msg, nil
}

// (Query) make msg - staking params
func MakeQueryStakingParamsMsg() (stakingtypes.QueryParamsRequest, error) {
	msg, _ := parseQueryStakingParamsArgs()
	return msg, nil
}
