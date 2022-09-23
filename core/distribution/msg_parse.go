package distribution

import (
	"context"

	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distcli "github.com/cosmos/cosmos-sdk/x/distribution/client/cli"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/grpc"
)

// Parsing - fund community pool
func parseFundCommunityPoolArgs(fundCommunityPoolMsg types.FundCommunityPoolMsg, privKey key.PrivateKey) (*disttypes.MsgFundCommunityPool, error) {
	depositorAddr := util.GetAddrByPrivKey(privKey)
	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(fundCommunityPoolMsg.Amount))
	if err != nil {
		return nil, err
	}

	msg := disttypes.NewMsgFundCommunityPool(amount, depositorAddr)
	return msg, nil
}

// Parsing - proposal community pool
func parseProposalCommunityPoolSpendArgs(communityPoolSpendMsg types.CommunityPoolSpendMsg, privKey key.PrivateKey, encodingConfig params.EncodingConfig) (*govtypes.MsgSubmitProposal, error) {
	var proposal disttypes.CommunityPoolSpendProposalWithDeposit
	var err error

	if communityPoolSpendMsg.JsonFilePath != "" {
		proposal, err = distcli.ParseCommunityPoolSpendProposalWithDeposit(encodingConfig.Marshaler, communityPoolSpendMsg.JsonFilePath)
		if err != nil {
			return nil, err
		}
	} else {
		proposal.Title = communityPoolSpendMsg.Title
		proposal.Description = communityPoolSpendMsg.Description
		proposal.Recipient = communityPoolSpendMsg.Recipient
		proposal.Amount = communityPoolSpendMsg.Amount
		proposal.Deposit = communityPoolSpendMsg.Deposit
	}

	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(proposal.Amount))
	if err != nil {
		return nil, err
	}

	deposit, err := sdk.ParseCoinsNormalized(util.DenomAdd(proposal.Deposit))
	if err != nil {
		return nil, err
	}

	from := util.GetAddrByPrivKey(privKey)
	recpAddr, err := sdk.AccAddressFromBech32(proposal.Recipient)
	if err != nil {
		return nil, err
	}

	content := disttypes.NewCommunityPoolSpendProposal(proposal.Title, proposal.Description, recpAddr, amount)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Parsing - withdraw rewards
func parseWithdrawRewardsArgs(withdrawRewardsMsg types.WithdrawRewardsMsg, privKey key.PrivateKey) ([]sdk.Msg, error) {
	delAddr := util.GetAddrByPrivKey(privKey)
	valAddr, err := sdk.ValAddressFromBech32(withdrawRewardsMsg.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	msgs := []sdk.Msg{disttypes.NewMsgWithdrawDelegatorReward(delAddr, valAddr)}
	if withdrawRewardsMsg.Commission {
		msgs = append(msgs, disttypes.NewMsgWithdrawValidatorCommission(valAddr))
	}

	for _, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return nil, err
		}
	}

	return msgs, nil
}

// Parsing - withdraw all rewards
func parseWithdrawAllRewardsArgs(privKey key.PrivateKey, grpcConn grpc.ClientConn, ctx context.Context) ([]sdk.Msg, error) {
	delAddr := util.GetAddrByPrivKey(privKey)
	queryClient := disttypes.NewQueryClient(grpcConn)
	delValsRes, err := queryClient.DelegatorValidators(
		ctx,
		&disttypes.QueryDelegatorValidatorsRequest{
			DelegatorAddress: delAddr.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	vals := delValsRes.Validators
	msgs := make([]sdk.Msg, 0, len(vals))
	for _, valAddr := range vals {
		val, err := sdk.ValAddressFromBech32(valAddr)
		if err != nil {
			return nil, err
		}

		msg := disttypes.NewMsgWithdrawDelegatorReward(delAddr, val)
		if err := msg.ValidateBasic(); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, err
}

// Parsing - set withdraw addr
func parseSetWithdrawAddrArgs(setWithdrawAddrMsg types.SetwithdrawAddrMsg, privKey key.PrivateKey) (*disttypes.MsgSetWithdrawAddress, error) {
	delAddr := util.GetAddrByPrivKey(privKey)
	withdrawAddr, err := sdk.AccAddressFromBech32(setWithdrawAddrMsg.WithdrawAddr)
	if err != nil {
		return nil, err
	}

	msg := disttypes.NewMsgSetWithdrawAddress(delAddr, withdrawAddr)

	return msg, nil
}

// Parsing - distribution params
func parseQueryDistributionParamsArgs() (disttypes.QueryParamsRequest, error) {
	return disttypes.QueryParamsRequest{}, nil
}

// Parsing - validator outstanding rewards
func parseValidatorOutstandingRewardsArgs(validatorOutstandingRewardsMsg types.ValidatorOutstandingRewardsMsg) (disttypes.QueryValidatorOutstandingRewardsRequest, error) {
	valAddr, err := sdk.ValAddressFromBech32(validatorOutstandingRewardsMsg.ValidatorAddr)
	if err != nil {
		return disttypes.QueryValidatorOutstandingRewardsRequest{}, err
	}

	return disttypes.QueryValidatorOutstandingRewardsRequest{
		ValidatorAddress: valAddr.String(),
	}, nil
}

// Parsing - commission
func parseQueryDistCommissionArgs(queryDistCommissionMsg types.QueryDistCommissionMsg) (disttypes.QueryValidatorCommissionRequest, error) {
	valAddr, err := sdk.ValAddressFromBech32(queryDistCommissionMsg.ValidatorAddr)
	if err != nil {
		return disttypes.QueryValidatorCommissionRequest{}, err
	}

	return disttypes.QueryValidatorCommissionRequest{
		ValidatorAddress: valAddr.String(),
	}, nil
}

// Parsing - distribution slashes
func parseDistSlashesArgs(queryDistSlashesMsg types.QueryDistSlashesMsg) (disttypes.QueryValidatorSlashesRequest, error) {
	valAddr, err := sdk.ValAddressFromBech32(queryDistSlashesMsg.ValidatorAddr)
	if err != nil {
		return disttypes.QueryValidatorSlashesRequest{}, nil
	}
	startHeightNumber := util.FromStringToUint64(queryDistSlashesMsg.StartHeight)
	endHeightNumber := util.FromStringToUint64(queryDistSlashesMsg.EndHeight)

	pageReq := core.PageRequest

	return disttypes.QueryValidatorSlashesRequest{
		ValidatorAddress: valAddr.String(),
		StartingHeight:   startHeightNumber,
		EndingHeight:     endHeightNumber,
		Pagination:       pageReq,
	}, nil
}

// Parsing - distribution rewards
func parseQueryDistRewardsArgs(queryDistRewardsMsg types.QueryDistRewardsMsg, privKey key.PrivateKey) (disttypes.QueryDelegationRewardsRequest, error) {
	delAddr := util.GetAddrByPrivKey(privKey)
	valAddr, err := sdk.ValAddressFromBech32(queryDistRewardsMsg.ValidatorAddr)
	if err != nil {
		return disttypes.QueryDelegationRewardsRequest{}, err
	}

	return disttypes.QueryDelegationRewardsRequest{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
	}, nil
}

// Parsing - community pool
func parseQueryCommunityPoolArgs() (disttypes.QueryCommunityPoolRequest, error) {
	return disttypes.QueryCommunityPoolRequest{}, nil
}
