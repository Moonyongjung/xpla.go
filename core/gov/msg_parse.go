package gov

import (
	"context"

	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govutils "github.com/cosmos/cosmos-sdk/x/gov/client/utils"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/grpc"
)

// Parsing - submit proposal
func parseSubmitProposalArgs(submitProposalMsg types.SubmitProposalMsg, privKey key.PrivateKey) (*govtypes.MsgSubmitProposal, error) {
	proposer := util.GetAddrByPrivKey(privKey)
	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(submitProposalMsg.Deposit))
	if err != nil {
		return nil, err
	}

	content := govtypes.ContentFromProposalType(
		submitProposalMsg.Title,
		submitProposalMsg.Description,
		govutils.NormalizeProposalType(submitProposalMsg.Type),
	)

	msg, err := govtypes.NewMsgSubmitProposal(content, amount, proposer)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Parsing - deposit
func parseGovDepositArgs(govDepositMsg types.GovDepositMsg, privKey key.PrivateKey) (*govtypes.MsgDeposit, error) {
	proposalId := util.FromStringToUint64(govDepositMsg.ProposalID)
	from := util.GetAddrByPrivKey(privKey)
	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(govDepositMsg.Deposit))
	if err != nil {
		return nil, err
	}

	msg := govtypes.NewMsgDeposit(from, proposalId, amount)

	return msg, nil
}

// Parsing - vote
func parseVoteArgs(voteMsg types.VoteMsg, privKey key.PrivateKey) (*govtypes.MsgVote, error) {
	proposalId := util.FromStringToUint64(voteMsg.ProposalID)
	from := util.GetAddrByPrivKey(privKey)
	byteVoteOption, err := govtypes.VoteOptionFromString(govutils.NormalizeVoteOption(voteMsg.Option))
	if err != nil {
		return nil, err
	}

	msg := govtypes.NewMsgVote(from, proposalId, byteVoteOption)
	return msg, nil
}

// Parsing - weighted vote
func parseWeightedVoteArgs(weightedVoteMsg types.WeightedVoteMsg, privKey key.PrivateKey) (*govtypes.MsgVoteWeighted, error) {
	proposalId := util.FromStringToUint64(weightedVoteMsg.ProposalID)
	from := util.GetAddrByPrivKey(privKey)

	options, err := weightedVoteOptionConverting(weightedVoteMsg)
	if err != nil {
		return nil, err
	}

	msg := govtypes.NewMsgVoteWeighted(from, proposalId, options)
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Parsing - proposal
func parseQueryProposalArgs(queryProposalMsg types.QueryProposalMsg) (govtypes.QueryProposalRequest, error) {
	return govtypes.QueryProposalRequest{
		ProposalId: util.FromStringToUint64(queryProposalMsg.ProposalID),
	}, nil
}

// Parsing - proposals
func parseQueryProposalsArgs(queryProposalsMsg types.QueryProposalsMsg) (govtypes.QueryProposalsRequest, error) {
	depositorAddr := queryProposalsMsg.Depositor
	voterAddr := queryProposalsMsg.Voter
	strProposalStatus := queryProposalsMsg.Status

	var proposalStatus govtypes.ProposalStatus

	if len(depositorAddr) != 0 {
		_, err := sdk.AccAddressFromBech32(depositorAddr)
		if err != nil {
			return govtypes.QueryProposalsRequest{}, err
		}
	}

	if len(voterAddr) != 0 {
		_, err := sdk.AccAddressFromBech32(voterAddr)
		if err != nil {
			return govtypes.QueryProposalsRequest{}, err
		}
	}

	if len(strProposalStatus) != 0 {
		proposalStatus1, err := govtypes.ProposalStatusFromString(govutils.NormalizeProposalStatus(strProposalStatus))
		proposalStatus = proposalStatus1
		if err != nil {
			return govtypes.QueryProposalsRequest{}, err
		}
	}

	return govtypes.QueryProposalsRequest{
		ProposalStatus: proposalStatus,
		Voter:          voterAddr,
		Depositor:      depositorAddr,
		Pagination:     core.PageRequest,
	}, nil
}

// Parsing - query deposit
func parseQueryDepositArgs(queryDepositMsg types.QueryDepositMsg, grpcConn grpc.ClientConn, ctx context.Context) (interface{}, string, error) {
	queryClient := govtypes.NewQueryClient(grpcConn)
	proposalId := util.FromStringToUint64(queryDepositMsg.ProposalID)

	proposalRes, err := queryClient.Proposal(
		ctx,
		&govtypes.QueryProposalRequest{ProposalId: proposalId},
	)
	if err != nil {
		return nil, "", err
	}

	depositorAddr, err := sdk.AccAddressFromBech32(queryDepositMsg.Depositor)
	if err != nil {
		return nil, "", err
	}

	propStatus := proposalRes.Proposal.Status
	if !(propStatus == govtypes.StatusVotingPeriod || propStatus == govtypes.StatusDepositPeriod) {
		params := govtypes.NewQueryDepositParams(proposalId, depositorAddr)
		return params, "params", nil
	}

	return govtypes.QueryDepositRequest{
		ProposalId: proposalId,
		Depositor:  queryDepositMsg.Depositor,
	}, "request", nil
}

// Parsing - query deposits
func parseQueryDepositsArgs(queryDepositMsg types.QueryDepositMsg, grpcConn grpc.ClientConn, ctx context.Context) (interface{}, string, error) {
	queryClient := govtypes.NewQueryClient(grpcConn)
	proposalId := util.FromStringToUint64(queryDepositMsg.ProposalID)

	proposalRes, err := queryClient.Proposal(
		ctx,
		&govtypes.QueryProposalRequest{ProposalId: proposalId},
	)
	if err != nil {
		return nil, "", err
	}

	propStatus := proposalRes.GetProposal().Status
	if !(propStatus == govtypes.StatusVotingPeriod || propStatus == govtypes.StatusDepositPeriod) {
		params := govtypes.NewQueryProposalParams(proposalId)
		return params, "params", nil
	}

	return govtypes.QueryDepositsRequest{
		ProposalId: proposalId,
		Pagination: core.PageRequest,
	}, "request", nil
}

// Parsing - tally
func parseGovTallyArgs(tallyMsg types.TallyMsg, grpcConn grpc.ClientConn, ctx context.Context) (govtypes.QueryTallyResultRequest, error) {
	queryClient := govtypes.NewQueryClient(grpcConn)
	proposalId := util.FromStringToUint64(tallyMsg.ProposalID)

	_, err := queryClient.Proposal(
		ctx,
		&govtypes.QueryProposalRequest{ProposalId: proposalId},
	)
	if err != nil {
		return govtypes.QueryTallyResultRequest{}, util.LogErr("failed to fetch proposal-id", proposalId, " : ", err)
	}

	return govtypes.QueryTallyResultRequest{
		ProposalId: proposalId,
	}, nil
}

// Parsing - gov params
func parseGovParamArgs(govParamsMsg types.GovParamsMsg) (govtypes.QueryParamsRequest, error) {
	if govParamsMsg.ParamType == "voting" ||
		govParamsMsg.ParamType == "tallying" ||
		govParamsMsg.ParamType == "deposit" {
		return govtypes.QueryParamsRequest{
			ParamsType: govParamsMsg.ParamType,
		}, nil
	} else {
		return govtypes.QueryParamsRequest{}, util.LogErr("argument must be one of (voting|tallying|deposit), was ", govParamsMsg.ParamType)
	}
}

// Parsing - query vote
func parseQueryVoteArgs(queryVoteMsg types.QueryVoteMsg, grpcConn grpc.ClientConn, ctx context.Context) (govtypes.QueryVoteRequest, error) {
	queryClient := govtypes.NewQueryClient(grpcConn)
	proposalId := util.FromStringToUint64(queryVoteMsg.ProposalID)

	_, err := queryClient.Proposal(
		ctx,
		&govtypes.QueryProposalRequest{ProposalId: proposalId},
	)
	if err != nil {
		return govtypes.QueryVoteRequest{}, err
	}

	return govtypes.QueryVoteRequest{
		ProposalId: proposalId,
		Voter:      queryVoteMsg.VoterAddr,
	}, nil
}

// Parsing - query votes
func parseQueryVotesArgs(queryVoteMsg types.QueryVoteMsg, grpcConn grpc.ClientConn, ctx context.Context) (interface{}, string, error) {
	queryClient := govtypes.NewQueryClient(grpcConn)
	proposalId := util.FromStringToUint64(queryVoteMsg.ProposalID)

	res, err := queryClient.Proposal(
		ctx,
		&govtypes.QueryProposalRequest{ProposalId: proposalId},
	)
	if err != nil {
		return govtypes.QueryVoteRequest{}, "", err
	}

	status := res.GetProposal().Status
	if !(status == govtypes.StatusVotingPeriod || status == govtypes.StatusDepositPeriod) {
		params := govtypes.NewQueryProposalVotesParams(proposalId, 0, 0)
		return params, "notPassed", nil
	}

	return govtypes.QueryVotesRequest{
		ProposalId: proposalId,
		Pagination: core.PageRequest,
	}, "passed", nil

}

func weightedVoteOptionConverting(weightedVoteMsg types.WeightedVoteMsg) (govtypes.WeightedVoteOptions, error) {
	weightedVoteOptions := govtypes.WeightedVoteOptions{}

	if weightedVoteMsg.Yes != "" {
		weightedVoteOptions = append(weightedVoteOptions, govtypes.WeightedVoteOption{
			Option: govtypes.OptionYes,
			Weight: sdk.MustNewDecFromStr(weightedVoteMsg.Yes),
		})
	}
	if weightedVoteMsg.Abstain != "" {
		weightedVoteOptions = append(weightedVoteOptions, govtypes.WeightedVoteOption{
			Option: govtypes.OptionAbstain,
			Weight: sdk.MustNewDecFromStr(weightedVoteMsg.Abstain),
		})
	}
	if weightedVoteMsg.No != "" {
		weightedVoteOptions = append(weightedVoteOptions, govtypes.WeightedVoteOption{
			Option: govtypes.OptionNo,
			Weight: sdk.MustNewDecFromStr(weightedVoteMsg.No),
		})
	}
	if weightedVoteMsg.NoWithVeto != "" {
		weightedVoteOptions = append(weightedVoteOptions, govtypes.WeightedVoteOption{
			Option: govtypes.OptionNoWithVeto,
			Weight: sdk.MustNewDecFromStr(weightedVoteMsg.NoWithVeto),
		})
	}

	return weightedVoteOptions, nil
}
