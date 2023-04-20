package queries

import (
	mbase "github.com/Moonyongjung/xpla.go/core/base"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	tmv1beta1 "cosmossdk.io/api/cosmos/base/tendermint/v1beta1"
	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
)

// Query client for bank module.
func (i IXplaClient) QueryBase() (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcBase(i)
	} else {
		return queryByLcdBase(i)
	}

}

func queryByGrpcBase(i IXplaClient) (string, error) {
	serviceClient := tmservice.NewServiceClient(i.Ixplac.GetGrpcClient())

	switch {
	// Node info
	case i.Ixplac.GetMsgType() == mbase.BaseNodeInfoMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(tmservice.GetNodeInfoRequest)
		res, err = serviceClient.GetNodeInfo(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Syncing
	case i.Ixplac.GetMsgType() == mbase.BaseSyncingMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(tmservice.GetSyncingRequest)
		res, err = serviceClient.GetSyncing(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Latest block
	case i.Ixplac.GetMsgType() == mbase.BaseLatestBlockMsgtype:
		if i.Ixplac.GetRpc() != "" {
			var height *int64
			return queryBlockByRpc(i, height)

		} else {
			convertMsg, _ := i.Ixplac.GetMsg().(tmservice.GetLatestBlockRequest)
			res, err = serviceClient.GetLatestBlock(
				i.Ixplac.GetContext(),
				&convertMsg,
			)
			if err != nil {
				return "", util.LogErr(errors.ErrGrpcRequest, err)
			}
		}

	// Block by height
	case i.Ixplac.GetMsgType() == mbase.BaseBlockByHeightMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(tmservice.GetBlockByHeightRequest)
		if i.Ixplac.GetRpc() != "" {
			height := &convertMsg.Height
			return queryBlockByRpc(i, height)

		} else {
			res, err = serviceClient.GetBlockByHeight(
				i.Ixplac.GetContext(),
				&convertMsg,
			)
			if err != nil {
				return "", util.LogErr(errors.ErrGrpcRequest, err)
			}
		}

	// Latest validator set
	case i.Ixplac.GetMsgType() == mbase.BaseLatestValidatorSetMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(tmservice.GetLatestValidatorSetRequest)
		res, err = serviceClient.GetLatestValidatorSet(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Validator set by height
	case i.Ixplac.GetMsgType() == mbase.BaseValidatorSetByHeightMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(tmservice.GetValidatorSetByHeightRequest)
		res, err = serviceClient.GetValidatorSetByHeight(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err = printProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

const (
	baseNodeInfoLabel      = "node_info"
	baseSyncingLabel       = "syncing"
	baseBlocksLabel        = "blocks"
	baseLatestLabel        = "latest"
	baseValidatorsetsLabel = "validatorsets"
)

func queryByLcdBase(i IXplaClient) (string, error) {
	url := util.MakeQueryLcdUrl(tmv1beta1.Service_ServiceDesc.Metadata.(string))

	switch {
	// Node info
	case i.Ixplac.GetMsgType() == mbase.BaseNodeInfoMsgType:
		url = url + baseNodeInfoLabel

	// Syncing
	case i.Ixplac.GetMsgType() == mbase.BaseSyncingMsgType:
		url = url + baseSyncingLabel

	// Latest block
	case i.Ixplac.GetMsgType() == mbase.BaseLatestBlockMsgtype:
		url = util.MakeQueryLabels("/", baseBlocksLabel, baseLatestLabel)

	// Block by height
	case i.Ixplac.GetMsgType() == mbase.BaseBlockByHeightMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(tmservice.GetBlockByHeightRequest)
		url = util.MakeQueryLabels("/", baseBlocksLabel, util.FromInt64ToString(convertMsg.Height))

	// Latest validator set
	case i.Ixplac.GetMsgType() == mbase.BaseLatestValidatorSetMsgType:
		url = url + util.MakeQueryLabels(baseValidatorsetsLabel, baseLatestLabel)

	// Validator set by height
	case i.Ixplac.GetMsgType() == mbase.BaseValidatorSetByHeightMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(tmservice.GetValidatorSetByHeightRequest)
		url = url + util.MakeQueryLabels(baseValidatorsetsLabel, util.FromInt64ToString(convertMsg.Height))

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func queryBlockByRpc(i IXplaClient, height *int64) (string, error) {
	client, err := cmclient.NewClientFromNode(i.Ixplac.GetRpc())
	if err != nil {
		return "", util.LogErr(errors.ErrGrpcRequest, err)
	}
	res, err := client.Block(i.Ixplac.GetContext(), height)
	if err != nil {
		return "", util.LogErr(errors.ErrGrpcRequest, err)
	}
	out, err := printObjectLegacy(i, res)
	if err != nil {
		return "", util.LogErr(errors.ErrGrpcRequest, err)
	}
	return string(out), nil
}
