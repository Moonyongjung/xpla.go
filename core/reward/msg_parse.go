package reward

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	rewardtypes "github.com/xpladev/xpla/x/reward/types"
)

// parsing - fund fee collector
func parseFundFeeCollectorArgs(fundFeeCollectorMsg types.FundFeeCollectorMsg, privKey key.PrivateKey) (rewardtypes.MsgFundFeeCollector, error) {
	addrByPrivKey, err := key.Bech32AddrString(privKey)
	if err != nil {
		return rewardtypes.MsgFundFeeCollector{}, err
	}

	if fundFeeCollectorMsg.DepositorAddr != addrByPrivKey {
		return rewardtypes.MsgFundFeeCollector{}, util.LogErr(errors.ErrAccountNotMatch, "wrong depositor address, not match private key")
	}

	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(fundFeeCollectorMsg.Amount))
	if err != nil {
		return rewardtypes.MsgFundFeeCollector{}, err
	}

	addr, err := sdk.AccAddressFromBech32(fundFeeCollectorMsg.DepositorAddr)
	if err != nil {
		return rewardtypes.MsgFundFeeCollector{}, err
	}

	msg := rewardtypes.NewMsgFundFeeCollector(amount, addr)

	return *msg, nil
}
