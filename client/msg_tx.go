package client

import (
	mauthz "github.com/Moonyongjung/xpla.go/core/authz"
	mbank "github.com/Moonyongjung/xpla.go/core/bank"
	mcrisis "github.com/Moonyongjung/xpla.go/core/crisis"
	mdist "github.com/Moonyongjung/xpla.go/core/distribution"
	mevm "github.com/Moonyongjung/xpla.go/core/evm"
	mfeegrant "github.com/Moonyongjung/xpla.go/core/feegrant"
	mgov "github.com/Moonyongjung/xpla.go/core/gov"
	mparams "github.com/Moonyongjung/xpla.go/core/params"
	mslashing "github.com/Moonyongjung/xpla.go/core/slashing"
	mstaking "github.com/Moonyongjung/xpla.go/core/staking"
	mupgrade "github.com/Moonyongjung/xpla.go/core/upgrade"
	mwasm "github.com/Moonyongjung/xpla.go/core/wasm"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
)

// Authz module

// Grant authorization to an address.
func (xplac *XplaClient) AuthzGrant(authzGrantMsg types.AuthzGrantMsg) *XplaClient {
	msg, err := mauthz.MakeAuthzGrantMsg(authzGrantMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mauthz.AuthzModule
	xplac.MsgType = mauthz.AuthzGrantMsgType
	xplac.Msg = msg
	return xplac
}

// Revoke authorization.
func (xplac *XplaClient) AuthzRevoke(authzRevokeMsg types.AuthzRevokeMsg) *XplaClient {
	msg, err := mauthz.MakeAuthzRevokeMsg(authzRevokeMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mauthz.AuthzModule
	xplac.MsgType = mauthz.AuthzRevokeMsgType
	xplac.Msg = msg
	return xplac
}

// Execute transaction on behalf of granter account.
func (xplac *XplaClient) AuthzExec(authzExecMsg types.AuthzExecMsg) *XplaClient {
	msg, err := mauthz.MakeAuthzExecMsg(authzExecMsg, xplac.EncodingConfig)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mauthz.AuthzModule
	xplac.MsgType = mauthz.AuthzExecMsgType
	xplac.Msg = msg
	return xplac
}

// Bank module

// Send funds from one account to another.
func (xplac *XplaClient) BankSend(bankSendMsg types.BankSendMsg) *XplaClient {
	msg, err := mbank.MakeBankSendMsg(bankSendMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mbank.BankModule
	xplac.MsgType = mbank.BankSendMsgType
	xplac.Msg = msg
	return xplac
}

// Crisis module

// Submit proof that an invariant broken to halt the chain.
func (xplac *XplaClient) InvariantBroken(invariantBrokenMsg types.InvariantBrokenMsg) *XplaClient {
	msg, err := mcrisis.MakeInvariantRouteMsg(invariantBrokenMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mcrisis.CrisisModule
	xplac.MsgType = mcrisis.CrisisInvariantBrokenMsgType
	xplac.Msg = msg
	return xplac
}

// Distribution module

// Funds the community pool with the specified amount.
func (xplac *XplaClient) FundCommunityPool(fundCommunityPoolMsg types.FundCommunityPoolMsg) *XplaClient {
	msg, err := mdist.MakeFundCommunityPoolMsg(fundCommunityPoolMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mdist.DistributionModule
	xplac.MsgType = mdist.DistributionFundCommunityPoolMsgType
	xplac.Msg = msg
	return xplac
}

// Submit a community pool spend proposal.
func (xplac *XplaClient) CommunityPoolSpend(communityPoolSpendMsg types.CommunityPoolSpendMsg) *XplaClient {
	msg, err := mdist.MakeProposalCommunityPoolSpendMsg(communityPoolSpendMsg, xplac.PrivateKey, xplac.EncodingConfig)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mdist.DistributionModule
	xplac.MsgType = mdist.DistributionProposalCommunityPoolSpendMsgType
	xplac.Msg = msg
	return xplac
}

// Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator.
func (xplac *XplaClient) WithdrawRewards(withdrawRewardsMsg types.WithdrawRewardsMsg) *XplaClient {
	msg, err := mdist.MakeWithdrawRewardsMsg(withdrawRewardsMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mdist.DistributionModule
	xplac.MsgType = mdist.DistributionWithdrawRewardsMsgType
	xplac.Msg = msg
	return xplac
}

// Withdraw all delegations rewards for a delegator.
func (xplac *XplaClient) WithdrawAllRewards() *XplaClient {
	msg, err := mdist.MakeWithdrawAllRewardsMsg(xplac.PrivateKey, xplac.Grpc, xplac.Context)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mdist.DistributionModule
	xplac.MsgType = mdist.DistributionWithdrawRewardsMsgType
	xplac.Msg = msg
	return xplac
}

// Change the default withdraw address for rewards associated with an address.
func (xplac *XplaClient) SetWithdrawAddr(setWithdrawAddrMsg types.SetwithdrawAddrMsg) *XplaClient {
	msg, err := mdist.MakeSetWithdrawAddrMsg(setWithdrawAddrMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mdist.DistributionModule
	xplac.MsgType = mdist.DistributionSetWithdrawAddrMsgType
	xplac.Msg = msg
	return xplac
}

// EVM module

// Send coind by using evm client.
func (xplac *XplaClient) EvmSendCoin(sendCoinMsg types.SendCoinMsg) *XplaClient {
	msg, err := mevm.MakeSendCoinMsg(sendCoinMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mevm.EvmModule
	xplac.MsgType = mevm.EvmSendCoinMsgType
	xplac.Msg = msg
	return xplac
}

// Deploy soldity contract.
func (xplac *XplaClient) DeploySolidityContract(deploySolContractMsg types.DeploySolContractMsg) *XplaClient {
	err := mevm.MakeDeploySolContractMsg(deploySolContractMsg)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mevm.EvmModule
	xplac.MsgType = mevm.EvmDeploySolContractMsgType
	xplac.Msg = nil
	return xplac
}

// Invoke (as execute) solidity contract.
func (xplac *XplaClient) InvokeSolidityContract(invokeSolContractMsg types.InvokeSolContractMsg) *XplaClient {
	msg, err := mevm.MakeInvokeSolContractMsg(invokeSolContractMsg)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mevm.EvmModule
	xplac.MsgType = mevm.EvmInvokeSolContractMsgType
	xplac.Msg = msg
	return xplac
}

// Feegrant module

// Grant fee allowance to an address.
func (xplac *XplaClient) Grant(grantMsg types.GrantMsg) *XplaClient {
	msg, err := mfeegrant.MakeGrantMsg(grantMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mfeegrant.FeegrantModule
	xplac.MsgType = mfeegrant.FeegrantGrantMsgType
	xplac.Msg = msg
	return xplac
}

// Revoke fee-grant.
func (xplac *XplaClient) RevokeGrant(revokeGrantMsg types.RevokeGrantMsg) *XplaClient {
	msg, err := mfeegrant.MakeRevokeGrantMsg(revokeGrantMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mfeegrant.FeegrantModule
	xplac.MsgType = mfeegrant.FeegrantRevokeGrantMsgType
	xplac.Msg = msg
	return xplac
}

// gov module

// Submit a proposal along with an initial deposit.
func (xplac *XplaClient) SubmitProposal(submitProposalMsg types.SubmitProposalMsg) *XplaClient {
	msg, err := mgov.MakeSubmitProposalMsg(submitProposalMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mgov.GovModule
	xplac.MsgType = mgov.GovSubmitProposalMsgType
	xplac.Msg = msg
	return xplac
}

// Deposit tokens for an active proposal.
func (xplac *XplaClient) GovDeposit(govDepositMsg types.GovDepositMsg) *XplaClient {
	msg, err := mgov.MakeGovDepositMsg(govDepositMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mgov.GovModule
	xplac.MsgType = mgov.GovDepositMsgType
	xplac.Msg = msg
	return xplac
}

// Vote for an active proposal, options: yes/no/no_with_veto/abstain.
func (xplac *XplaClient) Vote(voteMsg types.VoteMsg) *XplaClient {
	msg, err := mgov.MakeVoteMsg(voteMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mgov.GovModule
	xplac.MsgType = mgov.GovVoteMsgType
	xplac.Msg = msg
	return xplac
}

// Vote for an active proposal, options: yes/no/no_with_veto/abstain.
func (xplac *XplaClient) WeightedVote(weightedVoteMsg types.WeightedVoteMsg) *XplaClient {
	msg, err := mgov.MakeWeightedVoteMsg(weightedVoteMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mgov.GovModule
	xplac.MsgType = mgov.GovWeightedVoteMsgType
	xplac.Msg = msg
	return xplac
}

// Params module

// Submit a parameter change proposal.
func (xplac *XplaClient) ParamChange(paramChangeMsg types.ParamChangeMsg) *XplaClient {
	msg, err := mparams.MakeProposalParamChangeMsg(paramChangeMsg, xplac.PrivateKey, xplac.EncodingConfig)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mparams.ParamsModule
	xplac.MsgType = mparams.ParamsProposalParamChangeMsgType
	xplac.Msg = msg
	return xplac
}

// Slashing module

// Unjail validator previously jailed for downtime.
func (xplac *XplaClient) Unjail() *XplaClient {
	msg, err := mslashing.MakeUnjailMsg(xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mslashing.SlashingModule
	xplac.MsgType = mslashing.SlahsingUnjailMsgType
	xplac.Msg = msg
	return xplac
}

// Staking module

// Create new validator initialized with a self-delegation to it.
func (xplac *XplaClient) CreateValidator(createValidatorMsg types.CreateValidatorMsg) *XplaClient {
	msg, err := mstaking.MakeCreateValidatorMsg(createValidatorMsg, xplac.PrivateKey, xplac.Opts.OutputDocument)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mstaking.StakingModule
	xplac.MsgType = mstaking.StakingCreateValidatorMsgType
	xplac.Msg = msg
	return xplac
}

// Edit an existing validator account.
func (xplac *XplaClient) EditValidator(editValidatorMsg types.EditValidatorMsg) *XplaClient {
	msg, err := mstaking.MakeEditValidatorMsg(editValidatorMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mstaking.StakingModule
	xplac.MsgType = mstaking.StakingEditValidatorMsgType
	xplac.Msg = msg
	return xplac
}

// Delegate liquid tokens to a validator.
func (xplac *XplaClient) Delegate(delegateMsg types.DelegateMsg) *XplaClient {
	msg, err := mstaking.MakeDelegateMsg(delegateMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mstaking.StakingModule
	xplac.MsgType = mstaking.StakingDelegateMsgType
	xplac.Msg = msg
	return xplac
}

// Unbond shares from a validator.
func (xplac *XplaClient) Unbond(unbondMsg types.UnbondMsg) *XplaClient {
	msg, err := mstaking.MakeUnbondMsg(unbondMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mstaking.StakingModule
	xplac.MsgType = mstaking.StakingUnbondMsgType
	xplac.Msg = msg
	return xplac
}

// Redelegate illiquid tokens from one validator to another.
func (xplac *XplaClient) Redelegate(redelegateMsg types.RedelegateMsg) *XplaClient {
	msg, err := mstaking.MakeRedelegateMsg(redelegateMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mstaking.StakingModule
	xplac.MsgType = mstaking.StakingRedelegateMsgType
	xplac.Msg = msg
	return xplac
}

// Upgrade module

// Submit a software upgrade proposal.
func (xplac *XplaClient) SoftwareUpgrade(softwareUpgradeMsg types.SoftwareUpgradeMsg) *XplaClient {
	msg, err := mupgrade.MakeProposalSoftwareUpgradeMsg(softwareUpgradeMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mupgrade.UpgradeModule
	xplac.MsgType = mupgrade.UpgradeProposalSoftwareUpgradeMsgType
	xplac.Msg = msg
	return xplac
}

// Cancel the current software upgrade proposal.
func (xplac *XplaClient) CancelSoftwareUpgrade(cancelSoftwareUpgradeMsg types.CancelSoftwareUpgradeMsg) *XplaClient {
	msg, err := mupgrade.MakeCancelSoftwareUpgradeMsg(cancelSoftwareUpgradeMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mupgrade.UpgradeModule
	xplac.MsgType = mupgrade.UpgradeCancelSoftwareUpgradeMsgType
	xplac.Msg = msg
	return xplac
}

// Wasm module

// Upload a wasm binary.
func (xplac *XplaClient) StoreCode(filePath string) *XplaClient {
	addr := util.GetAddrByPrivKey(xplac.PrivateKey)
	msg, err := mwasm.MakeStoreCodeMsg(filePath, addr)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmStoreMsgType
	xplac.Msg = msg
	return xplac
}

// Instantiate a wasm contract.
func (xplac *XplaClient) InstantiateContract(instantiageMsg types.InstantiateMsg) *XplaClient {
	addr := util.GetAddrByPrivKey(xplac.PrivateKey)
	msg, err := mwasm.MakeInstantiateMsg(instantiageMsg, addr)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmInstantiateMsgType
	xplac.Msg = msg
	return xplac
}

// Execute a wasm contract.
func (xplac *XplaClient) ExecuteContract(executeMsg types.ExecuteMsg) *XplaClient {
	addr := util.GetAddrByPrivKey(xplac.PrivateKey)
	msg, err := mwasm.MakeExecuteMsg(executeMsg, addr)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmExecuteMsgType
	xplac.Msg = msg
	return xplac
}

// Clears admin for a contract to prevent further migrations.
func (xplac *XplaClient) ClearContractAdmin(clearContractAdminMsg types.ClearContractAdminMsg) *XplaClient {
	msg, err := mwasm.MakeClearContractAdminMsg(clearContractAdminMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmClearContractAdminMsgType
	xplac.Msg = msg
	return xplac
}

// Set new admin for a contract.
func (xplac *XplaClient) SetContractAdmin(setContractAdminMsg types.SetContractAdminMsg) *XplaClient {
	msg, err := mwasm.MakeSetContractAdmintMsg(setContractAdminMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmSetContractAdminMsgType
	xplac.Msg = msg
	return xplac
}

// Migrate a wasm contract to a new code version.
func (xplac *XplaClient) Migrate(migrateMsg types.MigrateMsg) *XplaClient {
	msg, err := mwasm.MakeMigrateMsg(migrateMsg, xplac.PrivateKey)
	if err != nil {
		xplac.Err = err
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmMigrateMsgType
	xplac.Msg = msg
	return xplac
}
