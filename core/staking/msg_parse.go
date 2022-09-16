package staking

import (
	"fmt"

	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Parsing - create validator
func parseCreateValidatorArgs(
	createValidatorMsg types.CreateValidatorMsg,
	output string,
) (sdk.Msg, error) {

	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(createValidatorMsg.HomeDir)

	nodeId, valPubKey, err := genutil.InitializeNodeValidatorFiles(serverCtx.Config)
	if err != nil {
		util.LogErr(err)
	}
	ip, _ := server.ExternalIP()

	if output != "" {
		if nodeId != "" && ip != "" {
			types.Memo = fmt.Sprintf("%s@%s:26656", nodeId, ip)
		}
	}

	website := createValidatorMsg.Website
	securityContact := createValidatorMsg.SecurityContact
	identity := createValidatorMsg.Identity
	details := createValidatorMsg.Details
	moniker := createValidatorMsg.Moniker

	var amount string
	if createValidatorMsg.Amount == "" {
		amount = types.DefaultAmount
	} else {
		amount = createValidatorMsg.Amount
	}

	var commissionRate string
	if createValidatorMsg.CommissionRate == "" {
		commissionRate = types.DefaultCommissionRate
	} else {
		commissionRate = createValidatorMsg.CommissionRate
	}

	var commissionMaxRate string
	if createValidatorMsg.CommissionMaxRate == "" {
		commissionMaxRate = types.DefaultCommissionMaxRate
	} else {
		commissionMaxRate = createValidatorMsg.CommissionMaxRate
	}

	var commissionMaxChangeRate string
	if createValidatorMsg.CommissionMaxChangeRate == "" {
		commissionMaxChangeRate = types.DefaultCommissionMaxChangeRate
	} else {
		commissionMaxChangeRate = createValidatorMsg.CommissionMaxChangeRate
	}

	var minSelfDelegation string
	if createValidatorMsg.MinSelfDelegation == "" {
		minSelfDelegation = types.DefaultMinSelfDelegation
	} else {
		minSelfDelegation = createValidatorMsg.MinSelfDelegation
	}

	description := stakingtypes.NewDescription(
		moniker, identity, website, securityContact, details)

	amountCoins, err := sdk.ParseCoinNormalized(util.DenomAdd(amount))
	if err != nil {
		return nil, err
	}

	buildCRates, err := buildCommissionRates(commissionRate, commissionMaxRate, commissionMaxChangeRate)
	if err != nil {
		return nil, err
	}

	intMinSelfDelegation, ok := sdk.NewIntFromString(minSelfDelegation)
	if !ok {
		return nil, util.LogErr("Wrong minSelfDelegation")
	}

	addr, err := sdk.AccAddressFromHex(valPubKey.Address().String())
	if err != nil {
		return nil, err
	}

	valAddr := sdk.ValAddress(addr)
	msg, err := stakingtypes.NewMsgCreateValidator(valAddr, valPubKey, amountCoins, description, buildCRates, intMinSelfDelegation)
	if err != nil {
		return nil, util.LogErr(err)
	}
	return msg, nil
}

// Parsing - edit validator
func parseEditValidatorArgs(editValidatorMsg types.EditValidatorMsg, privKey key.PrivateKey) (*stakingtypes.MsgEditValidator, error) {
	moniker := editValidatorMsg.Moniker
	identity := editValidatorMsg.Identity
	website := editValidatorMsg.Website
	security := editValidatorMsg.SecurityContact
	details := editValidatorMsg.Details
	description := stakingtypes.NewDescription(moniker, identity, website, security, details)

	var newRate *sdk.Dec

	commisionRate := editValidatorMsg.CommissionRate
	if commisionRate != "" {
		rate, err := sdk.NewDecFromStr(util.DenomAdd(commisionRate))
		if err != nil {
			return nil, util.LogErr(err)
		}
		newRate = &rate
	}

	var newMinSelfDelegation *sdk.Int

	minSelfDelegation := editValidatorMsg.MinSelfDelegation
	if minSelfDelegation != "" {
		msb, ok := sdk.NewIntFromString(minSelfDelegation)
		if !ok {
			return nil, util.LogErr("minimum self delegation must be a positive integer")
		}
		newMinSelfDelegation = &msb
	}

	addr := util.GetAddrByPrivKey(privKey)

	msg := stakingtypes.NewMsgEditValidator(sdk.ValAddress(addr), description, newRate, newMinSelfDelegation)

	return msg, nil
}

// Parsing - delegate
func parseDelegateArgs(delegateMsg types.DelegateMsg, privKey key.PrivateKey) (*stakingtypes.MsgDelegate, error) {
	amount, err := sdk.ParseCoinNormalized(util.DenomAdd(delegateMsg.Amount))
	if err != nil {
		return nil, err
	}
	delAddr := util.GetAddrByPrivKey(privKey)
	valAddr, err := sdk.ValAddressFromBech32(delegateMsg.ValAddr)
	if err != nil {
		return nil, err
	}

	msg := stakingtypes.NewMsgDelegate(delAddr, valAddr, amount)

	return msg, nil
}

// Parsing - redelegate
func parseRedelegateArgs(redelegateMsg types.RedelegateMsg, privKey key.PrivateKey) (*stakingtypes.MsgBeginRedelegate, error) {
	amount, err := sdk.ParseCoinNormalized(util.DenomAdd(redelegateMsg.Amount))
	if err != nil {
		return nil, err
	}
	delAddr := util.GetAddrByPrivKey(privKey)
	valSrcAddr, err := sdk.ValAddressFromBech32(redelegateMsg.ValSrcAddr)
	if err != nil {
		return nil, err
	}
	valDstAddr, err := sdk.ValAddressFromBech32(redelegateMsg.ValDstAddr)
	if err != nil {
		return nil, err
	}

	msg := stakingtypes.NewMsgBeginRedelegate(delAddr, valSrcAddr, valDstAddr, amount)
	return msg, nil
}

// Parsing - unbond
func parseUnbondArgs(unbondMsg types.UnbondMsg, privKey key.PrivateKey) (*stakingtypes.MsgUndelegate, error) {
	amount, err := sdk.ParseCoinNormalized(util.DenomAdd(unbondMsg.Amount))
	if err != nil {
		return nil, err
	}
	delAddr := util.GetAddrByPrivKey(privKey)
	valAddr, err := sdk.ValAddressFromBech32(unbondMsg.ValAddr)
	if err != nil {
		return nil, err
	}

	msg := stakingtypes.NewMsgUndelegate(delAddr, valAddr, amount)

	return msg, nil
}

// Parsing - validator
func parseQueryValidatorsArgs() (stakingtypes.QueryValidatorsRequest, error) {
	return stakingtypes.QueryValidatorsRequest{
		Pagination: core.PageRequest,
	}, nil
}

// Parsing - validators
func parseQueryValidatorArgs(queryValidatorMsg types.QueryValidatorMsg) (stakingtypes.QueryValidatorRequest, error) {
	return stakingtypes.QueryValidatorRequest{
		ValidatorAddr: queryValidatorMsg.ValidatorAddr,
	}, nil
}

// Parsing - query delegation
func parseQueryDelegationArgs(queryDelegationMsg types.QueryDelegationMsg) (stakingtypes.QueryDelegationRequest, error) {
	delAddr := queryDelegationMsg.DelegatorAddr
	valAddr := queryDelegationMsg.ValidatorAddr

	return stakingtypes.QueryDelegationRequest{
		DelegatorAddr: delAddr,
		ValidatorAddr: valAddr,
	}, nil
}

// Parsing - query delegations
func parseQueryDelegationsArgs(queryDelegationMsg types.QueryDelegationMsg) (stakingtypes.QueryDelegatorDelegationsRequest, error) {
	delAddr := queryDelegationMsg.DelegatorAddr

	return stakingtypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: delAddr,
		Pagination:    core.PageRequest,
	}, nil
}

// Parsing - query delegations to
func parseQueryDelegationsToArgs(queryDelegationMsg types.QueryDelegationMsg) (stakingtypes.QueryValidatorDelegationsRequest, error) {
	valAddr := queryDelegationMsg.ValidatorAddr

	return stakingtypes.QueryValidatorDelegationsRequest{
		ValidatorAddr: valAddr,
		Pagination:    core.PageRequest,
	}, nil
}

// Parsing - query unbonding delegation
func parseQueryUnbondingDelegationArgs(queryUnbondingDelegationMsg types.QueryUnbondingDelegationMsg) (stakingtypes.QueryUnbondingDelegationRequest, error) {
	delAddr := queryUnbondingDelegationMsg.DelegatorAddr
	valAddr := queryUnbondingDelegationMsg.ValidatorAddr

	return stakingtypes.QueryUnbondingDelegationRequest{
		DelegatorAddr: delAddr,
		ValidatorAddr: valAddr,
	}, nil
}

// Parsing - query unbonding delegations
func parseQueryUnbondingDelegationsArgs(queryUnbondingDelegationMsg types.QueryUnbondingDelegationMsg) (stakingtypes.QueryDelegatorUnbondingDelegationsRequest, error) {
	delAddr := queryUnbondingDelegationMsg.DelegatorAddr

	return stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
		DelegatorAddr: delAddr,
		Pagination:    core.PageRequest,
	}, nil
}

// Parsing - query unbonding delegations from
func parseQueryUnbondingDelegationsFromArgs(queryUnbondingDelegationMsg types.QueryUnbondingDelegationMsg) (stakingtypes.QueryValidatorUnbondingDelegationsRequest, error) {
	valAddr := queryUnbondingDelegationMsg.ValidatorAddr

	return stakingtypes.QueryValidatorUnbondingDelegationsRequest{
		ValidatorAddr: valAddr,
		Pagination:    core.PageRequest,
	}, nil
}

// Parsing - query redelegation
func parseQueryRedelegationArgs(queryRedelegationMsg types.QueryRedelegationMsg) (stakingtypes.QueryRedelegationsRequest, error) {
	delAddr := queryRedelegationMsg.DelegatorAddr
	valSrcAddr := queryRedelegationMsg.SrcValidatorAddr
	valDstAddr := queryRedelegationMsg.DstValidatorAddr

	return stakingtypes.QueryRedelegationsRequest{
		DelegatorAddr:    delAddr,
		DstValidatorAddr: valDstAddr,
		SrcValidatorAddr: valSrcAddr,
	}, nil
}

// Parsing - query redelegations
func parseQueryRedelegationsArgs(queryRedelegationMsg types.QueryRedelegationMsg) (stakingtypes.QueryRedelegationsRequest, error) {
	delAddr := queryRedelegationMsg.DelegatorAddr

	return stakingtypes.QueryRedelegationsRequest{
		DelegatorAddr: delAddr,
		Pagination:    core.PageRequest,
	}, nil
}

// Parsing - query redelegations from
func parseQueryRedelegationsFromArgs(queryRedelegationMsg types.QueryRedelegationMsg) (stakingtypes.QueryRedelegationsRequest, error) {
	valSrcAddr := queryRedelegationMsg.SrcValidatorAddr

	return stakingtypes.QueryRedelegationsRequest{
		SrcValidatorAddr: valSrcAddr,
		Pagination:       core.PageRequest,
	}, nil
}

// Parsing - historical
func parseHistoricalInfoArgs(historicalMsg types.HistoricalInfoMsg) (stakingtypes.QueryHistoricalInfoRequest, error) {
	height := historicalMsg.Height
	heightInt := util.FromStringToInt64(height)
	if heightInt < 0 {
		return stakingtypes.QueryHistoricalInfoRequest{}, util.LogErr("height argument provided must be a non-negative-integer")
	}

	return stakingtypes.QueryHistoricalInfoRequest{Height: heightInt}, nil
}

// Parsing - staking pool
func parseQueryStakingPoolArgs() (stakingtypes.QueryPoolRequest, error) {
	return stakingtypes.QueryPoolRequest{}, nil
}

// Parsing - staking params
func parseQueryStakingParamsArgs() (stakingtypes.QueryParamsRequest, error) {
	return stakingtypes.QueryParamsRequest{}, nil
}

// Build commission rate
func buildCommissionRates(rateStr, maxRateStr, maxChangeRateStr string) (commission stakingtypes.CommissionRates, err error) {
	if rateStr == "" || maxRateStr == "" || maxChangeRateStr == "" {
		return commission, util.LogErr("must specify all validator commission parameters")
	}
	rate, err := sdk.NewDecFromStr(rateStr)
	if err != nil {
		return commission, err
	}
	maxRate, err := sdk.NewDecFromStr(maxRateStr)
	if err != nil {
		return commission, err
	}
	maxChangeRate, err := sdk.NewDecFromStr(maxChangeRateStr)
	if err != nil {
		return commission, err
	}
	commission = stakingtypes.NewCommissionRates(rate, maxRate, maxChangeRate)

	return commission, nil
}
