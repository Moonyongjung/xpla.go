package client

import (
	"context"
	mgov "github.com/Moonyongjung/xpla.go/core/gov"
	"github.com/Moonyongjung/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govutils "github.com/cosmos/cosmos-sdk/x/gov/client/utils"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// Query client for gov module.
func queryGov(xplac *XplaClient) (string, error) {
	queryClient := govtypes.NewQueryClient(xplac.Grpc)

	switch {
	// Gov proposal
	case xplac.MsgType == mgov.GovQueryProposalMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryProposalRequest)
		res, err = queryClient.Proposal(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Gov proposals
	case xplac.MsgType == mgov.GovQueryProposalsMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryProposalsRequest)
		res, err = queryClient.Proposals(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Gov deposit parameter
	case xplac.MsgType == mgov.GovQueryDepositParamsMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryDepositParams)

		var deposit govtypes.Deposit

		clientCtx, err := clientForQuery(xplac)
		if err != nil {
			return "", err
		}

		resByTxQuery, err := govutils.QueryDepositByTxQuery(clientCtx, convertMsg)
		if err != nil {
			return "", err
		}
		clientCtx.Codec.MustUnmarshalJSON(resByTxQuery, &deposit)
		res = &deposit

	// Gov deposit
	case xplac.MsgType == mgov.GovQueryDepositRequestMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryDepositRequest)
		res, err = queryClient.Deposit(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Gov deposits parameter
	case xplac.MsgType == mgov.GovQueryDepositsParamsMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryProposalParams)

		var deposit govtypes.Deposit
		clientCtx, err := clientForQuery(xplac)
		if err != nil {
			return "", err
		}

		resByTxQuery, err := govutils.QueryDepositsByTxQuery(clientCtx, convertMsg)
		if err != nil {
			return "", err
		}
		clientCtx.Codec.MustUnmarshalJSON(resByTxQuery, &deposit)
		res = &deposit

	// Gov deposits
	case xplac.MsgType == mgov.GovQueryDepositsRequestMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryDepositsRequest)
		res, err = queryClient.Deposits(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Gov tally
	case xplac.MsgType == mgov.GovTallyMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryTallyResultRequest)
		res, err = queryClient.TallyResult(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Gov params
	case xplac.MsgType == mgov.GovQueryGovParamsMsgType:
		votingRes, err := queryClient.Params(
			context.Background(),
			&govtypes.QueryParamsRequest{ParamsType: "voting"},
		)
		if err != nil {
			return "", err
		}

		tallyRes, err := queryClient.Params(
			context.Background(),
			&govtypes.QueryParamsRequest{ParamsType: "tallying"},
		)
		if err != nil {
			return "", err
		}

		depositRes, err := queryClient.Params(
			context.Background(),
			&govtypes.QueryParamsRequest{ParamsType: "deposit"},
		)
		if err != nil {
			return "", err
		}

		govAllParams := govtypes.NewParams(
			votingRes.GetVotingParams(),
			tallyRes.GetTallyParams(),
			depositRes.GetDepositParams(),
		)

		bytes, err := util.JsonMarshalData(govAllParams)
		if err != nil {
			return "", err
		}
		return string(bytes), nil

	// Gov params of voting
	case xplac.MsgType == mgov.GovQueryGovParamVotingMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryParamsRequest)
		resParams, err := queryClient.Params(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

		bytes, err := util.JsonMarshalData(resParams.GetVotingParams())
		if err != nil {
			return "", err
		}
		return string(bytes), nil

	// Gov params of tally
	case xplac.MsgType == mgov.GovQueryGovParamTallyingMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryParamsRequest)
		resParams, err := queryClient.Params(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

		bytes, err := util.JsonMarshalData(resParams.GetTallyParams())
		if err != nil {
			return "", err
		}
		return string(bytes), nil

	// Gov params of deposit
	case xplac.MsgType == mgov.GovQueryGovParamDepositMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryParamsRequest)
		resParams, err := queryClient.Params(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

		bytes, err := util.JsonMarshalData(resParams.GetDepositParams())
		if err != nil {
			return "", err
		}
		return string(bytes), nil

	// Gov proposer
	case xplac.MsgType == mgov.GovQueryProposerMsgType:
		convertMsg, _ := xplac.Msg.(string)
		proposalId := util.FromStringToUint64(convertMsg)

		clientCtx, err := clientForQuery(xplac)
		if err != nil {
			return "", err
		}

		prop, err := govutils.QueryProposerByTxQuery(clientCtx, proposalId)
		if err != nil {
			return "", err
		}

		bytes, err := util.JsonMarshalData(prop)
		if err != nil {
			return "", err
		}
		return string(bytes), nil

	// Gov vote
	case xplac.MsgType == mgov.GovQueryVoteMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryVoteRequest)
		resVote, err := queryClient.Vote(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

		clientCtx, err := clientForQuery(xplac)
		if err != nil {
			return "", err
		}

		voterAddr, err := sdk.AccAddressFromBech32(convertMsg.Voter)
		if err != nil {
			return "", err
		}

		vote := resVote.GetVote()
		if vote.Empty() {
			params := govtypes.NewQueryVoteParams(convertMsg.ProposalId, voterAddr)
			resByTxQuery, err := govutils.QueryVoteByTxQuery(clientCtx, params)
			if err != nil {
				return "", err
			}

			if err := clientCtx.Codec.UnmarshalJSON(resByTxQuery, &vote); err != nil {
				return "", err
			}
		}

		res = &resVote.Vote

	// Gov votes not passed
	case xplac.MsgType == mgov.GovQueryVotesNotPassedMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryProposalVotesParams)
		clientCtx, err := clientForQuery(xplac)
		if err != nil {
			return "", err
		}
		resByTxQuery, err := govutils.QueryVotesByTxQuery(clientCtx, convertMsg)
		if err != nil {
			return "", err
		}

		var votes govtypes.Votes

		clientCtx.LegacyAmino.MustUnmarshalJSON(resByTxQuery, &votes)
		out, err := printObjectLegacy(xplac, votes)
		if err != nil {
			return "", err
		}
		return string(out), nil

	// Gov votes passed
	case xplac.MsgType == mgov.GovQueryVotesPassedMsgType:
		convertMsg, _ := xplac.Msg.(govtypes.QueryVotesRequest)
		res, err = queryClient.Votes(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err = printProto(xplac, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
