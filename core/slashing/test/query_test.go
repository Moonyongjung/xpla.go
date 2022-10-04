package slashing_test

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/Moonyongjung/xpla.go/core"
	mslashing "github.com/Moonyongjung/xpla.go/core/slashing"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/slashing/testslashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

type TestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *xapp.XplaApp
	queryClient slashingtypes.QueryClient
	addrDels    []sdk.AccAddress
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

func (suite *TestSuite) SetupTest() {
	app := util.Setup(false, 5)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	app.SlashingKeeper.SetParams(ctx, testslashing.TestParams())

	addrDels := AddTestAddrsIncremental(app, ctx, 2, app.StakingKeeper.TokensFromConsensusPower(ctx, 200))

	info1 := slashingtypes.NewValidatorSigningInfo(sdk.ConsAddress(addrDels[0]), int64(4), int64(3),
		time.Unix(2, 0), false, int64(10))
	info2 := slashingtypes.NewValidatorSigningInfo(sdk.ConsAddress(addrDels[1]), int64(5), int64(4),
		time.Unix(2, 0), false, int64(10))

	app.SlashingKeeper.SetValidatorSigningInfo(ctx, sdk.ConsAddress(addrDels[0]), info1)
	app.SlashingKeeper.SetValidatorSigningInfo(ctx, sdk.ConsAddress(addrDels[1]), info2)

	suite.app = app
	suite.ctx = ctx
	suite.addrDels = addrDels

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	slashingtypes.RegisterQueryServer(queryHelper, app.SlashingKeeper)
	queryClient := slashingtypes.NewQueryClient(queryHelper)
	suite.queryClient = queryClient
}

func (suite *TestSuite) TestGRPCQueryParams() {
	queryClient := suite.queryClient
	msg, _ := mslashing.MakeQuerySlashingParamsMsg()
	paramsResp, err := queryClient.Params(context.Background(), &msg)

	suite.NoError(err)
	suite.Equal(testslashing.TestParams(), paramsResp.Params)
}

func (suite *TestSuite) TestGRPCSigningInfo() {
	queryClient := suite.queryClient

	encodingConfig := util.MakeEncodingConfig()
	signingInfoMsg := types.SigningInfoMsg{
		ConsAddr: "",
	}
	msg, _ := mslashing.MakeQuerySigningInfoMsg(signingInfoMsg, encodingConfig)
	infoResp, err := queryClient.SigningInfo(context.Background(), &msg)
	suite.Error(err)
	suite.Nil(infoResp)

	consAddr := sdk.ConsAddress(suite.addrDels[0])
	info, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, consAddr)
	suite.True(found)

	signingInfoMsg = types.SigningInfoMsg{
		ConsAddr: consAddr.String(),
	}
	msg, _ = mslashing.MakeQuerySigningInfoMsg(signingInfoMsg, encodingConfig)
	infoResp, err = queryClient.SigningInfo(context.Background(),
		&msg)
	suite.NoError(err)
	suite.Equal(info, infoResp.ValSigningInfo)

	m1, _ := key.NewMnemonic()
	k1, _ := key.NewPrivKey(m1)
	addr1 := util.GetAddrByPrivKey(k1)
	pubk1 := k1.PubKey().String()

	newValInfo := slashingtypes.NewValidatorSigningInfo(sdk.ConsAddress(addr1), int64(5), int64(4),
		time.Unix(2, 0), false, int64(10))

	suite.app.SlashingKeeper.SetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(addr1), newValInfo)
	getValInfo, _ := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, consAddr)

	signingInfoMsg = types.SigningInfoMsg{
		ConsPubKey: pubk1,
	}
	msg, _ = mslashing.MakeQuerySigningInfoMsg(signingInfoMsg, encodingConfig)
	infoResp, err = queryClient.SigningInfo(context.Background(),
		&msg)
	suite.NoError(err)
	suite.Equal(getValInfo, infoResp.ValSigningInfo)
}

func (suite *TestSuite) TestGRPCSigningInfos() {
	queryClient := suite.queryClient

	var signingInfos []slashingtypes.ValidatorSigningInfo

	suite.app.SlashingKeeper.IterateValidatorSigningInfos(suite.ctx, func(consAddr sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo) (stop bool) {
		signingInfos = append(signingInfos, info)
		return false
	})

	// verify all values are returned without pagination
	var infoResp, err = queryClient.SigningInfos(context.Background(),
		&slashingtypes.QuerySigningInfosRequest{Pagination: nil})
	suite.NoError(err)
	suite.Equal(signingInfos, infoResp.Info)

	pagination := types.Pagination{
		Limit:      1,
		CountTotal: true,
	}
	pr, _ := core.ReadPageRequest(pagination)
	core.PageRequest = pr

	msg, _ := mslashing.MakeQuerySigningInfosMsg()

	infoResp, err = queryClient.SigningInfos(context.Background(), &msg)
	suite.NoError(err)
	suite.Len(infoResp.Info, 1)
	suite.Equal(signingInfos[0], infoResp.Info[0])
	suite.NotNil(infoResp.Pagination.NextKey)
	suite.Equal(uint64(2), infoResp.Pagination.Total)
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(app *xapp.XplaApp, ctx sdk.Context, accNum int, accAmt sdk.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createIncrementalAccounts)
}

func addTestAddrs(app *xapp.XplaApp, ctx sdk.Context, accNum int, accAmt sdk.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))

	for _, addr := range testAddrs {
		initAccountWithCoins(app, ctx, addr, initCoins)
	}

	return testAddrs
}

func initAccountWithCoins(app *xapp.XplaApp, ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) {
	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
}

// createIncrementalAccounts is a strategy used by addTestAddrs() in order to generated addresses in ascending order.
func createIncrementalAccounts(accNum int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (accNum + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") // base address string

		buffer.WriteString(numString) // adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addr, _ := util.TestAddr(buffer.String(), bech)

		addresses = append(addresses, addr)
		buffer.Reset()
	}

	return addresses
}
