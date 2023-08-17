package client_test

import (
	"fmt"
	"testing"

	"github.com/Moonyongjung/xpla.go/client"
	"github.com/Moonyongjung/xpla.go/client/xplago_helper"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/client/flags"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	"github.com/stretchr/testify/suite"
)

var (
	validatorNumber = 1
)

type ClientHTTPTestSuite struct {
	suite.Suite

	xplac       *client.XplaClient
	apis        []string
	fromAddr    string
	fromPrivKey cryptotypes.PrivKey

	cfg     network.Config
	network *network.Network
}

func NewClientHTTPTestSuite(cfg network.Config) *ClientHTTPTestSuite {
	return &ClientHTTPTestSuite{cfg: cfg}
}

func (s *ClientHTTPTestSuite) SetupSuite() {
	mnemonic, err := key.NewMnemonic()
	s.Require().NoError(err)

	s.fromPrivKey, err = testutil.NewTestSecpPrivKey(mnemonic)
	s.Require().NoError(err)

	s.fromAddr, err = key.Bech32AddrString(s.fromPrivKey)
	s.Require().NoError(err)

	newAddr, err := sdk.AccAddressFromBech32(s.fromAddr)
	s.Require().NoError(err)

	s.network = network.New(s.T(), s.cfg)
	s.Require().NoError(s.network.WaitForNextBlock())

	val := s.network.Validators[0]

	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	s.xplac = xplago_helper.NewTestXplaClient()
	s.apis = []string{
		s.network.Validators[0].APIAddress,
		s.network.Validators[0].AppConfig.GRPC.Address,
	}
}

func (s *ClientHTTPTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *ClientHTTPTestSuite) TestLoadAccount() {
	val := s.network.Validators[0].Address

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.LoadAccount(val)
		s.Require().NoError(err)
		s.Require().Equal(val.String(), res.GetAddress().String())
	}
	s.xplac = xplago_helper.ResetXplac(s.xplac)
}

func (s *ClientHTTPTestSuite) TestSimulate() {
	val1 := s.network.Validators[0].Address
	s.xplac.WithPrivateKey(s.fromPrivKey)

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		authzGrantMsg := types.AuthzGrantMsg{
			Granter:           s.fromAddr,
			Grantee:           val1.String(),
			AuthorizationType: "send",
			SpendLimit:        "1000",
		}

		xplac := s.xplac.AuthzGrant(authzGrantMsg)
		s.Require().NoError(xplac.Err)

		builder := xplac.EncodingConfig.TxConfig.NewTxBuilder()

		convertMsg, _ := xplac.Msg.(authz.MsgGrant)
		builder.SetMsgs(&convertMsg)

		simulate, err := xplac.Simulate(builder)
		s.Require().NoError(err)
		s.Require().Equal(uint64(55608), simulate.GasInfo.GasUsed)

	}
	s.xplac = xplago_helper.ResetXplac(s.xplac)
}

func TestClientHTTPTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.ChainID = testutil.TestChainId
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewClientHTTPTestSuite(cfg))
}
