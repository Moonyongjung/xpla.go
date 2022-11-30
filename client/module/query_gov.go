package module

import (
	mgov "github.com/Moonyongjung/xpla.go/core/gov"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	govv1beta1 "cosmossdk.io/api/cosmos/gov/v1beta1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govutils "github.com/cosmos/cosmos-sdk/x/gov/client/utils"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// Query client for gov module.
func (i IXplaClient) QueryGov() (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcGov(i)
	} else {
		return queryByLcdGov(i)
	}
}

func queryByGrpcGov(i IXplaClient) (string, error) {
	queryClient := govtypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Gov proposal
	case i.Ixplac.GetMsgType() == mgov.GovQueryProposalMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryProposalRequest)
		res, err = queryClient.Proposal(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Gov proposals
	case i.Ixplac.GetMsgType() == mgov.GovQueryProposalsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryProposalsRequest)
		res, err = queryClient.Proposals(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Gov deposit parameter
	case i.Ixplac.GetMsgType() == mgov.GovQueryDepositParamsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryDepositParams)

		var deposit govtypes.Deposit

		clientCtx, err := clientForQuery(i)
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
	case i.Ixplac.GetMsgType() == mgov.GovQueryDepositRequestMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryDepositRequest)
		res, err = queryClient.Deposit(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Gov deposits parameter
	case i.Ixplac.GetMsgType() == mgov.GovQueryDepositsParamsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryProposalParams)

		var deposit govtypes.Deposit
		clientCtx, err := clientForQuery(i)
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
	case i.Ixplac.GetMsgType() == mgov.GovQueryDepositsRequestMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryDepositsRequest)
		res, err = queryClient.Deposits(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Gov tally
	case i.Ixplac.GetMsgType() == mgov.GovTallyMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryTallyResultRequest)
		res, err = queryClient.TallyResult(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Gov params
	case i.Ixplac.GetMsgType() == mgov.GovQueryGovParamsMsgType:
		votingRes, err := queryClient.Params(
			i.Ixplac.GetContext(),
			&govtypes.QueryParamsRequest{ParamsType: "voting"},
		)
		if err != nil {
			return "", err
		}

		tallyRes, err := queryClient.Params(
			i.Ixplac.GetContext(),
			&govtypes.QueryParamsRequest{ParamsType: "tallying"},
		)
		if err != nil {
			return "", err
		}

		depositRes, err := queryClient.Params(
			i.Ixplac.GetContext(),
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
	case i.Ixplac.GetMsgType() == mgov.GovQueryGovParamVotingMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryParamsRequest)
		resParams, err := queryClient.Params(
			i.Ixplac.GetContext(),
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
	case i.Ixplac.GetMsgType() == mgov.GovQueryGovParamTallyingMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryParamsRequest)
		resParams, err := queryClient.Params(
			i.Ixplac.GetContext(),
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
	case i.Ixplac.GetMsgType() == mgov.GovQueryGovParamDepositMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryParamsRequest)
		resParams, err := queryClient.Params(
			i.Ixplac.GetContext(),
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
	case i.Ixplac.GetMsgType() == mgov.GovQueryProposerMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(string)
		proposalId := util.FromStringToUint64(convertMsg)

		clientCtx, err := clientForQuery(i)
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
	case i.Ixplac.GetMsgType() == mgov.GovQueryVoteMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryVoteRequest)
		resVote, err := queryClient.Vote(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

		clientCtx, err := clientForQuery(i)
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
	case i.Ixplac.GetMsgType() == mgov.GovQueryVotesNotPassedMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryProposalVotesParams)
		clientCtx, err := clientForQuery(i)
		if err != nil {
			return "", err
		}
		resByTxQuery, err := govutils.QueryVotesByTxQuery(clientCtx, convertMsg)
		if err != nil {
			return "", err
		}

		var votes govtypes.Votes

		clientCtx.LegacyAmino.MustUnmarshalJSON(resByTxQuery, &votes)
		out, err := printObjectLegacy(i, votes)
		if err != nil {
			return "", err
		}
		return string(out), nil

	// Gov votes passed
	case i.Ixplac.GetMsgType() == mgov.GovQueryVotesPassedMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryVotesRequest)
		res, err = queryClient.Votes(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err = printProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

const (
	govProposalsLabel = "proposals"
	govDeposistsLabel = "deposits"
	govTallyLabel     = "tally"
	govParamsLabel    = "params"
	govVotesLabel     = "votes"
)

func queryByLcdGov(i IXplaClient) (string, error) {
	url := util.MakeQueryLcdUrl(govv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Gov proposal
	case i.Ixplac.GetMsgType() == mgov.GovQueryProposalMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryProposalRequest)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalId))

		util.LogInfo(convertMsg)

	// Gov proposals
	case i.Ixplac.GetMsgType() == mgov.GovQueryProposalsMsgType:
		url = url + govProposalsLabel

	// Gov deposit parameter
	case i.Ixplac.GetMsgType() == mgov.GovQueryDepositParamsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryDepositParams)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalID), govDeposistsLabel, convertMsg.Depositor.String())

	// Gov deposit
	case i.Ixplac.GetMsgType() == mgov.GovQueryDepositRequestMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryDepositRequest)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalId), govDeposistsLabel, convertMsg.Depositor)

	// Gov deposits parameter
	case i.Ixplac.GetMsgType() == mgov.GovQueryDepositsParamsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryProposalParams)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalID), govDeposistsLabel)

	// Gov deposits
	case i.Ixplac.GetMsgType() == mgov.GovQueryDepositsRequestMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryDepositsRequest)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalId), govDeposistsLabel)

	// Gov tally
	case i.Ixplac.GetMsgType() == mgov.GovTallyMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryTallyResultRequest)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalId), govTallyLabel)

	// Gov params
	case i.Ixplac.GetMsgType() == mgov.GovQueryGovParamsMsgType:
		return "", util.LogErr("unsupported querying all gov params by using LCD. query each parameter(voting|tallying|deposit)")

	// Gov params of voting
	case i.Ixplac.GetMsgType() == mgov.GovQueryGovParamVotingMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryParamsRequest)

		url = url + util.MakeQueryLabels(govParamsLabel, convertMsg.ParamsType)

	// Gov params of tally
	case i.Ixplac.GetMsgType() == mgov.GovQueryGovParamTallyingMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryParamsRequest)

		url = url + util.MakeQueryLabels(govParamsLabel, convertMsg.ParamsType)

	// Gov params of deposit
	case i.Ixplac.GetMsgType() == mgov.GovQueryGovParamDepositMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryParamsRequest)

		url = url + util.MakeQueryLabels(govParamsLabel, convertMsg.ParamsType)

	// Gov proposer
	case i.Ixplac.GetMsgType() == mgov.GovQueryProposerMsgType:
		return "", util.LogErr("unsupported querying proposer by using LCD")

	// Gov vote
	case i.Ixplac.GetMsgType() == mgov.GovQueryVoteMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryVoteRequest)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalId), govVotesLabel, convertMsg.Voter)

	// Gov votes not passed
	case i.Ixplac.GetMsgType() == mgov.GovQueryVotesNotPassedMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryProposalVotesParams)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalID), govVotesLabel)

	// Gov votes passed
	case i.Ixplac.GetMsgType() == mgov.GovQueryVotesPassedMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(govtypes.QueryVotesRequest)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalId), govVotesLabel)

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil

}
