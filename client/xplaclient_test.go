package client_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/Moonyongjung/xpla.go/client"
	mbank "github.com/Moonyongjung/xpla.go/core/bank"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
	"github.com/Moonyongjung/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/stretchr/testify/assert"
)

var (
	testBroadcastMode  = "sync"
	testTimeoutHeight  = "1000"
	testLcdUrl         = "https://cube-lcd.xpla.dev"
	testGrpcUrl        = "https://cube-grpc.xpla.dev"
	testRpcUrl         = "https://cube-rpc.xpla.dev"
	testEvmRpcUrl      = "https://cube-evm-rpc.xpla.dev"
	testOutputDocument = "./document.json"
)

func TestNewXplaClient(t *testing.T) {
	// mnemonic, err := key.NewMnemonic()
	// assert.NoError(t, err)

	// priKey, err := key.NewPrivKey(mnemonic)
	// assert.NoError(t, err)

	// addr, err := key.Bech32AddrString(priKey)
	// assert.NoError(t, err)

	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := testutil.RandomAccounts(r, 2)

	from := accounts[0]
	feegranter := accounts[1]
	feeAmount := util.FromUint64ToString(util.MulUint64(
		util.FromStringToUint64(types.DefaultGasLimit),
		util.FromStringToUint64(types.DefaultGasPrice),
	))
	testPagination := types.Pagination{
		PageKey:    "",
		Offset:     0,
		Limit:      0,
		CountTotal: false,
		Reverse:    true,
	}

	newClientOption := client.Options{
		PrivateKey:     from.PrivKey,
		AccountNumber:  util.FromIntToString(types.DefaultAccNum),
		Sequence:       util.FromIntToString(types.DefaultAccSeq),
		BroadcastMode:  testBroadcastMode,
		GasLimit:       types.DefaultGasLimit,
		GasPrice:       types.DefaultGasPrice,
		GasAdjustment:  types.DefaultGasAdjustment,
		FeeAmount:      feeAmount,
		SignMode:       signing.SignMode_SIGN_MODE_DIRECT,
		FeeGranter:     feegranter.Address,
		TimeoutHeight:  testTimeoutHeight,
		LcdURL:         testLcdUrl,
		GrpcURL:        testGrpcUrl,
		RpcURL:         testRpcUrl,
		EvmRpcURL:      testEvmRpcUrl,
		Pagination:     testPagination,
		OutputDocument: testOutputDocument,
	}

	xplac := client.NewXplaClient(testutil.TestChainId).WithOptions(newClientOption)
	xplac.Total()

	totalMsg, err := mbank.MakeTotalSupplyMsg()
	assert.NoError(t, err)

	assert.Equal(t, testutil.TestChainId, xplac.GetChainId())
	assert.Equal(t, from.PrivKey, xplac.GetPrivateKey())
	assert.Equal(t, context.Background(), xplac.GetContext())
	assert.Equal(t, testLcdUrl, xplac.GetLcdURL())
	assert.Equal(t, testGrpcUrl, xplac.GetGrpcUrl())
	assert.Equal(t, xplac.Grpc, xplac.GetGrpcClient())
	assert.Equal(t, testRpcUrl, xplac.GetRpc())
	assert.Equal(t, testEvmRpcUrl, xplac.GetEvmRpc())
	assert.Equal(t, testBroadcastMode, xplac.GetBroadcastMode())
	assert.Equal(t, util.FromIntToString(types.DefaultAccNum), xplac.GetAccountNumber())
	assert.Equal(t, util.FromIntToString(types.DefaultAccSeq), xplac.GetSequence())
	assert.Equal(t, types.DefaultGasLimit, xplac.GetGasLimit())
	assert.Equal(t, types.DefaultGasPrice, xplac.GetGasPrice())
	assert.Equal(t, types.DefaultGasAdjustment, xplac.GetGasAdjustment())
	assert.Equal(t, feeAmount, xplac.GetFeeAmount())
	assert.Equal(t, signing.SignMode_SIGN_MODE_DIRECT, xplac.GetSignMode())
	assert.Equal(t, feegranter.Address, xplac.GetFeeGranter())
	assert.Equal(t, testTimeoutHeight, xplac.GetTimeoutHeight())
	assert.Equal(t, testPagination.Reverse, xplac.GetPagination().Reverse)
	assert.Equal(t, testOutputDocument, xplac.GetOutputDocument())
	assert.Equal(t, mbank.BankModule, xplac.GetModule())
	assert.Equal(t, mbank.BankTotalMsgType, xplac.GetMsgType())
	assert.Equal(t, mbank.BankTotalMsgType, xplac.GetMsgType())
	assert.Equal(t, totalMsg, xplac.GetMsg())

}
