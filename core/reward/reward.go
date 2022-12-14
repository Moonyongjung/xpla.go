package reward

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"

	rewardtypes "github.com/xpladev/xpla/x/reward/types"
)

// (Tx) make msg - Fund fee collector
func MakeFundFeeCollectorMsg(fundFeeCollectorMsg types.FundFeeCollectorMsg, privKey key.PrivateKey) (rewardtypes.MsgFundFeeCollector, error) {
	msg, err := parseFundFeeCollectorArgs(fundFeeCollectorMsg, privKey)
	if err != nil {
		return rewardtypes.MsgFundFeeCollector{}, err
	}

	return msg, nil
}

// (Query) make msg - query reward params
func MakeQueryRewardParamsMsg() (rewardtypes.QueryParamsRequest, error) {
	return rewardtypes.QueryParamsRequest{}, nil
}

// (Query) make msg - query reward pool
func MakeQueryRewardPoolMsg() (rewardtypes.QueryPoolRequest, error) {
	return rewardtypes.QueryPoolRequest{}, nil
}
