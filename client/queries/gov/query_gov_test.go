package gov_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Moonyongjung/xpla.go/client"
	"github.com/Moonyongjung/xpla.go/client/queries/qtest"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util/testutil"
	"github.com/gogo/protobuf/jsonpb"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtestutil "github.com/cosmos/cosmos-sdk/x/gov/client/testutil"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"
)

var validatorNumber = 1

type IntegrationTestSuite struct {
	suite.Suite

	xplac *client.XplaClient
	apis  []string

	cfg     network.Config
	network *network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = network.New(s.T(), s.cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	val := s.network.Validators[0]

	// create a proposal with deposit
	_, err = govtestutil.MsgSubmitProposal(val.ClientCtx, val.Address.String(),
		"Text Proposal 1", "Where is the title!?", govtypes.ProposalTypeText,
		fmt.Sprintf("--%s=%s", govcli.FlagDeposit, sdk.NewCoin(s.cfg.BondDenom, govtypes.DefaultMinDepositTokens).String()))
	s.Require().NoError(err)
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	// vote for proposal
	_, err = govtestutil.MsgVote(val.ClientCtx, val.Address.String(), "1", "yes")
	s.Require().NoError(err)

	// create a proposal without deposit
	_, err = govtestutil.MsgSubmitProposal(val.ClientCtx, val.Address.String(),
		"Text Proposal 2", "Where is the title!?", govtypes.ProposalTypeText)
	s.Require().NoError(err)
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	// create a proposal3 with deposit
	_, err = govtestutil.MsgSubmitProposal(val.ClientCtx, val.Address.String(),
		"Text Proposal 3", "Where is the title!?", govtypes.ProposalTypeText,
		fmt.Sprintf("--%s=%s", govcli.FlagDeposit, sdk.NewCoin(s.cfg.BondDenom, govtypes.DefaultMinDepositTokens).String()))
	s.Require().NoError(err)
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	// vote for proposal3 as val
	_, err = govtestutil.MsgVote(val.ClientCtx, val.Address.String(), "3", "yes=0.6,no=0.3,abstain=0.05,no_with_veto=0.05")
	s.Require().NoError(err)

	s.xplac = qtest.NewTestXplaClient()
	s.apis = []string{
		s.network.Validators[0].APIAddress,
		s.network.Validators[0].AppConfig.GRPC.Address,
	}

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestQueryProposal() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		queryProposalMsg := types.QueryProposalMsg{
			ProposalID: "1",
		}
		res, err := s.xplac.QueryProposal(queryProposalMsg).Query()
		s.Require().NoError(err)

		var queryProposalResponse govtypes.QueryProposalResponse

		jsonpb.Unmarshal(strings.NewReader(res), &queryProposalResponse)

		var content govtypes.Content
		s.xplac.EncodingConfig.InterfaceRegistry.UnpackAny(queryProposalResponse.Proposal.Content, &content)

		s.Require().Equal("/cosmos.gov.v1beta1.TextProposal", queryProposalResponse.Proposal.Content.TypeUrl)
		s.Require().Equal("Text Proposal 1", content.GetTitle())
		s.Require().Equal("Where is the title!?", content.GetDescription())
	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestQueryProposals() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		queryProposalsMsg := types.QueryProposalsMsg{}

		res, err := s.xplac.QueryProposals(queryProposalsMsg).Query()
		s.Require().NoError(err)

		var queryProposalsResponse govtypes.QueryProposalsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryProposalsResponse)

		s.Require().Equal(3, len(queryProposalsResponse.Proposals))
	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestDeposit() {
	val := s.network.Validators[0].Address.String()

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		queryDepositMsg := types.QueryDepositMsg{
			ProposalID: "1",
			Depositor:  val,
		}

		res, err := s.xplac.QueryDeposit(queryDepositMsg).Query()
		s.Require().NoError(err)

		var queryDepositResponse govtypes.QueryDepositResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryDepositResponse)

		s.Require().Equal(val, queryDepositResponse.Deposit.Depositor)
		s.Require().Equal(val, queryDepositResponse.Deposit.Depositor)
		s.Require().Equal(govtypes.DefaultMinDepositTokens, queryDepositResponse.Deposit.Amount[0].Amount)
		s.Require().Equal("stake", queryDepositResponse.Deposit.Amount[0].Denom)
	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestVote() {
	val := s.network.Validators[0].Address.String()

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		queryVoteMsg := types.QueryVoteMsg{
			ProposalID: "1",
			VoterAddr:  val,
		}
		res1, err := s.xplac.QueryVote(queryVoteMsg).Query()
		s.Require().NoError(err)

		if i == 0 {
			var queryVoteResponse govtypes.QueryVoteResponse
			jsonpb.Unmarshal(strings.NewReader(res1), &queryVoteResponse)

			s.Require().Equal(val, queryVoteResponse.Vote.Voter)
			s.Require().Equal(uint64(1), queryVoteResponse.Vote.ProposalId)
		} else {
			var vote govtypes.Vote
			jsonpb.Unmarshal(strings.NewReader(res1), &vote)

			s.Require().Equal(val, vote.Voter)
			s.Require().Equal(uint64(1), vote.ProposalId)
		}
	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestTally() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		tallyMsg := types.TallyMsg{
			ProposalID: "3",
		}

		res, err := s.xplac.Tally(tallyMsg).Query()
		s.Require().NoError(err)

		var queryTallyResultResponse govtypes.QueryTallyResultResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryTallyResultResponse)

		s.T().Log(queryTallyResultResponse.Tally)

		s.Require().Equal("60000000000000000000", queryTallyResultResponse.Tally.Yes.String())
		s.Require().Equal("5000000000000000000", queryTallyResultResponse.Tally.Abstain.String())
		s.Require().Equal("30000000000000000000", queryTallyResultResponse.Tally.No.String())
		s.Require().Equal("5000000000000000000", queryTallyResultResponse.Tally.NoWithVeto.String())

	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestGovParams() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)

			// only query tally params
			govParamsMsg := types.GovParamsMsg{
				ParamType: "tallying",
			}

			res1, err := s.xplac.GovParams(govParamsMsg).Query()
			s.Require().NoError(err)

			var queryParamsResponse1 govtypes.QueryParamsResponse
			jsonpb.Unmarshal(strings.NewReader(res1), &queryParamsResponse1)

			// can check tally
			s.Require().Equal("0.334000000000000000", queryParamsResponse1.TallyParams.Quorum.String())
			s.Require().Equal("0.500000000000000000", queryParamsResponse1.TallyParams.Threshold.String())
			s.Require().Equal("0.334000000000000000", queryParamsResponse1.TallyParams.VetoThreshold.String())
			s.Require().Equal("0s", queryParamsResponse1.VotingParams.VotingPeriod.String())
			s.Require().Equal(0, len(queryParamsResponse1.DepositParams.MinDeposit))
			s.Require().Equal("0s", queryParamsResponse1.DepositParams.MaxDepositPeriod.String())

			// only query voting params
			govParamsMsg = types.GovParamsMsg{
				ParamType: "voting",
			}

			res2, err := s.xplac.GovParams(govParamsMsg).Query()
			s.Require().NoError(err)

			var queryParamsResponse2 govtypes.QueryParamsResponse
			jsonpb.Unmarshal(strings.NewReader(res2), &queryParamsResponse2)

			// can check voting
			s.Require().Equal("0.000000000000000000", queryParamsResponse2.TallyParams.Quorum.String())
			s.Require().Equal("0.000000000000000000", queryParamsResponse2.TallyParams.Threshold.String())
			s.Require().Equal("0.000000000000000000", queryParamsResponse2.TallyParams.VetoThreshold.String())
			s.Require().Equal("48h0m0s", queryParamsResponse2.VotingParams.VotingPeriod.String())
			s.Require().Equal(0, len(queryParamsResponse2.DepositParams.MinDeposit))
			s.Require().Equal("0s", queryParamsResponse2.DepositParams.MaxDepositPeriod.String())

			// only query deposit params
			govParamsMsg = types.GovParamsMsg{
				ParamType: "deposit",
			}

			res3, err := s.xplac.GovParams(govParamsMsg).Query()
			s.Require().NoError(err)

			var queryParamsResponse3 govtypes.QueryParamsResponse
			jsonpb.Unmarshal(strings.NewReader(res3), &queryParamsResponse3)

			// can check deposit
			s.Require().Equal("0.000000000000000000", queryParamsResponse3.TallyParams.Quorum.String())
			s.Require().Equal("0.000000000000000000", queryParamsResponse3.TallyParams.Threshold.String())
			s.Require().Equal("0.000000000000000000", queryParamsResponse3.TallyParams.VetoThreshold.String())
			s.Require().Equal("0s", queryParamsResponse3.VotingParams.VotingPeriod.String())
			s.Require().Equal("stake", queryParamsResponse3.DepositParams.MinDeposit[0].Denom)
			s.Require().Equal("10000000", queryParamsResponse3.DepositParams.MinDeposit[0].Amount.String())
			s.Require().Equal("48h0m0s", queryParamsResponse3.DepositParams.MaxDepositPeriod.String())

		} else {
			s.xplac.WithGrpc(api)

			// only query tally params
			govParamsMsg := types.GovParamsMsg{
				ParamType: "tallying",
			}

			res1, err := s.xplac.GovParams(govParamsMsg).Query()
			s.Require().NoError(err)

			var tallyParams govtypes.TallyParams
			jsonpb.Unmarshal(strings.NewReader(res1), &tallyParams)

			s.Require().Equal("0.334000000000000000", tallyParams.Quorum.String())
			s.Require().Equal("0.500000000000000000", tallyParams.Threshold.String())
			s.Require().Equal("0.334000000000000000", tallyParams.VetoThreshold.String())

			// only query voting params
			govParamsMsg = types.GovParamsMsg{
				ParamType: "voting",
			}

			res2, err := s.xplac.GovParams(govParamsMsg).Query()
			s.Require().NoError(err)
			s.Require().True(strings.Contains(res2, "172800000000000"))

			// only query deposit params
			govParamsMsg = types.GovParamsMsg{
				ParamType: "deposit",
			}

			res3, err := s.xplac.GovParams(govParamsMsg).Query()
			s.Require().NoError(err)

			var depositParams govtypes.DepositParams
			jsonpb.Unmarshal(strings.NewReader(res3), &depositParams)

			s.Require().Equal("10000000", depositParams.MinDeposit[0].Amount.String())
			s.Require().True(strings.Contains(res3, "172800000000000"))

			// query all gov params (not support LCD)
			res4, err := s.xplac.GovParams().Query()
			s.Require().NoError(err)

			var queryParamsResponse4 govtypes.QueryParamsResponse
			jsonpb.Unmarshal(strings.NewReader(res4), &queryParamsResponse4)

			expectedResult := `{"voting_params":{"voting_period":172800000000000},"tally_params":{"quorum":"0.334000000000000000","threshold":"0.500000000000000000","veto_threshold":"0.334000000000000000"},"deposit_params":{"min_deposit":[{"denom":"stake","amount":"10000000"}],"max_deposit_period":172800000000000}}`
			s.Require().Equal(res4, expectedResult)
		}

	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

// func (s *IntegrationTestSuite) TestProposer() {
// 	// not supported LCD
// 	s.xplac.WithGrpc(s.apis[1])

// 	s.T().Log(s.apis[1])

// 	ProposerMsg := types.ProposerMsg{
// 		ProposalID: "2",
// 	}

// 	res, err := s.xplac.Proposer(ProposerMsg).Query()
// 	s.Require().NoError(err)

// 	var proposer govcliutils.Proposer
// 	json.Unmarshal([]byte(res), &proposer)

// 	s.T().Log(proposer)

// 	s.Require().Equal(1, 2)

// 	s.xplac = qtest.ResetXplac(s.xplac)
// }

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.ChainID = testutil.TestChainId
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
