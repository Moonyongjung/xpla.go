package client

import (
	mauth "github.com/Moonyongjung/xpla.go/core/auth"
	mauthz "github.com/Moonyongjung/xpla.go/core/authz"
	mbank "github.com/Moonyongjung/xpla.go/core/bank"
	mdist "github.com/Moonyongjung/xpla.go/core/distribution"
	mevidence "github.com/Moonyongjung/xpla.go/core/evidence"
	mevm "github.com/Moonyongjung/xpla.go/core/evm"
	mfeegrant "github.com/Moonyongjung/xpla.go/core/feegrant"
	mgov "github.com/Moonyongjung/xpla.go/core/gov"
	mmint "github.com/Moonyongjung/xpla.go/core/mint"
	mparams "github.com/Moonyongjung/xpla.go/core/params"
	mreward "github.com/Moonyongjung/xpla.go/core/reward"
	mslashing "github.com/Moonyongjung/xpla.go/core/slashing"
	mstaking "github.com/Moonyongjung/xpla.go/core/staking"
	mupgrade "github.com/Moonyongjung/xpla.go/core/upgrade"
	mwasm "github.com/Moonyongjung/xpla.go/core/wasm"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
)

// Auth module

// Query the current auth parameters.
func (xplac *XplaClient) AuthParams() *XplaClient {
	msg, err := mauth.MakeAuthParamMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mauth.AuthModule
	xplac.MsgType = mauth.AuthQueryParamsMsgType
	xplac.Msg = msg
	return xplac
}

// Query for account by address.
func (xplac *XplaClient) AccAddress(queryAccAddresMsg types.QueryAccAddressMsg) *XplaClient {
	msg, err := mauth.MakeQueryAccAddressMsg(queryAccAddresMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mauth.AuthModule
	xplac.MsgType = mauth.AuthQueryAccAddressMsgType
	xplac.Msg = msg
	return xplac
}

// Query all accounts.
func (xplac *XplaClient) Accounts() *XplaClient {
	msg, err := mauth.MakeQueryAccountsMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mauth.AuthModule
	xplac.MsgType = mauth.AuthQueryAccountsMsgType
	xplac.Msg = msg
	return xplac
}

// Query for paginated transactions that match a set of events.
func (xplac *XplaClient) TxsByEvents(txsByEventsMsg types.QueryTxsByEventsMsg) *XplaClient {
	msg, err := mauth.MakeTxsByEventsMsg(txsByEventsMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mauth.AuthModule
	xplac.MsgType = mauth.AuthQueryTxsByEventsMsgType
	xplac.Msg = msg
	return xplac
}

// Query for a transaction by hash <addr>/<seq> combination or comma-separated signatures in a committed block.
func (xplac *XplaClient) Tx(queryTxMsg types.QueryTxMsg) *XplaClient {
	msg, err := mauth.MakeQueryTxMsg(queryTxMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mauth.AuthModule
	xplac.MsgType = mauth.AuthQueryTxMsgType
	xplac.Msg = msg
	return xplac
}

// Authz module

// Query grants for granter-grantee pair and optionally a msg-type-url.
// Also, it is able to support querying grants granted by granter and granted to a grantee.
func (xplac *XplaClient) QueryAuthzGrants(queryAuthzGrantMsg types.QueryAuthzGrantMsg) *XplaClient {
	if queryAuthzGrantMsg.Grantee != "" && queryAuthzGrantMsg.Granter != "" {
		msg, err := mauthz.MakeQueryAuthzGrantsMsg(queryAuthzGrantMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mauthz.AuthzModule
		xplac.MsgType = mauthz.AuthzQueryGrantMsgType
		xplac.Msg = msg
	} else if queryAuthzGrantMsg.Grantee != "" && queryAuthzGrantMsg.Granter == "" {
		msg, err := mauthz.MakeQueryAuthzGrantsByGranteeMsg(queryAuthzGrantMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mauthz.AuthzModule
		xplac.MsgType = mauthz.AuthzQueryGrantsByGranteeMsgType
		xplac.Msg = msg
	} else if queryAuthzGrantMsg.Grantee == "" && queryAuthzGrantMsg.Granter != "" {
		msg, err := mauthz.MakeQueryAuthzGrantsByGranterMsg(queryAuthzGrantMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mauthz.AuthzModule
		xplac.MsgType = mauthz.AuthzQueryGrantsByGranterMsgType
		xplac.Msg = msg
	} else {
		xplac.Err = util.LogErr("No query grants parameters")
	}
	return xplac
}

// Bank module

// Query for account balances by address
func (xplac *XplaClient) BankBalances(bankBalancesMsg types.BankBalancesMsg) *XplaClient {
	if bankBalancesMsg.Denom == "" {
		msg, err := mbank.MakeBankAllBalancesMsg(bankBalancesMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mbank.BankModule
		xplac.MsgType = mbank.BankAllBalancesMsgType
		xplac.Msg = msg
	} else {
		msg, err := mbank.MakeBankBalanceMsg(bankBalancesMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mbank.BankModule
		xplac.MsgType = mbank.BankBalanceMsgType
		xplac.Msg = msg
	}
	return xplac

}

// Query the client metadata for coin denominations.
func (xplac *XplaClient) DenomMetadata(denomMetadataMsg ...types.DenomMetadataMsg) *XplaClient {
	if len(denomMetadataMsg) == 0 {
		msg, err := mbank.MakeDenomsMetaDataMsg()
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mbank.BankModule
		xplac.MsgType = mbank.BankDenomsMetadataMsgType
		xplac.Msg = msg
	} else if len(denomMetadataMsg) == 1 {
		msg, err := mbank.MakeDenomMetaDataMsg(denomMetadataMsg[0])
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mbank.BankModule
		xplac.MsgType = mbank.BankDenomMetadataMsgType
		xplac.Msg = msg
	} else {
		xplac.Err = util.LogErr("Need one parameter")
	}
	return xplac
}

// Query the total supply of coins of the chain.
func (xplac *XplaClient) Total(totalMsg ...types.TotalMsg) *XplaClient {
	if len(totalMsg) == 0 {
		msg, err := mbank.MakeTotalSupplyMsg()
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mbank.BankModule
		xplac.MsgType = mbank.BankTotalMsgType
		xplac.Msg = msg
	} else if len(totalMsg) == 1 {
		msg, err := mbank.MakeSupplyOfMsg(totalMsg[0])
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mbank.BankModule
		xplac.MsgType = mbank.BankTotalSupplyOfMsgType
		xplac.Msg = msg
	} else {
		xplac.Err = util.LogErr("Need one parameter")
	}
	return xplac
}

// Distribution module

// Query distribution parameters.
func (xplac *XplaClient) DistributionParams() *XplaClient {
	msg, err := mdist.MakeQueryDistributionParamsMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mdist.DistributionModule
	xplac.MsgType = mdist.DistributionQueryDistributionParamsMsgType
	xplac.Msg = msg
	return xplac
}

// Query distribution outstanding (un-withdrawn) rewards for a validator and all thier delegations.
func (xplac *XplaClient) ValidatorOutstandingRewards(validatorOutstandingRewardsMsg types.ValidatorOutstandingRewardsMsg) *XplaClient {
	msg, err := mdist.MakeValidatorOutstandingRewardsMsg(validatorOutstandingRewardsMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mdist.DistributionModule
	xplac.MsgType = mdist.DistributionValidatorOutstandingRewardsMSgType
	xplac.Msg = msg
	return xplac
}

// Query distribution validator commission.
func (xplac *XplaClient) DistCommission(queryDistCommissionMsg types.QueryDistCommissionMsg) *XplaClient {
	msg, err := mdist.MakeQueryDistCommissionMsg(queryDistCommissionMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mdist.DistributionModule
	xplac.MsgType = mdist.DistributionQueryDistCommissionMsgType
	xplac.Msg = msg
	return xplac
}

// Query distribution validator slashes.
func (xplac *XplaClient) DistSlashes(queryDistSlashesMsg types.QueryDistSlashesMsg) *XplaClient {
	msg, err := mdist.MakeQueryDistSlashesMsg(queryDistSlashesMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mdist.DistributionModule
	xplac.MsgType = mdist.DistributionQuerySlashesMsgType
	xplac.Msg = msg
	return xplac
}

// Query all ditribution delegator rewards or rewards from a particular validator.
func (xplac *XplaClient) DistRewards(queryDistRewardsMsg types.QueryDistRewardsMsg) *XplaClient {
	msg, err := mdist.MakeyQueryDistRewardsMsg(queryDistRewardsMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mdist.DistributionModule
	xplac.MsgType = mdist.DistributionQueryRewardsMsgType
	xplac.Msg = msg
	return xplac
}

// Query the amount of coins in the community pool.
func (xplac *XplaClient) CommunityPool() *XplaClient {
	msg, err := mdist.MakeQueryCommunityPoolMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mdist.DistributionModule
	xplac.MsgType = mdist.DistributionQueryCommunityPoolMsgType
	xplac.Msg = msg
	return xplac
}

// Evidence module

// Query for evidence by hash or for all (paginated) submitted evidence.
func (xplac *XplaClient) QueryEvidence(queryEvidenceMsg ...types.QueryEvidenceMsg) *XplaClient {
	if len(queryEvidenceMsg) == 0 {
		msg, err := mevidence.MakeQueryAllEvidenceMsg()
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mevidence.EvidenceModule
		xplac.MsgType = mevidence.EvidenceQueryAllMsgType
		xplac.Msg = msg
	} else if len(queryEvidenceMsg) == 1 {
		msg, err := mevidence.MakeQueryEvidenceMsg(queryEvidenceMsg[0])
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mevidence.EvidenceModule
		xplac.MsgType = mevidence.EvidenceQueryMsgType
		xplac.Msg = msg
	} else {
		xplac.Err = util.LogErr("Need one parameter")
	}
	return xplac
}

// EVM module

// Call(as query) solidity contract.
func (xplac *XplaClient) CallSolidityContract(callSolContractMsg types.CallSolContractMsg) *XplaClient {
	msg, err := mevm.MakeCallSolContractMsg(callSolContractMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mevm.EvmModule
	xplac.MsgType = mevm.EvmCallSolContractMsgType
	xplac.Msg = msg
	return xplac
}

// Query a transaction which is ethereum type information by retrieving hash.
func (xplac *XplaClient) GetTransactionByHash(getTransactionByHashMsg types.GetTransactionByHashMsg) *XplaClient {
	msg, err := mevm.MakeGetTransactionByHashMsg(getTransactionByHashMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mevm.EvmModule
	xplac.MsgType = mevm.EvmGetTransactionByHashMsgType
	xplac.Msg = msg
	return xplac
}

// Query a block which is ethereum type information by retrieving hash or block height(as number).
func (xplac *XplaClient) GetBlockByHashOrHeight(getBlockByHashHeightMsg types.GetBlockByHashHeightMsg) *XplaClient {
	msg, err := mevm.MakeGetBlockByHashHeightMsg(getBlockByHashHeightMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mevm.EvmModule
	xplac.MsgType = mevm.EvmGetBlockByHashHeightMsgType
	xplac.Msg = msg
	return xplac
}

// Query a account information which includes account address(hex and bech32), balance and etc.
func (xplac *XplaClient) AccountInfo(accountInfoMsg types.AccountInfoMsg) *XplaClient {
	msg, err := mevm.MakeQueryAccountInfoMsg(accountInfoMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mevm.EvmModule
	xplac.MsgType = mevm.EvmQueryAccountInfoMsgType
	xplac.Msg = msg
	return xplac
}

// Query suggested gas price.
func (xplac *XplaClient) SuggestGasPrice() *XplaClient {
	xplac.Module = mevm.EvmModule
	xplac.MsgType = mevm.EvmSuggestGasPriceMsgType
	xplac.Msg = nil
	return xplac
}

// Query chain ID of ethereum type.
func (xplac *XplaClient) EthChainID() *XplaClient {
	xplac.Module = mevm.EvmModule
	xplac.MsgType = mevm.EvmQueryChainIdMsgType
	xplac.Msg = nil
	return xplac
}

// Query latest block height(as number)
func (xplac *XplaClient) EthBlockNumber() *XplaClient {
	xplac.Module = mevm.EvmModule
	xplac.MsgType = mevm.EvmQueryCurrentBlockNumberMsgType
	xplac.Msg = nil
	return xplac
}

// Feegrant module

// Query details of fee grants.
func (xplac *XplaClient) QueryGrants(queryGrantMsg types.QueryGrantMsg) *XplaClient {
	if queryGrantMsg.Grantee != "" && queryGrantMsg.Granter != "" {
		msg, err := mfeegrant.MakeQueryGrantMsg(queryGrantMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mfeegrant.FeegrantModule
		xplac.MsgType = mfeegrant.FeegrantQueryGrantMsgType
		xplac.Msg = msg
	} else if queryGrantMsg.Grantee != "" && queryGrantMsg.Granter == "" {
		msg, err := mfeegrant.MakeQueryGrantsByGranteeMsg(queryGrantMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mfeegrant.FeegrantModule
		xplac.MsgType = mfeegrant.FeegrantQueryGrantsByGranteeMsgType
		xplac.Msg = msg
	} else if queryGrantMsg.Grantee == "" && queryGrantMsg.Granter != "" {
		msg, err := mfeegrant.MakeQueryGrantsByGranterMsg(queryGrantMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mfeegrant.FeegrantModule
		xplac.MsgType = mfeegrant.FeegrantQueryGrantsByGranterMsgType
		xplac.Msg = msg
	} else {
		xplac.Err = util.LogErr("No query grants parameters")
	}

	return xplac
}

// Gov module

// Query details of a singla proposal.
func (xplac *XplaClient) QueryProposal(queryProposal types.QueryProposalMsg) *XplaClient {
	msg, err := mgov.MakeQueryProposalMsg(queryProposal)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mgov.GovModule
	xplac.MsgType = mgov.GovQueryProposalMsgType
	xplac.Msg = msg
	return xplac
}

// Query proposals with optional filters.
func (xplac *XplaClient) QueryProposals(queryProposals types.QueryProposalsMsg) *XplaClient {
	msg, err := mgov.MakeQueryProposalsMsg(queryProposals)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mgov.GovModule
	xplac.MsgType = mgov.GovQueryProposalsMsgType
	xplac.Msg = msg
	return xplac
}

// Query details of a deposit or deposits on a proposal.
func (xplac *XplaClient) QueryDeposit(queryDepositMsg types.QueryDepositMsg) *XplaClient {
	if queryDepositMsg.Depositor != "" {
		msg, argsType, err := mgov.MakeQueryDepositMsg(queryDepositMsg, xplac.Grpc, xplac.Context)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		if argsType == "params" {
			xplac.Module = mgov.GovModule
			xplac.MsgType = mgov.GovQueryDepositParamsMsgType
			xplac.Msg = msg
		} else {
			xplac.Module = mgov.GovModule
			xplac.MsgType = mgov.GovQueryDepositRequestMsgType
			xplac.Msg = msg
		}
	} else {
		msg, argsType, err := mgov.MakeQueryDepositsMsg(queryDepositMsg, xplac.Grpc, xplac.Context)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		if argsType == "params" {
			xplac.Module = mgov.GovModule
			xplac.MsgType = mgov.GovQueryDepositsParamsMsgType
			xplac.Msg = msg
		} else {
			xplac.Module = mgov.GovModule
			xplac.MsgType = mgov.GovQueryDepositsRequestMsgType
			xplac.Msg = msg
		}
	}
	return xplac
}

// Query details of a single vote or votes on a proposal.
func (xplac *XplaClient) QueryVote(queryVoteMsg types.QueryVoteMsg) *XplaClient {
	if queryVoteMsg.VoterAddr != "" {
		msg, err := mgov.MakeQueryVoteMsg(queryVoteMsg, xplac.Grpc, xplac.Context)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mgov.GovModule
		xplac.MsgType = mgov.GovQueryVoteMsgType
		xplac.Msg = msg

	} else {
		msg, status, err := mgov.MakeQueryVotesMsg(queryVoteMsg, xplac.Grpc, xplac.Context)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		if status == "notPassed" {
			xplac.Module = mgov.GovModule
			xplac.MsgType = mgov.GovQueryVotesNotPassedMsgType
			xplac.Msg = msg
		} else {
			xplac.Module = mgov.GovModule
			xplac.MsgType = mgov.GovQueryVotesPassedMsgType
			xplac.Msg = msg
		}
	}
	return xplac
}

// Query the tally of a proposal vote.
func (xplac *XplaClient) Tally(tallyMsg types.TallyMsg) *XplaClient {
	msg, err := mgov.MakeGovTallyMsg(tallyMsg, xplac.Grpc, xplac.Context)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mgov.GovModule
	xplac.MsgType = mgov.GovTallyMsgType
	xplac.Msg = msg
	return xplac
}

// Query parameters of the governance process or the parameters (voting|tallying|deposit) of the governance process.
func (xplac *XplaClient) GovParams(govParamsMsg ...types.GovParamsMsg) *XplaClient {
	if len(govParamsMsg) == 0 {
		xplac.Module = mgov.GovModule
		xplac.MsgType = mgov.GovQueryGovParamsMsgType
		xplac.Msg = nil
	} else if len(govParamsMsg) == 1 {
		msg, err := mgov.MakeGovParamsMsg(govParamsMsg[0])
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mgov.GovModule
		switch govParamsMsg[0].ParamType {
		case "voting":
			xplac.MsgType = mgov.GovQueryGovParamVotingMsgType
		case "tallying":
			xplac.MsgType = mgov.GovQueryGovParamTallyingMsgType
		case "deposit":
			xplac.MsgType = mgov.GovQueryGovParamDepositMsgType
		}
		xplac.Msg = msg
	} else {
		xplac.Err = util.LogErr("Need one parameter")
	}
	return xplac
}

// Query the proposer of a governance proposal.
func (xplac *XplaClient) Proposer(proposerMsg types.ProposerMsg) *XplaClient {
	xplac.Module = mgov.GovModule
	xplac.MsgType = mgov.GovQueryProposerMsgType
	xplac.Msg = proposerMsg.ProposalID
	return xplac
}

// Mint module

// Query the current minting parameters.
func (xplac *XplaClient) MintParams() *XplaClient {
	msg, err := mmint.MakeQueryMintParamsMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mmint.MintModule
	xplac.MsgType = mmint.MintQueryMintParamsMsgType
	xplac.Msg = msg
	return xplac
}

// Query the current minting inflation value.
func (xplac *XplaClient) Inflation() *XplaClient {
	msg, err := mmint.MakeQueryInflationMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mmint.MintModule
	xplac.MsgType = mmint.MintQueryInflationMsgType
	xplac.Msg = msg
	return xplac
}

// Query the current minting annual provisions value.
func (xplac *XplaClient) AnnualProvisions() *XplaClient {
	msg, err := mmint.MakeQueryAnnualProvisionsMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mmint.MintModule
	xplac.MsgType = mmint.MintQueryAnnualProvisionsMsgType
	xplac.Msg = msg
	return xplac
}

// Params module

// Query for raw parameters by subspace and key.
func (xplac *XplaClient) QuerySubspace(subspaceMsg types.SubspaceMsg) *XplaClient {
	msg, err := mparams.MakeQueryParamsSubspaceMsg(subspaceMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mparams.ParamsModule
	xplac.MsgType = mparams.ParamsQuerySubpsaceMsgType
	xplac.Msg = msg

	return xplac
}

// Reward module

// Query reward params
func (xplac *XplaClient) RewardParams() *XplaClient {
	msg, err := mreward.MakeQueryRewardParamsMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mreward.RewardModule
	xplac.MsgType = mreward.RewardQueryRewardParamsMsgType
	xplac.Msg = msg

	return xplac
}

// Query reward pool
func (xplac *XplaClient) RewardPool() *XplaClient {
	msg, err := mreward.MakeQueryRewardPoolMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mreward.RewardModule
	xplac.MsgType = mreward.RewardQueryRewardPoolMsgType
	xplac.Msg = msg

	return xplac
}

// Slashing module

// Query the current slashing parameters.
func (xplac *XplaClient) SlashingParams() *XplaClient {
	msg, err := mslashing.MakeQuerySlashingParamsMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mslashing.SlashingModule
	xplac.MsgType = mslashing.SlahsingQuerySlashingParamsMsgType
	xplac.Msg = msg
	return xplac
}

// Query a validator's signing information or signing information of all validators.
func (xplac *XplaClient) SigningInfos(signingInfoMsg ...types.SigningInfoMsg) *XplaClient {
	if len(signingInfoMsg) == 0 {
		msg, err := mslashing.MakeQuerySigningInfosMsg()
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mslashing.SlashingModule
		xplac.MsgType = mslashing.SlashingQuerySigningInfosMsgType
		xplac.Msg = msg
	} else if len(signingInfoMsg) == 1 {
		msg, err := mslashing.MakeQuerySigningInfoMsg(signingInfoMsg[0], xplac.EncodingConfig)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mslashing.SlashingModule
		xplac.MsgType = mslashing.SlashingQuerySigningInfoMsgType
		xplac.Msg = msg
	} else {
		xplac.Err = util.LogErr("Need one parameter")
	}
	return xplac
}

// Staking module

// Query a validator or for all validators.
func (xplac *XplaClient) QueryValidators(queryValidatorMsg ...types.QueryValidatorMsg) *XplaClient {
	if len(queryValidatorMsg) == 0 {
		msg, err := mstaking.MakeQueryValidatorsMsg()
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mstaking.StakingModule
		xplac.MsgType = mstaking.StakingQueryValidatorsMsgType
		xplac.Msg = msg
	} else if len(queryValidatorMsg) == 1 {
		msg, err := mstaking.MakeQueryValidatorMsg(queryValidatorMsg[0])
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mstaking.StakingModule
		xplac.MsgType = mstaking.StakingQueryValidatorMsgType
		xplac.Msg = msg
	} else {
		xplac.Err = util.LogErr("Need one parameter")
	}
	return xplac
}

// Query a delegation based on address and validator address, all out going redelegations from a validator or all delegations made by on delegator.
func (xplac *XplaClient) QueryDelegation(queryDelegationMsg types.QueryDelegationMsg) *XplaClient {
	if queryDelegationMsg.DelegatorAddr != "" && queryDelegationMsg.ValidatorAddr != "" {
		msg, err := mstaking.MakeQueryDelegationMsg(queryDelegationMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mstaking.StakingModule
		xplac.MsgType = mstaking.StakingQueryDelegationMsgType
		xplac.Msg = msg
	} else if queryDelegationMsg.DelegatorAddr != "" {
		msg, err := mstaking.MakeQueryDelegationsMsg(queryDelegationMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mstaking.StakingModule
		xplac.MsgType = mstaking.StakingQueryDelegationsMsgType
		xplac.Msg = msg
	} else if queryDelegationMsg.ValidatorAddr != "" {
		msg, err := mstaking.MakeQueryDelegationsToMsg(queryDelegationMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mstaking.StakingModule
		xplac.MsgType = mstaking.StakingQueryDelegationsToMsgType
		xplac.Msg = msg
	} else {
		xplac.Err = util.LogErr("Wrong delegation msg")
	}
	return xplac
}

// Query all unbonding delegatations from a validator, an unbonding-delegation record based on delegator and validator address or all unbonding-delegations records for one delegator.
func (xplac *XplaClient) QueryUnbondingDelegation(queryUnbondingDelegationMsg types.QueryUnbondingDelegationMsg) *XplaClient {
	if queryUnbondingDelegationMsg.DelegatorAddr != "" && queryUnbondingDelegationMsg.ValidatorAddr != "" {
		msg, err := mstaking.MakeQueryUnbondingDelegationMsg(queryUnbondingDelegationMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mstaking.StakingModule
		xplac.MsgType = mstaking.StakingQueryUnbondingDelegationMsgType
		xplac.Msg = msg
	} else if queryUnbondingDelegationMsg.DelegatorAddr != "" {
		msg, err := mstaking.MakeQueryUnbondingDelegationsMsg(queryUnbondingDelegationMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mstaking.StakingModule
		xplac.MsgType = mstaking.StakingQueryUnbondingDelegationsMsgType
		xplac.Msg = msg
	} else if queryUnbondingDelegationMsg.ValidatorAddr != "" {
		msg, err := mstaking.MakeQueryUnbondingDelegationsFromMsg(queryUnbondingDelegationMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mstaking.StakingModule
		xplac.MsgType = mstaking.StakingQueryUnbondingDelegationsFromMsgType
		xplac.Msg = msg
	} else {
		xplac.Err = util.LogErr("Wrong unbonding delegation msg")
	}
	return xplac
}

// Query a redelegation record based on delegator and a source and destination validator.
// Also, query all outgoing redelegatations from a validator or all redelegations records for one delegator.
func (xplac *XplaClient) QueryRedelegation(queryRedelegationMsg types.QueryRedelegationMsg) *XplaClient {
	if queryRedelegationMsg.DelegatorAddr != "" &&
		queryRedelegationMsg.SrcValidatorAddr != "" &&
		queryRedelegationMsg.DstValidatorAddr != "" {
		msg, err := mstaking.MakeQueryRedelegationMsg(queryRedelegationMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mstaking.StakingModule
		xplac.MsgType = mstaking.StakingQueryRedelegationMsgType
		xplac.Msg = msg
	} else if queryRedelegationMsg.DelegatorAddr != "" {
		msg, err := mstaking.MakeQueryRedelegationsMsg(queryRedelegationMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mstaking.StakingModule
		xplac.MsgType = mstaking.StakingQueryRedelegationsMsgType
		xplac.Msg = msg
	} else if queryRedelegationMsg.SrcValidatorAddr != "" {
		msg, err := mstaking.MakeQueryRedelegationsFromMsg(queryRedelegationMsg)
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mstaking.StakingModule
		xplac.MsgType = mstaking.StakingQueryRedelegationsFromMsgType
		xplac.Msg = msg
	} else {
		xplac.Err = util.LogErr("Wrong redelegation msg")
	}
	return xplac
}

// Query historical info at given height.
func (xplac *XplaClient) HistoricalInfo(historicalInfoMsg types.HistoricalInfoMsg) *XplaClient {
	msg, err := mstaking.MakeHistoricalInfoMsg(historicalInfoMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mstaking.StakingModule
	xplac.MsgType = mstaking.StakingHistoricalInfoMsgType
	xplac.Msg = msg
	return xplac
}

// Query the current staking pool values.
func (xplac *XplaClient) StakingPool() *XplaClient {
	msg, err := mstaking.MakeQueryStakingPoolMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mstaking.StakingModule
	xplac.MsgType = mstaking.StakingQueryStakingPoolMsgType
	xplac.Msg = msg
	return xplac
}

// Query the current staking parameters information.
func (xplac *XplaClient) StakingParams() *XplaClient {
	msg, err := mstaking.MakeQueryStakingParamsMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mstaking.StakingModule
	xplac.MsgType = mstaking.StakingQueryStakingParamsMsgType
	xplac.Msg = msg
	return xplac
}

// Upgrade module

// Block header for height at which a completed upgrade was applied.
func (xplac *XplaClient) UpgradeApplied(appliedMsg types.AppliedMsg) *XplaClient {
	msg, err := mupgrade.MakeAppliedMsg(appliedMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mupgrade.UpgradeModule
	xplac.MsgType = mupgrade.UpgradeAppliedMsgType
	xplac.Msg = msg
	return xplac
}

// Query the list of module versions.
func (xplac *XplaClient) ModulesVersion(queryModulesVersionMsg ...types.QueryModulesVersionMsg) *XplaClient {
	if len(queryModulesVersionMsg) == 0 {
		msg, err := mupgrade.MakeQueryAllModuleVersionMsg()
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mupgrade.UpgradeModule
		xplac.MsgType = mupgrade.UpgradeQueryAllModuleVersionsMsgType
		xplac.Msg = msg
	} else if len(queryModulesVersionMsg) == 1 {
		msg, err := mupgrade.MakeQueryModuleVersionMsg(queryModulesVersionMsg[0])
		if err != nil {
			xplac.Err = err
			return xplac
		}
		xplac.Module = mupgrade.UpgradeModule
		xplac.MsgType = mupgrade.UpgradeQueryModuleVersionsMsgType
		xplac.Msg = msg
	} else {
		xplac.Err = util.LogErr("Need one parameter")
	}
	return xplac
}

// Query upgrade plan(if one exists).
func (xplac *XplaClient) Plan() *XplaClient {
	msg, err := mupgrade.MakePlanMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mupgrade.UpgradeModule
	xplac.MsgType = mupgrade.UpgradePlanMsgType
	xplac.Msg = msg
	return xplac
}

// Wasm module

// Calls contract with given address with query data and prints the returned result.
func (xplac *XplaClient) QueryContract(queryMsg types.QueryMsg) *XplaClient {
	addr := util.GetAddrByPrivKey(xplac.Opts.PrivateKey)
	msg, err := mwasm.MakeQueryMsg(queryMsg, addr)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmQueryContractMsgType
	xplac.Msg = msg
	return xplac
}

// Query list all wasm bytecode on the chain.
func (xplac *XplaClient) ListCode() *XplaClient {
	msg, err := mwasm.MakeListcodeMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmListCodeMsgType
	xplac.Msg = msg
	return xplac
}

// Query list wasm all bytecode on the chain for given code ID.
func (xplac *XplaClient) ListContractByCode(listContractByCodeMsg types.ListContractByCodeMsg) *XplaClient {
	msg, err := mwasm.MakeListContractByCodeMsg(listContractByCodeMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmListContractByCodeMsgType
	xplac.Msg = msg
	return xplac
}

// Downloads wasm bytecode for given code ID.
func (xplac *XplaClient) Download(downloadMsg types.DownloadMsg) *XplaClient {
	msg, err := mwasm.MakeDownloadMsg(downloadMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmDownloadMsgType
	xplac.Msg = msg
	return xplac
}

// Prints out metadata of a code ID.
func (xplac *XplaClient) CodeInfo(codeInfoMsg types.CodeInfoMsg) *XplaClient {
	msg, err := mwasm.MakeCodeInfoMsg(codeInfoMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmCodeInfoMsgType
	xplac.Msg = msg
	return xplac
}

// Prints out metadata of a contract given its address.
func (xplac *XplaClient) ContractInfo(contractInfoMsg types.ContractInfoMsg) *XplaClient {
	msg, err := mwasm.MakeContractInfoMsg(contractInfoMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmContractInfoMsgType
	xplac.Msg = msg
	return xplac
}

// Prints out all internal state of a contract given its address.
func (xplac *XplaClient) ContractStateAll(contractStateAllMsg types.ContractStateAllMsg) *XplaClient {
	msg, err := mwasm.MakeContractStateAllMsg(contractStateAllMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmContractStateAllMsgType
	xplac.Msg = msg
	return xplac
}

// Prints out the code history for a contract given its address.
func (xplac *XplaClient) ContractHistory(contractHistoryMsg types.ContractHistoryMsg) *XplaClient {
	msg, err := mwasm.MakeContractHistoryMsg(contractHistoryMsg)
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmContractHistoryMsgType
	xplac.Msg = msg
	return xplac
}

// Query list all pinned code IDs.
func (xplac *XplaClient) Pinned() *XplaClient {
	msg, err := mwasm.MakePinnedMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmPinnedMsgType
	xplac.Msg = msg
	return xplac
}

// Get libwasmvm version.
func (xplac *XplaClient) LibwasmvmVersion() *XplaClient {
	msg, err := mwasm.MakeLibwasmvmVersionMsg()
	if err != nil {
		xplac.Err = err
		return xplac
	}
	xplac.Module = mwasm.WasmModule
	xplac.MsgType = mwasm.WasmLibwasmvmVersionMsgType
	xplac.Msg = msg
	return xplac
}
