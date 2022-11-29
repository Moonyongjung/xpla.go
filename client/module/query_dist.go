package module

import (
	mdist "github.com/Moonyongjung/xpla.go/core/distribution"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	distv1beta1 "cosmossdk.io/api/cosmos/distribution/v1beta1"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// Query client for distribution module.
func (i IXplaClient) QueryDistribution() (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcDist(i)
	} else {
		return queryByLcdDist(i)
	}
}

func queryByGrpcDist(i IXplaClient) (string, error) {
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

	// Distribution total rewards
	case i.Ixplac.GetMsgType() == mdist.DistributionQueryTotalRewardsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryDelegationTotalRewardsRequest)
		res, err = queryClient.DelegationTotalRewards(
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

const (
	distParamsLabel             = "params"
	distValidatorLabel          = "validators"
	distDelegatorLabel          = "delegators"
	distOutstandingRewardsLabel = "outstanding_rewards"
	distCommissionLabel         = "commission"
	distSlashesLabel            = "slashes"
	distRewardsLabel            = "rewards"
	distCommunityPoolLabel      = "community_pool"
)

func queryByLcdDist(i IXplaClient) (string, error) {
	url := util.MakeQueryLcdUrl(distv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Distribution params
	case i.Ixplac.GetMsgType() == mdist.DistributionQueryDistributionParamsMsgType:
		url = url + distParamsLabel

	// Distribution validator outstanding rewards
	case i.Ixplac.GetMsgType() == mdist.DistributionValidatorOutstandingRewardsMSgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryValidatorOutstandingRewardsRequest)

		url = url + util.MakeQueryLabels(distValidatorLabel, convertMsg.ValidatorAddress, distOutstandingRewardsLabel)

	// Distribution commission
	case i.Ixplac.GetMsgType() == mdist.DistributionQueryDistCommissionMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryValidatorCommissionRequest)

		url = url + util.MakeQueryLabels(distValidatorLabel, convertMsg.ValidatorAddress, distCommissionLabel)

	// Distribution slashes
	case i.Ixplac.GetMsgType() == mdist.DistributionQuerySlashesMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryValidatorSlashesRequest)

		url = url + util.MakeQueryLabels(distValidatorLabel, convertMsg.ValidatorAddress, distSlashesLabel)

	// Distribution rewards
	case i.Ixplac.GetMsgType() == mdist.DistributionQueryRewardsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryDelegationRewardsRequest)

		url = url + util.MakeQueryLabels(distDelegatorLabel, convertMsg.DelegatorAddress, distRewardsLabel, convertMsg.ValidatorAddress)

	// Distribution total rewards
	case i.Ixplac.GetMsgType() == mdist.DistributionQueryTotalRewardsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryDelegationTotalRewardsRequest)

		url = url + util.MakeQueryLabels(distDelegatorLabel, convertMsg.DelegatorAddress, distRewardsLabel)

	// Distribution community pool
	case i.Ixplac.GetMsgType() == mdist.DistributionQueryCommunityPoolMsgType:
		url = url + distCommunityPoolLabel

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err := CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil
}
