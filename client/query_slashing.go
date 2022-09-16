package client

import (
	"context"

	mslashing "github.com/Moonyongjung/xpla.go/core/slashing"
	"github.com/Moonyongjung/xpla.go/util"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// Query client for slashing module.
func querySlashing(xplac *XplaClient) (string, error) {
	queryClient := slashingtypes.NewQueryClient(xplac.Grpc)

	switch {
	// Slashing parameters
	case xplac.MsgType == mslashing.SlahsingQuerySlashingParamsMsgType:
		convertMsg, _ := xplac.Msg.(slashingtypes.QueryParamsRequest)
		res, err = queryClient.Params(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Slashing signing information
	case xplac.MsgType == mslashing.SlashingQuerySigningInfosMsgType:
		convertMsg, _ := xplac.Msg.(slashingtypes.QuerySigningInfosRequest)
		res, err = queryClient.SigningInfos(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Slashing signing information
	case xplac.MsgType == mslashing.SlashingQuerySigningInfoMsgType:
		convertMsg, _ := xplac.Msg.(slashingtypes.QuerySigningInfoRequest)
		res, err = queryClient.SigningInfo(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err = printProto(xplac, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
