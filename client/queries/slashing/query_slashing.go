package slashing

import (
	"github.com/Moonyongjung/xpla.go/client/queries"
	mslashing "github.com/Moonyongjung/xpla.go/core/slashing"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"
	"github.com/gogo/protobuf/proto"

	slashingv1beta1 "cosmossdk.io/api/cosmos/slashing/v1beta1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

var out []byte
var res proto.Message
var err error

// Query client for slashing module.
func QuerySlashing(i queries.IXplaClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcSlashing(i)
	} else {
		return queryByLcdSlashing(i)
	}
}

func queryByGrpcSlashing(i queries.IXplaClient) (string, error) {
	queryClient := slashingtypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Slashing parameters
	case i.Ixplac.GetMsgType() == mslashing.SlashingQuerySlashingParamsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(slashingtypes.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Slashing signing information
	case i.Ixplac.GetMsgType() == mslashing.SlashingQuerySigningInfosMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(slashingtypes.QuerySigningInfosRequest)
		res, err = queryClient.SigningInfos(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Slashing signing information
	case i.Ixplac.GetMsgType() == mslashing.SlashingQuerySigningInfoMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(slashingtypes.QuerySigningInfoRequest)
		res, err = queryClient.SigningInfo(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err = queries.PrintProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

const (
	slashingParamsLabel       = "params"
	slashingSigningInfosLabel = "signing_infos"
)

func queryByLcdSlashing(i queries.IXplaClient) (string, error) {
	url := util.MakeQueryLcdUrl(slashingv1beta1.Query_ServiceDesc.Metadata.(string))
	switch {
	// Slashing parameters
	case i.Ixplac.GetMsgType() == mslashing.SlashingQuerySlashingParamsMsgType:
		url = url + slashingParamsLabel

	// Slashing signing information
	case i.Ixplac.GetMsgType() == mslashing.SlashingQuerySigningInfosMsgType:
		url = url + slashingSigningInfosLabel

	// Slashing signing information
	case i.Ixplac.GetMsgType() == mslashing.SlashingQuerySigningInfoMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(slashingtypes.QuerySigningInfoRequest)

		url = url + util.MakeQueryLabels(slashingSigningInfosLabel, convertMsg.ConsAddress)

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil
}
