package client

import (
	mmint "github.com/Moonyongjung/xpla.go/core/mint"
	"github.com/Moonyongjung/xpla.go/util"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// Query client for mint module.
func queryMint(xplac *XplaClient) (string, error) {
	queryClient := minttypes.NewQueryClient(xplac.Grpc)

	switch {
	// Mint parameters
	case xplac.MsgType == mmint.MintQueryMintParamsMsgType:
		convertMsg, _ := xplac.Msg.(minttypes.QueryParamsRequest)
		res, err = queryClient.Params(
			xplac.Context,
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Mint inflation
	case xplac.MsgType == mmint.MintQueryInflationMsgType:
		convertMsg, _ := xplac.Msg.(minttypes.QueryInflationRequest)
		res, err = queryClient.Inflation(
			xplac.Context,
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Mint annual provisions
	case xplac.MsgType == mmint.MintQueryAnnualProvisionsMsgType:
		convertMsg, _ := xplac.Msg.(minttypes.QueryAnnualProvisionsRequest)
		res, err = queryClient.AnnualProvisions(
			xplac.Context,
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
