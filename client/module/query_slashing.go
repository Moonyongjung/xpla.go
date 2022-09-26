package module

import (
	mslashing "github.com/Moonyongjung/xpla.go/core/slashing"
	"github.com/Moonyongjung/xpla.go/util"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// Query client for slashing module.
func (i IXplaClient) QuerySlashing() (string, error) {
	queryClient := slashingtypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Slashing parameters
	case i.Ixplac.GetMsgType() == mslashing.SlahsingQuerySlashingParamsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(slashingtypes.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Slashing signing information
	case i.Ixplac.GetMsgType() == mslashing.SlashingQuerySigningInfosMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(slashingtypes.QuerySigningInfosRequest)
		res, err = queryClient.SigningInfos(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Slashing signing information
	case i.Ixplac.GetMsgType() == mslashing.SlashingQuerySigningInfoMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(slashingtypes.QuerySigningInfoRequest)
		res, err = queryClient.SigningInfo(
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
