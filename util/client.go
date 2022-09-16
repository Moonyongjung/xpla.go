package util

import (
	"github.com/Moonyongjung/xpla.go/types"
	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/evmos/ethermint/crypto/hd"
)

const (
	BackendFile   = "file"
	BackendMemory = "memory"
)

// Provide cosmos sdk client.
func NewClient() (cmclient.Context, error) {
	clientCtx := cmclient.Context{}
	encodingConfig := types.MakeEncodingConfig()
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
