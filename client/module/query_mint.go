package module

import (
	mmint "github.com/Moonyongjung/xpla.go/core/mint"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	mintv1beta1 "cosmossdk.io/api/cosmos/mint/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// Query client for mint module.
func (i IXplaClient) QueryMint() (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcMint(i)
	} else {
		return queryByLcdMint(i)
	}
}

func queryByGrpcMint(i IXplaClient) (string, error) {
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
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Mint inflation
	case i.Ixplac.GetMsgType() == mmint.MintQueryInflationMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(minttypes.QueryInflationRequest)
		res, err = queryClient.Inflation(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Mint annual provisions
	case i.Ixplac.GetMsgType() == mmint.MintQueryAnnualProvisionsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(minttypes.QueryAnnualProvisionsRequest)
		res, err = queryClient.AnnualProvisions(
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
	mintParamsLabel           = "params"
	mintInflationLabel        = "inflation"
	mintAnnualProvisionsLabel = "annual_provisions"
)

func queryByLcdMint(i IXplaClient) (string, error) {
	url := util.MakeQueryLcdUrl(mintv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Mint parameters
	case i.Ixplac.GetMsgType() == mmint.MintQueryMintParamsMsgType:
		url = url + mintParamsLabel

	// Mint inflation
	case i.Ixplac.GetMsgType() == mmint.MintQueryInflationMsgType:
		url = url + mintInflationLabel

	// Mint annual provisions
	case i.Ixplac.GetMsgType() == mmint.MintQueryAnnualProvisionsMsgType:
		url = url + mintAnnualProvisionsLabel

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil
}
