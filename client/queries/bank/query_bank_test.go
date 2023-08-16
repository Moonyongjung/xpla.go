package bank_test

import (
	"strings"
	"testing"

	"github.com/Moonyongjung/xpla.go/client"
	"github.com/Moonyongjung/xpla.go/client/queries/qtest"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
	"github.com/Moonyongjung/xpla.go/util/testutil"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var (
	validatorNumber = 2

	mainDenom = "amain"
	altDenom  = "aalt"
)

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

	genesisState := s.cfg.GenesisState
	var bankGenesis banktypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[banktypes.ModuleName], &bankGenesis))

	bankGenesis.DenomMetadata = []banktypes.Metadata{
		{
			Description: "main token for test",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "main",
					Exponent: 18,
					Aliases:  []string{"MAIN"},
				},
				{
					Denom:    mainDenom,
					Exponent: 0,
					Aliases:  []string{"attomain"},
				},
			},
			Base:    mainDenom,
			Display: "main",
		},
		{
			Description: "alt token for test",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    altDenom,
					Exponent: 0,
					Aliases:  []string{"attoalt"},
				},
			},
			Base:    altDenom,
			Display: "alt",
		},
	}

	bankGenesisBz, err := s.cfg.Codec.MarshalJSON(&bankGenesis)
	s.Require().NoError(err)
	genesisState[banktypes.ModuleName] = bankGenesisBz
	s.cfg.GenesisState = genesisState

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

func (s *IntegrationTestSuite) TestAllBalancesAndBalance() {
	validator := s.network.Validators[0]
	addr := validator.Address.String()

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		bankBalancesMsg := types.BankBalancesMsg{
			Address: addr,
		}

		res, err := s.xplac.BankBalances(bankBalancesMsg).Query()
		s.Require().NoError(err)

		var allBalancesResponse banktypes.QueryAllBalancesResponse
		jsonpb.Unmarshal(strings.NewReader(res), &allBalancesResponse)

		bal1, err := util.FromStringToBigInt("1000000000000000000000")
		s.Require().NoError(err)
		bal2, err := util.FromStringToBigInt("400000000000000000000")
		s.Require().NoError(err)

		denom1 := "node0token"
		denom2 := "stake"

		s.Require().Equal(bal1, allBalancesResponse.Balances[0].Amount.BigInt())
		s.Require().Equal(bal2, allBalancesResponse.Balances[1].Amount.BigInt())
		s.Require().Equal(denom1, allBalancesResponse.Balances[0].Denom)
		s.Require().Equal(denom2, allBalancesResponse.Balances[1].Denom)

		if i == 1 {
			// LCD not supported
			bankBalancesMsg = types.BankBalancesMsg{
				Address: addr,
				Denom:   denom1,
			}
			res, err = s.xplac.BankBalances(bankBalancesMsg).Query()
			s.Require().NoError(err)

			var balanceResponse banktypes.QueryBalanceResponse
			jsonpb.Unmarshal(strings.NewReader(res), &balanceResponse)

			s.Require().Equal(denom1, balanceResponse.Balance.Denom)
			s.Require().Equal(bal1, balanceResponse.Balance.Amount.BigInt())
		}
	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestDenomMetadata() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.DenomMetadata().Query()
		s.Require().NoError(err)

		var denomsMetadataResponse banktypes.QueryDenomsMetadataResponse
		jsonpb.Unmarshal(strings.NewReader(res), &denomsMetadataResponse)

		s.Require().Equal(2, len(denomsMetadataResponse.Metadatas))
		s.Require().Equal(altDenom, denomsMetadataResponse.Metadatas[0].Base)
		s.Require().Equal(mainDenom, denomsMetadataResponse.Metadatas[1].Base)

		denomMetadataMsg := types.DenomMetadataMsg{
			Denom: mainDenom,
		}
		res, err = s.xplac.DenomMetadata(denomMetadataMsg).Query()
		s.Require().NoError(err)

		var denomMetadataResponse banktypes.QueryDenomMetadataResponse
		jsonpb.Unmarshal(strings.NewReader(res), &denomMetadataResponse)

		s.Require().Equal(mainDenom, denomMetadataResponse.Metadata.Base)
	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestBankTotal() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		bal1, err := util.FromStringToBigInt("1000000000000000000000")
		s.Require().NoError(err)
		bal2, err := util.FromStringToBigInt("1000000000000000000000")
		s.Require().NoError(err)
		bal3, err := util.FromStringToBigInt("1000000020597259368396")
		s.Require().NoError(err)

		denom1 := "node0token"
		denom2 := "node1token"
		denom3 := "stake"

		res, err := s.xplac.Total().Query()
		s.Require().NoError(err)

		var totalSupplyResponse banktypes.QueryTotalSupplyResponse
		jsonpb.Unmarshal(strings.NewReader(res), &totalSupplyResponse)

		s.Require().Equal(denom1, totalSupplyResponse.Supply[0].Denom)
		s.Require().Equal(bal1, totalSupplyResponse.Supply[0].Amount.BigInt())
		s.Require().Equal(denom2, totalSupplyResponse.Supply[1].Denom)
		s.Require().Equal(bal2, totalSupplyResponse.Supply[1].Amount.BigInt())
		s.Require().Equal(denom3, totalSupplyResponse.Supply[2].Denom)
		s.Require().Equal(bal3, totalSupplyResponse.Supply[2].Amount.BigInt())

		totalMsg := types.TotalMsg{
			Denom: denom1,
		}
		res, err = s.xplac.Total(totalMsg).Query()
		s.Require().NoError(err)

		var supplyOfResponse banktypes.QuerySupplyOfResponse
		jsonpb.Unmarshal(strings.NewReader(res), &supplyOfResponse)

		s.Require().Equal(denom1, supplyOfResponse.Amount.Denom)
		s.Require().Equal(bal1, supplyOfResponse.Amount.Amount.BigInt())
	}
	s.xplac = qtest.ResetXplac(s.xplac)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.ChainID = testutil.TestChainId
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
