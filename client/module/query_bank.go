package module

import (
	mbank "github.com/Moonyongjung/xpla.go/core/bank"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	bankv1beta1 "cosmossdk.io/api/cosmos/bank/v1beta1"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// Query client for bank module.
func (i IXplaClient) QueryBank() (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcBank(i)
	} else {
		return queryByLcdBank(i)
	}

}

func queryByGrpcBank(i IXplaClient) (string, error) {
	queryClient := banktypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Bank balances
	case i.Ixplac.GetMsgType() == mbank.BankAllBalancesMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(*banktypes.QueryAllBalancesRequest)
		res, err = queryClient.AllBalances(
			i.Ixplac.GetContext(),
			convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Bank balance
	case i.Ixplac.GetMsgType() == mbank.BankBalanceMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(*banktypes.QueryBalanceRequest)
		res, err = queryClient.Balance(
			i.Ixplac.GetContext(),
			convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Bank denominations metadata
	case i.Ixplac.GetMsgType() == mbank.BankDenomsMetadataMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(banktypes.QueryDenomsMetadataRequest)
		res, err = queryClient.DenomsMetadata(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Bank denomination metadata
	case i.Ixplac.GetMsgType() == mbank.BankDenomMetadataMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(banktypes.QueryDenomMetadataRequest)
		res, err = queryClient.DenomMetadata(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Bank total
	case i.Ixplac.GetMsgType() == mbank.BankTotalMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(banktypes.QueryTotalSupplyRequest)
		res, err = queryClient.TotalSupply(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Bank total supply
	case i.Ixplac.GetMsgType() == mbank.BankTotalSupplyOfMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(banktypes.QuerySupplyOfRequest)
		res, err = queryClient.SupplyOf(
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

const (
	bankBalancesLabel      = "balances"
	bankDenomMetadataLabel = "denoms_metadata"
	bankSupplyLabel        = "supply"
)

func queryByLcdBank(i IXplaClient) (string, error) {
	url := util.MakeQueryLcdUrl(bankv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Bank balances
	case i.Ixplac.GetMsgType() == mbank.BankAllBalancesMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(*banktypes.QueryAllBalancesRequest)
		url = url + util.MakeQueryLabels(bankBalancesLabel, convertMsg.Address)

	// Bank balance
	case i.Ixplac.GetMsgType() == mbank.BankBalanceMsgType:
		// not supported now.
		convertMsg, _ := i.Ixplac.GetMsg().(*banktypes.QueryBalanceRequest)
		url = url + util.MakeQueryLabels(bankBalancesLabel, convertMsg.Address, convertMsg.Denom)

	// Bank denominations metadata
	case i.Ixplac.GetMsgType() == mbank.BankDenomsMetadataMsgType:
		url = url + bankDenomMetadataLabel

	// Bank denomination metadata
	case i.Ixplac.GetMsgType() == mbank.BankDenomMetadataMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(banktypes.QueryDenomMetadataRequest)
		url = url + util.MakeQueryLabels(bankDenomMetadataLabel, convertMsg.Denom)

	// Bank total
	case i.Ixplac.GetMsgType() == mbank.BankTotalMsgType:
		url = url + bankSupplyLabel

	// Bank total supply
	case i.Ixplac.GetMsgType() == mbank.BankTotalSupplyOfMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(banktypes.QuerySupplyOfRequest)
		url = url + util.MakeQueryLabels(bankSupplyLabel, convertMsg.Denom)

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err := CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil
}
