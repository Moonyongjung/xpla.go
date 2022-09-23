package client

import (
	mdist "github.com/Moonyongjung/xpla.go/core/distribution"
	"github.com/Moonyongjung/xpla.go/util"

	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// Query client for distribution module.
func queryDistribution(xplac *XplaClient) (string, error) {
	queryClient := disttypes.NewQueryClient(xplac.Grpc)

	switch {
	// Distribution params
	case xplac.MsgType == mdist.DistributionQueryDistributionParamsMsgType:
		convertMsg, _ := xplac.Msg.(disttypes.QueryParamsRequest)
		res, err = queryClient.Params(
			xplac.Context,
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Distribution validator outstanding rewards
	case xplac.MsgType == mdist.DistributionValidatorOutstandingRewardsMSgType:
		convertMsg, _ := xplac.Msg.(disttypes.QueryValidatorOutstandingRewardsRequest)
		res, err = queryClient.ValidatorOutstandingRewards(
			xplac.Context,
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Distribution commission
	case xplac.MsgType == mdist.DistributionQueryDistCommissionMsgType:
		convertMsg, _ := xplac.Msg.(disttypes.QueryValidatorCommissionRequest)
		res, err = queryClient.ValidatorCommission(
			xplac.Context,
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Distribution slashes
	case xplac.MsgType == mdist.DistributionQuerySlashesMsgType:
		convertMsg, _ := xplac.Msg.(disttypes.QueryValidatorSlashesRequest)
		res, err = queryClient.ValidatorSlashes(
			xplac.Context,
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Distribution rewards
	case xplac.MsgType == mdist.DistributionQueryRewardsMsgType:
		convertMsg, _ := xplac.Msg.(disttypes.QueryDelegationRewardsRequest)
		res, err = queryClient.DelegationRewards(
			xplac.Context,
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Distribution community pool
	case xplac.MsgType == mdist.DistributionQueryCommunityPoolMsgType:
		convertMsg, _ := xplac.Msg.(disttypes.QueryCommunityPoolRequest)
		res, err = queryClient.CommunityPool(
			xplac.Context,
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
