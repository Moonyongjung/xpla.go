package module

import (
	mparams "github.com/Moonyongjung/xpla.go/core/params"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

// Query client for params module.
func (i IXplaClient) QueryParams() (string, error) {
	queryClient := proposal.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Params subspace
	case i.Ixplac.GetMsgType() == mparams.ParamsQuerySubpsaceMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(proposal.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err = printProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
