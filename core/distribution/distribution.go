package distribution

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"golang.org/x/net/context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/grpc"
	"github.com/xpladev/xpla/app/params"
)

// (Tx) make msg - proposal community pool
func MakeProposalCommunityPoolSpendMsg(communityPoolSpendMsg types.CommunityPoolSpendMsg, privKey key.PrivateKey, encodingConfig params.EncodingConfig) (*govtypes.MsgSubmitProposal, error) {
	msg, err := parseProposalCommunityPoolSpendArgs(communityPoolSpendMsg, privKey, encodingConfig)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// (Tx) make msg - fund community pool
func MakeFundCommunityPoolMsg(fundCommunityPoolMsg types.FundCommunityPoolMsg, privKey key.PrivateKey) (*disttypes.MsgFundCommunityPool, error) {
	msg, err := parseFundCommunityPoolArgs(fundCommunityPoolMsg, privKey)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// (Tx) make msg - withdraw rewards
func MakeWithdrawRewardsMsg(withdrawRewardsMsg types.WithdrawRewardsMsg, privKey key.PrivateKey) ([]sdk.Msg, error) {
	msg, err := parseWithdrawRewardsArgs(withdrawRewardsMsg, privKey)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// (Tx) make msg - withdraw all rewards
func MakeWithdrawAllRewardsMsg(privKey key.PrivateKey, grpcConn grpc.ClientConn, ctx context.Context) ([]sdk.Msg, error) {
	msg, err := parseWithdrawAllRewardsArgs(privKey, grpcConn, ctx)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// (Tx) make msg - withdraw address
func MakeSetWithdrawAddrMsg(setWithdrawAddrMsg types.SetwithdrawAddrMsg, privKey key.PrivateKey) (*disttypes.MsgSetWithdrawAddress, error) {
	msg, err := parseSetWithdrawAddrArgs(setWithdrawAddrMsg, privKey)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// (Query) make msg - distribution params
func MakeQueryDistributionParamsMsg() (disttypes.QueryParamsRequest, error) {
	msg, err := parseQueryDistributionParamsArgs()
	if err != nil {
		return disttypes.QueryParamsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - validator outstanding rewards
func MakeValidatorOutstandingRewardsMsg(validatorOutstandingRewardsMsg types.ValidatorOutstandingRewardsMsg) (disttypes.QueryValidatorOutstandingRewardsRequest, error) {
	msg, err := parseValidatorOutstandingRewardsArgs(validatorOutstandingRewardsMsg)
	if err != nil {
		return disttypes.QueryValidatorOutstandingRewardsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - commission
func MakeQueryDistCommissionMsg(queryDistCommissionMsg types.QueryDistCommissionMsg) (disttypes.QueryValidatorCommissionRequest, error) {
	msg, err := parseQueryDistCommissionArgs(queryDistCommissionMsg)
	if err != nil {
		return disttypes.QueryValidatorCommissionRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - distribution slashes
func MakeQueryDistSlashesMsg(queryDistSlashesMsg types.QueryDistSlashesMsg) (disttypes.QueryValidatorSlashesRequest, error) {
	msg, err := parseDistSlashesArgs(queryDistSlashesMsg)
	if err != nil {
		return disttypes.QueryValidatorSlashesRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - distribution rewards
func MakeyQueryDistRewardsMsg(queryDistRewardsMsg types.QueryDistRewardsMsg, privKey key.PrivateKey) (disttypes.QueryDelegationRewardsRequest, error) {
	msg, err := parseQueryDistRewardsArgs(queryDistRewardsMsg, privKey)
	if err != nil {
		return disttypes.QueryDelegationRewardsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - community pool
func MakeQueryCommunityPoolMsg() (disttypes.QueryCommunityPoolRequest, error) {
	msg, err := parseQueryCommunityPoolArgs()
	if err != nil {
		return disttypes.QueryCommunityPoolRequest{}, err
	}

	return msg, nil
}
