package client

import (
	"context"

	mparams "github.com/Moonyongjung/xpla.go/core/params"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

// Query client for params module.
func queryParams(xplac *XplaClient) (string, error) {
	queryClient := proposal.NewQueryClient(xplac.Grpc)

	switch {
	// Params subspace
	case xplac.MsgType == mparams.ParamsQuerySubpsaceMsgType:
		convertMsg, _ := xplac.Msg.(proposal.QueryParamsRequest)
		res, err = queryClient.Params(
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
