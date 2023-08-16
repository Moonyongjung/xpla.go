package mint_test

import (
	"strings"
	"testing"

	"github.com/Moonyongjung/xpla.go/client"
	"github.com/Moonyongjung/xpla.go/client/queries/qtest"
	"github.com/Moonyongjung/xpla.go/util/testutil"
	"github.com/gogo/protobuf/jsonpb"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/suite"
)

var validatorNumber = 2

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
	s.xplac = qtest.NewTestXplaClient()
	s.apis = []string{
		s.network.Validators[0].APIAddress,
		s.network.Validators[0].AppConfig.GRPC.Address,
	}

	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestParams() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.MintParams().Query()
		s.Require().NoError(err)

		var queryParamsResponse minttypes.QueryParamsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryParamsResponse)

		s.Require().Equal("stake", queryParamsResponse.Params.MintDenom)
		s.Require().Equal("0.130000000000000000", queryParamsResponse.Params.InflationRateChange.String())
		s.Require().Equal("0.200000000000000000", queryParamsResponse.Params.InflationMax.String())
		s.Require().Equal("0.070000000000000000", queryParamsResponse.Params.InflationMin.String())
		s.Require().Equal("0.670000000000000000", queryParamsResponse.Params.GoalBonded.String())
		s.Require().Equal(uint64(6311520), queryParamsResponse.Params.BlocksPerYear)
	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestInflation() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.Inflation().Query()
		s.Require().NoError(err)

		var queryInflationResponse minttypes.QueryInflationResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryInflationResponse)

		s.Require().Equal("0.130000014448822130", queryInflationResponse.Inflation.String())
	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestAnnualProvisions() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.AnnualProvisions().Query()
		s.Require().NoError(err)

		var queryAnnualProvisionsResponse minttypes.QueryAnnualProvisionsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryAnnualProvisionsResponse)

		s.Require().Equal("130000014448822130000.000000000000000000", queryAnnualProvisionsResponse.AnnualProvisions.String())

	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.ChainID = testutil.TestChainId
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
