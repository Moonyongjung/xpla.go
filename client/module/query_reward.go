package module

import (
	mreward "github.com/Moonyongjung/xpla.go/core/reward"
	"github.com/Moonyongjung/xpla.go/util"

	rewardtypes "github.com/xpladev/xpla/x/reward/types"
)

// Query client for reward module.
func (i IXplaClient) QueryReward() (string, error) {
	queryClient := rewardtypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Reward params
	case i.Ixplac.GetMsgType() == mreward.RewardQueryRewardParamsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(rewardtypes.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Reward pool
	case i.Ixplac.GetMsgType() == mreward.RewardQueryRewardPoolMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(rewardtypes.QueryPoolRequest)
		res, err = queryClient.Pool(
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
