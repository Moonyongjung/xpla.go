package authz_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Moonyongjung/xpla.go/client"
	"github.com/Moonyongjung/xpla.go/client/queries/qtest"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util/testutil"
	"github.com/gogo/protobuf/jsonpb"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/authz/client/cli"
	authztestutil "github.com/cosmos/cosmos-sdk/x/authz/client/testutil"
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

	val1 := s.network.Validators[0]
	val2 := s.network.Validators[1]

	s.Require().NoError(s.network.WaitForNextBlock())

	_, err := authztestutil.ExecGrant(
		val1,
		[]string{
			val2.Address.String(),
			"send",
			fmt.Sprintf("--%s=100stake", cli.FlagSpendLimit),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val1.Address.String()),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		},
	)

	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

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

func (s *IntegrationTestSuite) TestAuthzGrant() {
	granter := s.network.Validators[0].Address.String()
	grantee := s.network.Validators[1].Address.String()

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		msg1 := types.QueryAuthzGrantMsg{
			Granter: granter,
			Grantee: grantee,
		}
		res1, err := s.xplac.QueryAuthzGrants(msg1).Query()
		s.Require().NoError(err)

		var queryGrantsResponse authz.QueryGrantsResponse
		jsonpb.Unmarshal(strings.NewReader(res1), &queryGrantsResponse)

		s.Require().Equal(1, len(queryGrantsResponse.Grants))

		msg2 := types.QueryAuthzGrantMsg{
			Granter: granter,
		}
		res2, err := s.xplac.QueryAuthzGrants(msg2).Query()
		s.Require().NoError(err)

		var queryGranterGrantsResponse authz.QueryGranterGrantsResponse
		jsonpb.Unmarshal(strings.NewReader(res2), &queryGranterGrantsResponse)

		s.Require().Equal(1, len(queryGranterGrantsResponse.Grants))

		msg3 := types.QueryAuthzGrantMsg{
			Grantee: grantee,
		}
		res3, err := s.xplac.QueryAuthzGrants(msg3).Query()
		s.Require().NoError(err)

		var queryGranteeGrantsResponse authz.QueryGranteeGrantsResponse
		jsonpb.Unmarshal(strings.NewReader(res3), &queryGranteeGrantsResponse)

		s.Require().Equal(1, len(queryGranteeGrantsResponse.Grants))
	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.ChainID = testutil.TestChainId
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
