package params

import (
	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"
	"github.com/gogo/protobuf/proto"

	paramsv1beta1 "cosmossdk.io/api/cosmos/params/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

var out []byte
var res proto.Message
var err error

// Query client for params module.
func QueryParams(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcParams(i)
	} else {
		return queryByLcdParams(i)
	}

}

func queryByGrpcParams(i core.QueryClient) (string, error) {
	queryClient := proposal.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Params subspace
	case i.Ixplac.GetMsgType() == ParamsQuerySubpsaceMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(proposal.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err = core.PrintProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

const (
	paramsParamsLabel = "params"
)

func queryByLcdParams(i core.QueryClient) (string, error) {
	url := util.MakeQueryLcdUrl(paramsv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Params subspace
	case i.Ixplac.GetMsgType() == ParamsQuerySubpsaceMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(proposal.QueryParamsRequest)

		parsedSubspace := convertMsg.Subspace
		parsedKey := convertMsg.Key

		subspace := "?subspace=" + parsedSubspace
		key := "&key=" + parsedKey

		url = url + paramsParamsLabel + subspace + key

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil
}
