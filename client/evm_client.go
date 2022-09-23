package client

import (
	"context"
	"net/http"

	"github.com/Moonyongjung/xpla.go/util"
	"github.com/ethereum/go-ethereum/ethclient"
	erpc "github.com/ethereum/go-ethereum/rpc"
)

const (
	defaultEvmGasLimit         = "21000"
	defaultSolidityValue       = "0"
	defaultEvmTxReceiptTimeout = 100
)

type EvmClient struct {
	Ctx    context.Context
	Client *ethclient.Client
}

// Make new evm client using RPC URL which normally TCP port number is 8545.
// It supports that sending transaction, contract deployment, executing/querying contract and etc.
func NewEvmClient(evmRpcUrl string, ctx context.Context) (*EvmClient, error) {
	// Target blockchain node URL
	httpDefaultTransport := http.DefaultTransport
	defaultTransportPointer, ok := httpDefaultTransport.(*http.Transport)
	if !ok {
		return nil, util.LogErr("default transport pointer err")
	}
	defaultTransport := *defaultTransportPointer
	defaultTransport.DisableKeepAlives = true

	httpClient := &http.Client{Transport: &defaultTransport}
	rpcClient, err := erpc.DialHTTPWithClient(evmRpcUrl, httpClient)
	if err != nil {
		return nil, err
	}

	ethClient := ethclient.NewClient(rpcClient)

	return &EvmClient{ctx, ethClient}, nil
}
