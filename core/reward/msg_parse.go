package reward

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	rewardtypes "github.com/xpladev/xpla/x/reward/types"
)

// parsing - fund fee collector
func parseFundFeeCollectorArgs(fundFeeCollectorMsg types.FundFeeCollectorMsg, privKey key.PrivateKey) (*rewardtypes.MsgFundFeeCollector, error) {
	addrByPrivKey, err := key.Bech32AddrString(privKey)
	if err != nil {
		return nil, err
	}

	if fundFeeCollectorMsg.DepositorAddr != addrByPrivKey {
		return nil, util.LogErr("wrong depositor address, not match private key")
	}

	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(fundFeeCollectorMsg.Amount))
	if err != nil {
		return nil, err
	}

	addr, err := sdk.AccAddressFromBech32(fundFeeCollectorMsg.DepositorAddr)
	if err != nil {
		return nil, err
	}

	msg := rewardtypes.NewMsgFundFeeCollector(amount, addr)

	return msg, nil
}

// parsing - query reward params
func parseQueryRewardParamsArgs() (rewardtypes.QueryParamsRequest, error) {
	return rewardtypes.QueryParamsRequest{}, nil
}

// parsing - query reward pool
func parseQueryRewardPoolArgs() (rewardtypes.QueryPoolRequest, error) {
	return rewardtypes.QueryPoolRequest{}, nil
}
