package module

import (
	mmint "github.com/Moonyongjung/xpla.go/core/mint"
	"github.com/Moonyongjung/xpla.go/util"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// Query client for mint module.
func (i IXplaClient) QueryMint() (string, error) {
	queryClient := minttypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Mint parameters
	case i.Ixplac.GetMsgType() == mmint.MintQueryMintParamsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(minttypes.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Mint inflation
	case i.Ixplac.GetMsgType() == mmint.MintQueryInflationMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(minttypes.QueryInflationRequest)
		res, err = queryClient.Inflation(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Mint annual provisions
	case i.Ixplac.GetMsgType() == mmint.MintQueryAnnualProvisionsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(minttypes.QueryAnnualProvisionsRequest)
		res, err = queryClient.AnnualProvisions(
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
