package gov

import (
	"context"

	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/grpc"
)

// (Tx) make msg - submit proposal
func MakeSubmitProposalMsg(submitProposalMsg types.SubmitProposalMsg, privKey key.PrivateKey) (govtypes.MsgSubmitProposal, error) {
	msg, err := parseSubmitProposalArgs(submitProposalMsg, privKey)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, err
	}

	return msg, nil
}

// (Tx) make msg - deposit
func MakeGovDepositMsg(govDepositMsg types.GovDepositMsg, privKey key.PrivateKey) (govtypes.MsgDeposit, error) {
	msg, err := parseGovDepositArgs(govDepositMsg, privKey)
	if err != nil {
		return govtypes.MsgDeposit{}, err
	}

	return msg, nil
}

// (Tx) make msg - vote
func MakeVoteMsg(voteMsg types.VoteMsg, privKey key.PrivateKey) (govtypes.MsgVote, error) {
	msg, err := parseVoteArgs(voteMsg, privKey)
	if err != nil {
		return govtypes.MsgVote{}, err
	}

	return msg, nil
}

// (Tx) make msg - weighted vote
func MakeWeightedVoteMsg(weightedVoteMsg types.WeightedVoteMsg, privKey key.PrivateKey) (govtypes.MsgVoteWeighted, error) {
	msg, err := parseWeightedVoteArgs(weightedVoteMsg, privKey)
	if err != nil {
		return govtypes.MsgVoteWeighted{}, err
	}

	return msg, nil
}

// (Query) make msg - proposal
func MakeQueryProposalMsg(queryProposalMsg types.QueryProposalMsg) (govtypes.QueryProposalRequest, error) {
	msg, err := parseQueryProposalArgs(queryProposalMsg)
	if err != nil {
		return govtypes.QueryProposalRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - proposals
func MakeQueryProposalsMsg(queryProposalsMsg types.QueryProposalsMsg) (govtypes.QueryProposalsRequest, error) {
	msg, err := parseQueryProposalsArgs(queryProposalsMsg)
	if err != nil {
		return govtypes.QueryProposalsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - query deposit
func MakeQueryDepositMsg(queryDepositMsg types.QueryDepositMsg, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, string, error) {
	msg, argsType, err := parseQueryDepositArgs(queryDepositMsg, grpcConn, ctx, lcdUrl, queryType)
	if err != nil {
		return nil, "", err
	}

	return msg, argsType, nil
}

// (Query) make msg - query deposits
func MakeQueryDepositsMsg(queryDepositMsg types.QueryDepositMsg, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, string, error) {
	msg, argsType, err := parseQueryDepositsArgs(queryDepositMsg, grpcConn, ctx, lcdUrl, queryType)
	if err != nil {
		return nil, "", err
	}

	return msg, argsType, nil
}

// (Query) make msg - tally
func MakeGovTallyMsg(tallyMsg types.TallyMsg, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, error) {
	msg, err := parseGovTallyArgs(tallyMsg, grpcConn, ctx, lcdUrl, queryType)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// (Query) make msg - gov params
func MakeGovParamsMsg(govParamsMsg types.GovParamsMsg) (govtypes.QueryParamsRequest, error) {
	msg, err := parseGovParamArgs(govParamsMsg)
	if err != nil {
		return govtypes.QueryParamsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - query vote
func MakeQueryVoteMsg(queryVoteMsg types.QueryVoteMsg, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (govtypes.QueryVoteRequest, error) {
	msg, err := parseQueryVoteArgs(queryVoteMsg, grpcConn, ctx, lcdUrl, queryType)
	if err != nil {
		return govtypes.QueryVoteRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - query votes
func MakeQueryVotesMsg(queryVoteMsg types.QueryVoteMsg, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, string, error) {
	msg, status, err := parseQueryVotesArgs(queryVoteMsg, grpcConn, ctx, lcdUrl, queryType)
	if err != nil {
		return nil, "", err
	}

	return msg, status, nil
}
