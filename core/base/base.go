package base

import (
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
)

// (Query) make msg - node info
func MakeBaseNodeInfoMsg() (tmservice.GetNodeInfoRequest, error) {
	return tmservice.GetNodeInfoRequest{}, nil
}

// (Query) make msg - syncing
func MakeBaseSyncingMsg() (tmservice.GetSyncingRequest, error) {
	return tmservice.GetSyncingRequest{}, nil
}

// (Query) make msg - latest block
func MakeBaseLatestBlockMsg() (tmservice.GetLatestBlockRequest, error) {
	return tmservice.GetLatestBlockRequest{}, nil
}

// (Query) make msg - get block by height
func MakeBaseBlockByHeightMsg(blockMsg types.BlockMsg) (tmservice.GetBlockByHeightRequest, error) {
	return tmservice.GetBlockByHeightRequest{
		Height: util.FromStringToInt64(blockMsg.Height),
	}, nil
}

// (Query) make msg - latest validator set
func MakeLatestValidatorSetMsg() (tmservice.GetLatestValidatorSetRequest, error) {
	return tmservice.GetLatestValidatorSetRequest{}, nil
}

// (Query) make msg - latest validator set
func MakeValidatorSetByHeightMsg(validatorSetMsg types.ValidatorSetMsg) (tmservice.GetValidatorSetByHeightRequest, error) {
	return tmservice.GetValidatorSetByHeightRequest{
		Height: util.FromStringToInt64(validatorSetMsg.Height),
	}, nil
}
