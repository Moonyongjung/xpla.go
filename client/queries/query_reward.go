package queries

import (
	mreward "github.com/Moonyongjung/xpla.go/core/reward"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	rewardtypes "github.com/xpladev/xpla/x/reward/types"
)

// Query client for reward module.
func (i IXplaClient) QueryReward() (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcReward(i)
	} else {
		return queryByLcdReward(i)
	}
}

func queryByGrpcReward(i IXplaClient) (string, error) {
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
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Reward pool
	case i.Ixplac.GetMsgType() == mreward.RewardQueryRewardPoolMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(rewardtypes.QueryPoolRequest)
		res, err = queryClient.Pool(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err = printProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

const (
	rewardParamsLabel = "params"
	rewardPoolLabel   = "pool"
)

func queryByLcdReward(i IXplaClient) (string, error) {
	url := "/xpla/reward/v1beta1/"

	switch {
	// Reward params
	case i.Ixplac.GetMsgType() == mreward.RewardQueryRewardParamsMsgType:
		url = url + rewardParamsLabel

	// Reward pool
	case i.Ixplac.GetMsgType() == mreward.RewardQueryRewardPoolMsgType:
		url = url + rewardPoolLabel

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil

}
