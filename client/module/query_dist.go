package module

import (
	mdist "github.com/Moonyongjung/xpla.go/core/distribution"
	"github.com/Moonyongjung/xpla.go/util"

	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// Query client for distribution module.
func (i IXplaClient) QueryDistribution() (string, error) {
	queryClient := disttypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Distribution params
	case i.Ixplac.GetMsgType() == mdist.DistributionQueryDistributionParamsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Distribution validator outstanding rewards
	case i.Ixplac.GetMsgType() == mdist.DistributionValidatorOutstandingRewardsMSgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryValidatorOutstandingRewardsRequest)
		res, err = queryClient.ValidatorOutstandingRewards(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Distribution commission
	case i.Ixplac.GetMsgType() == mdist.DistributionQueryDistCommissionMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryValidatorCommissionRequest)
		res, err = queryClient.ValidatorCommission(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Distribution slashes
	case i.Ixplac.GetMsgType() == mdist.DistributionQuerySlashesMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryValidatorSlashesRequest)
		res, err = queryClient.ValidatorSlashes(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Distribution rewards
	case i.Ixplac.GetMsgType() == mdist.DistributionQueryRewardsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryDelegationRewardsRequest)
		res, err = queryClient.DelegationRewards(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Distribution community pool
	case i.Ixplac.GetMsgType() == mdist.DistributionQueryCommunityPoolMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryCommunityPoolRequest)
		res, err = queryClient.CommunityPool(
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
