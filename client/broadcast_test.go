package client_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/Moonyongjung/xpla.go/client"
	"github.com/Moonyongjung/xpla.go/client/xplago_helper"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
	"github.com/Moonyongjung/xpla.go/util/testutil"
	"github.com/gogo/protobuf/jsonpb"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
)

var (
	validatorNumberForBroadcast = 2

	testSendAmount = "1000"
)

type ClientBroadcastTestSuite struct {
	suite.Suite

	xplac    *client.XplaClient
	apis     []string
	accounts []simtypes.Account

	cfg     network.Config
	network *network.Network
}

func NewClientBroadcastTestSuite(cfg network.Config) *ClientBroadcastTestSuite {
	return &ClientBroadcastTestSuite{cfg: cfg}
}

func (s *ClientBroadcastTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	src := rand.NewSource(1)
	r := rand.New(src)
	s.accounts = testutil.RandomSecp256k1Accounts(r, 2)

	balanceBigInt, err := util.FromStringToBigInt("1000000000000000000000000000")
	s.Require().NoError(err)

	genesisState := s.cfg.GenesisState

	// add genesis account
	var authGenesis authtypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[authtypes.ModuleName], &authGenesis))

	var genAccounts []authtypes.GenesisAccount

	genAccounts = append(genAccounts, authtypes.NewBaseAccount(s.accounts[0].Address, nil, 0, 0))
	genAccounts = append(genAccounts, authtypes.NewBaseAccount(s.accounts[1].Address, nil, 0, 0))

	accounts, err := authtypes.PackAccounts(genAccounts)
	s.Require().NoError(err)

	authGenesis.Accounts = accounts

	authGenesisBz, err := s.cfg.Codec.MarshalJSON(&authGenesis)
	s.Require().NoError(err)
	genesisState[authtypes.ModuleName] = authGenesisBz

	// add balances
	var bankGenesis banktypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[banktypes.ModuleName], &bankGenesis))

	bankGenesis.Balances = []banktypes.Balance{
		{
			Address: s.accounts[0].Address.String(),
			Coins: sdk.Coins{
				sdk.NewCoin(types.XplaDenom, sdk.NewIntFromBigInt(balanceBigInt)),
			},
		},
		{
			Address: s.accounts[1].Address.String(),
			Coins: sdk.Coins{
				sdk.NewCoin(types.XplaDenom, sdk.NewIntFromBigInt(balanceBigInt)),
			},
		},
	}

	bankGenesisBz, err := s.cfg.Codec.MarshalJSON(&bankGenesis)
	s.Require().NoError(err)
	genesisState[banktypes.ModuleName] = bankGenesisBz

	// staking
	var stakingGenesis stakingtypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesis))

	stakingGenesis.Params.BondDenom = types.XplaDenom

	stakingGenesisBz, err := s.cfg.Codec.MarshalJSON(&stakingGenesis)
	s.Require().NoError(err)
	genesisState[stakingtypes.ModuleName] = stakingGenesisBz

	s.cfg.GenesisState = genesisState
	s.network = network.New(s.T(), s.cfg)
	s.Require().NoError(s.network.WaitForNextBlock())

	s.xplac = xplago_helper.NewTestXplaClient()
	s.apis = []string{
		s.network.Validators[0].APIAddress,
		s.network.Validators[0].AppConfig.GRPC.Address,
	}
}

func (s *ClientBroadcastTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *ClientBroadcastTestSuite) TestBroadcast() {
	from := s.accounts[0]
	to := s.accounts[1]

	s.xplac.WithPrivateKey(from.PrivKey)

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
			newSeq := util.FromStringToInt(s.xplac.GetSequence()) + 1
			s.xplac.WithSequence(util.FromIntToString(newSeq))
		}

		// check before send
		bankBalancesMsg := types.BankBalancesMsg{
			Address: to.Address.String(),
		}
		beforeToRes, err := s.xplac.BankBalances(bankBalancesMsg).Query()
		s.Require().NoError(err)

		var beforeQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
		jsonpb.Unmarshal(strings.NewReader(beforeToRes), &beforeQueryAllBalancesResponse)

		// broadcast transaction - bank send
		bankSendMsg := types.BankSendMsg{
			FromAddress: from.Address.String(),
			ToAddress:   to.Address.String(),
			Amount:      testSendAmount,
		}
		txbytes, err := s.xplac.BankSend(bankSendMsg).CreateAndSignTx()
		s.Require().NoError(err)

		_, err = s.xplac.Broadcast(txbytes)
		s.Require().NoError(err)
		s.Require().NoError(s.network.WaitForNextBlock())

		// check after send
		bankBalancesMsg = types.BankBalancesMsg{
			Address: to.Address.String(),
		}
		afterToRes, err := s.xplac.BankBalances(bankBalancesMsg).Query()
		s.Require().NoError(err)

		var afterQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
		jsonpb.Unmarshal(strings.NewReader(afterToRes), &afterQueryAllBalancesResponse)

		s.Require().Equal(
			testSendAmount,
			afterQueryAllBalancesResponse.Balances[0].Amount.Sub(beforeQueryAllBalancesResponse.Balances[0].Amount).String(),
		)
	}
	s.xplac = xplago_helper.ResetXplac(s.xplac)
}

func TestClientBroadcastTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.ChainID = testutil.TestChainId
	cfg.BondDenom = types.XplaDenom
	cfg.MinGasPrices = fmt.Sprintf("0%s", types.XplaDenom)
	cfg.NumValidators = validatorNumberForBroadcast
	suite.Run(t, NewClientBroadcastTestSuite(cfg))
}
