package util

import (
	"context"
	"net/http"

	"github.com/Moonyongjung/xpla.go/types"
	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/ethclient"
	erpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/evmos/ethermint/crypto/hd"
)

const (
	BackendFile   = "file"
	BackendMemory = "memory"
)

// Provide cosmos sdk client.
func NewClient() (cmclient.Context, error) {
	clientCtx := cmclient.Context{}
	encodingConfig := MakeEncodingConfig()
	clientKeyring, err := NewKeyring(BackendMemory, "")
	if err != nil {
		return cmclient.Context{}, err
	}

	clientCtx = clientCtx.
		WithTxConfig(encodingConfig.TxConfig).
		WithCodec(encodingConfig.Marshaler).
		WithLegacyAmino(encodingConfig.Amino).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithKeyringOptions(hd.EthSecp256k1Option()).
		WithKeyring(clientKeyring).
		WithAccountRetriever(authtypes.AccountRetriever{})

	return clientCtx, nil
}

const (
	DefaultEvmGasLimit         = "21000"
	DefaultSolidityValue       = "0"
	DefaultEvmTxReceiptTimeout = 100
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
		return nil, LogErr("default transport pointer err")
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

// Provide cosmos sdk keyring
func NewKeyring(backendType string, keyringPath string) (keyring.Keyring, error) {
	if backendType == BackendMemory {
		k, err := keyring.New(
			types.XplaToolDefaultName,
			keyring.BackendMemory,
			"",
			nil,
			hd.EthSecp256k1Option(),
		)
		if err != nil {
			return nil, err
		}

		return k, nil

	} else if backendType == BackendFile {
		k, err := keyring.New(
			types.XplaToolDefaultName,
			keyring.BackendFile,
			keyringPath,
			nil,
			hd.EthSecp256k1Option(),
		)
		if err != nil {
			return nil, err
		}

		return k, nil
	} else {
		return nil, LogErr("invalid keyring backend type")
	}
}

// Provide cosmos sdk tx factory.
func NewFactory(clientCtx cmclient.Context) tx.Factory {
	txFactory := tx.Factory{}.
		WithTxConfig(clientCtx.TxConfig).
		WithKeybase(clientCtx.Keyring).
		WithAccountRetriever(clientCtx.AccountRetriever)

	return txFactory
}
