package module

import (
	"context"

	mevm "github.com/Moonyongjung/xpla.go/core/evm"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/util"
	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/gogo/protobuf/grpc"
	"github.com/gogo/protobuf/proto"
)

var out []byte
var res proto.Message
var err error

type IXplaClient struct {
	Ixplac ModuleClient
}

type ModuleClient interface {
	GetChainId() string
	GetPrivateKey() key.PrivateKey
	GetEncoding() params.EncodingConfig
	GetContext() context.Context
	GetURL() string
	GetGrpcUrl() string
	GetGrpcClient() grpc.ClientConn
	GetRpc() string
	GetEvmRpc() string
	GetBroadcastMode() string
	GetAccountNumber() string
	GetSequence() string
	GetGasLimit() string
	GetGasPrice() string
	GetGasAdjustment() string
	GetFeeAmount() string
	GetSignMode() signing.SignMode
	GetFeeGranter() sdk.AccAddress
	GetTimeoutHeight() string
	GetPagination() *query.PageRequest
	GetOutputDocument() string
	GetModule() string
	GetMsg() interface{}
	GetMsgType() string
}

func NewIXplaClient(moduleClient ModuleClient) *IXplaClient {
	return &IXplaClient{Ixplac: moduleClient}
}

// For invoke(as execute) contract, parameters are packed by using ABI.
func GetAbiPack(callName string, args ...interface{}) ([]byte, error) {
	contractAbi, err := mevm.XplaSolContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	var abiByteData []byte

	if args == nil {
		abiByteData, err = contractAbi.Pack(callName)
		if err != nil {
			return nil, err
		}
	} else {
		abiByteData, err = contractAbi.Pack(callName, args...)
		if err != nil {
			return nil, err
		}
	}

	return abiByteData, nil
}

// After call(as query) solidity contract, the response of chain is unpacked by ABI.
func GetAbiUnpack(callName string, data []byte) ([]interface{}, error) {
	contractAbi, err := mevm.XplaSolContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	unpacked, err := contractAbi.Unpack(callName, data)
	if err != nil {
		return nil, err
	}

	return unpacked, nil
}

// Print protobuf message by using cosmos sdk codec.
func printProto(i IXplaClient, toPrint proto.Message) ([]byte, error) {
	out, err := i.Ixplac.GetEncoding().Marshaler.MarshalJSON(toPrint)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Print object by using cosmos sdk legacy amino.
func printObjectLegacy(i IXplaClient, toPrint interface{}) ([]byte, error) {
	out, err := i.Ixplac.GetEncoding().Amino.MarshalJSON(toPrint)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// For auth module and gov module, make cosmos sdk client for querying.
func clientForQuery(i IXplaClient) (cmclient.Context, error) {
	client, err := cmclient.NewClientFromNode(i.Ixplac.GetRpc())
	if err != nil {
		return cmclient.Context{}, err
	}

	clientCtx, err := util.NewClient()
	if err != nil {
		return cmclient.Context{}, err
	}

	clientCtx = clientCtx.
		WithNodeURI(i.Ixplac.GetRpc()).
		WithClient(client)

	return clientCtx, nil
}
