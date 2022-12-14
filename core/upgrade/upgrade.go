package upgrade

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// (Tx) make msg - software upgrade
func MakeProposalSoftwareUpgradeMsg(softwareUpgradeMsg types.SoftwareUpgradeMsg, privKey key.PrivateKey) (govtypes.MsgSubmitProposal, error) {
	msg, err := parseProposalSoftwareUpgradeArgs(softwareUpgradeMsg, privKey)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, err
	}

	return msg, nil
}

// (Tx) make msg - cancel software upgrade
func MakeCancelSoftwareUpgradeMsg(cancelSoftwareUpgradeMsg types.CancelSoftwareUpgradeMsg, privKey key.PrivateKey) (govtypes.MsgSubmitProposal, error) {
	msg, err := parseCancelSoftwareUpgradeArgs(cancelSoftwareUpgradeMsg, privKey)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, err
	}

	return msg, nil
}

// (Query) make msg - applied
func MakeAppliedMsg(appliedMsg types.AppliedMsg) (upgradetypes.QueryAppliedPlanRequest, error) {
	return upgradetypes.QueryAppliedPlanRequest{
		Name: appliedMsg.UpgradeName,
	}, nil
}

// (Query) make msg - module version
func MakeQueryModuleVersionMsg(queryModulesVersionMsg types.QueryModulesVersionMsg) (upgradetypes.QueryModuleVersionsRequest, error) {
	return upgradetypes.QueryModuleVersionsRequest{
		ModuleName: queryModulesVersionMsg.ModuleName,
	}, nil
}

// (Query) make msg - all module versions
func MakeQueryAllModuleVersionMsg() (upgradetypes.QueryModuleVersionsRequest, error) {
	return upgradetypes.QueryModuleVersionsRequest{}, nil
}

// (Query) make msg - plan
func MakePlanMsg() (upgradetypes.QueryCurrentPlanRequest, error) {
	return upgradetypes.QueryCurrentPlanRequest{}, nil
}
