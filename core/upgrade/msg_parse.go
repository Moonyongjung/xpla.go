package upgrade

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// Parsing - software upgrade
func parseProposalSoftwareUpgradeArgs(softwareUpgradeMsg types.SoftwareUpgradeMsg, privKey key.PrivateKey) (govtypes.MsgSubmitProposal, error) {
	plan := upgradetypes.Plan{
		Name:   softwareUpgradeMsg.UpgradeName,
		Height: util.FromStringToInt64(softwareUpgradeMsg.UpgradeHeight),
		Info:   softwareUpgradeMsg.UpgradeInfo,
	}
	content := upgradetypes.NewSoftwareUpgradeProposal(
		softwareUpgradeMsg.Title,
		softwareUpgradeMsg.Description,
		plan,
	)
	from, err := util.GetAddrByPrivKey(privKey)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	deposit, err := sdk.ParseCoinsNormalized(util.DenomAdd(softwareUpgradeMsg.Deposit))
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	return *msg, nil
}

// Parsing - cancel software upgrade
func parseCancelSoftwareUpgradeArgs(cancelSoftwareUpgradeMsg types.CancelSoftwareUpgradeMsg, privKey key.PrivateKey) (govtypes.MsgSubmitProposal, error) {
	from, err := util.GetAddrByPrivKey(privKey)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	deposit, err := sdk.ParseCoinsNormalized(util.DenomAdd(cancelSoftwareUpgradeMsg.Deposit))
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}
	content := upgradetypes.NewCancelSoftwareUpgradeProposal(
		cancelSoftwareUpgradeMsg.Deposit,
		cancelSoftwareUpgradeMsg.Description,
	)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	return *msg, nil
}
