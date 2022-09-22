package client

import (
	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	grpc1 "github.com/gogo/protobuf/grpc"
	"google.golang.org/grpc"
)

// The xpla client is a client for performing all functions within the xpla.go library.
// The user mandatorily inputs chain ID and private key.
type XplaClient struct {
	ChainId        string
	PrivateKey     key.PrivateKey
	EncodingConfig params.EncodingConfig
	Grpc           grpc1.ClientConn

	Opts Options

	Module  string
	MsgType string
	Msg     interface{}
	Err     error
}

// Optional parameters of xpla client.
type Options struct {
	AccountNumber  string
	Sequence       string
	BroadcastMode  string
	GasLimit       string
	GasPrice       string
	GasAdjustment  string
	FeeAmount      string
	SignMode       signing.SignMode
	FeeGranter     sdk.AccAddress
	TimeoutHeight  string
	LcdURL         string
	GrpcURL        string
	RpcURL         string
	EvmRpcURL      string
	Pagination     types.Pagination
	OutputDocument string
}

// Make new xpla client.
func NewXplaClient(
	chainId string,
	privKey key.PrivateKey,
) *XplaClient {
	var xplac XplaClient
	return xplac.
		WithChainId(chainId).
		WithPrivateKey(privKey).
		WithEncoding(types.MakeEncodingConfig())
}

// Set options of xpla client.
func (xplac *XplaClient) WithOptions(
	options Options,
) *XplaClient {
	return xplac.
		WithAccountNumber(options.AccountNumber).
		WithBroadcastMode(options.BroadcastMode).
		WithSequence(options.Sequence).
		WithGasLimit(options.GasLimit).
		WithGasPrice(util.DenomRemove(options.GasPrice)).
		WithGasAdjustment(options.GasAdjustment).
		WithFeeAmount(options.FeeAmount).
		WithSignMode(options.SignMode).
		WithFeeGranter(options.FeeGranter).
		WithTimeoutHeight(options.TimeoutHeight).
		WithURL(options.LcdURL).
		WithGrpc(options.GrpcURL).
		WithRpc(options.RpcURL).
		WithEvmRpc(options.EvmRpcURL).
		WithPagination(options.Pagination).
		WithOutputDocument(options.OutputDocument)
}

func init() {
	util.SetChainConfig()
}

// Set chain ID
func (xplac *XplaClient) WithChainId(chainId string) *XplaClient {
	xplac.ChainId = chainId
	return xplac
}

// Set private key
func (xplac *XplaClient) WithPrivateKey(privateKey key.PrivateKey) *XplaClient {
	xplac.PrivateKey = privateKey
	return xplac
}

// Set encoding configuration
func (xplac *XplaClient) WithEncoding(encodingConfig params.EncodingConfig) *XplaClient {
	xplac.EncodingConfig = encodingConfig
	return xplac
}

// Set LCD URL
func (xplac *XplaClient) WithURL(lcdURL string) *XplaClient {
	xplac.Opts.LcdURL = lcdURL
	return xplac
}

// Set GRPC URL to query or broadcast tx
func (xplac *XplaClient) WithGrpc(grpcUrl string) *XplaClient {
	connUrl := util.GrpcUrlParsing(grpcUrl)
	c, _ := grpc.Dial(
		connUrl, grpc.WithInsecure(),
	)
	xplac.Grpc = c
	xplac.Opts.GrpcURL = grpcUrl
	return xplac
}

// Set RPC URL of tendermint core
func (xplac *XplaClient) WithRpc(rpcUrl string) *XplaClient {
	xplac.Opts.RpcURL = rpcUrl
	return xplac
}

// Set RPC URL for evm module
func (xplac *XplaClient) WithEvmRpc(evmRpcUrl string) *XplaClient {
	xplac.Opts.EvmRpcURL = evmRpcUrl
	return xplac
}

// Set broadcast mode
func (xplac *XplaClient) WithBroadcastMode(broadcastMode string) *XplaClient {
	xplac.Opts.BroadcastMode = broadcastMode
	return xplac
}

// Set account number
func (xplac *XplaClient) WithAccountNumber(accountNumber string) *XplaClient {
	xplac.Opts.AccountNumber = accountNumber
	return xplac
}

// Set account sequence
func (xplac *XplaClient) WithSequence(sequence string) *XplaClient {
	xplac.Opts.Sequence = sequence
	return xplac
}

// Set gas limit
func (xplac *XplaClient) WithGasLimit(gasLimit string) *XplaClient {
	xplac.Opts.GasLimit = gasLimit
	return xplac
}

// Set Gas price
func (xplac *XplaClient) WithGasPrice(gasPrice string) *XplaClient {
	xplac.Opts.GasPrice = gasPrice
	return xplac
}

// Set Gas adjustment
func (xplac *XplaClient) WithGasAdjustment(gasAdjustment string) *XplaClient {
	xplac.Opts.GasAdjustment = gasAdjustment
	return xplac
}

// Set fee amount
func (xplac *XplaClient) WithFeeAmount(feeAmount string) *XplaClient {
	xplac.Opts.FeeAmount = feeAmount
	return xplac
}

// Set transaction sign mode
func (xplac *XplaClient) WithSignMode(signMode signing.SignMode) *XplaClient {
	xplac.Opts.SignMode = signMode
	return xplac
}

// Set fee granter
func (xplac *XplaClient) WithFeeGranter(feeGranter sdk.AccAddress) *XplaClient {
	xplac.Opts.FeeGranter = feeGranter
	return xplac
}

// Set timeout block height
func (xplac *XplaClient) WithTimeoutHeight(timeoutHeight string) *XplaClient {
	xplac.Opts.TimeoutHeight = timeoutHeight
	return xplac
}

// Set pagination
func (xplac *XplaClient) WithPagination(pagination types.Pagination) *XplaClient {
	emptyPagination := types.Pagination{}
	if pagination != emptyPagination {
		pageReq, err := core.ReadPageRequest(pagination)
		if err != nil {
			xplac.Err = err
		}
		core.PageRequest = pageReq
	} else {
		core.PageRequest = core.DefaultPagination()
	}

	return xplac
}

// Set output document name
func (xplac *XplaClient) WithOutputDocument(outputDocument string) *XplaClient {
	xplac.Opts.OutputDocument = outputDocument
	return xplac
}
