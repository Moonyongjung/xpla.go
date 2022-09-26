package module

import (
	mbank "github.com/Moonyongjung/xpla.go/core/bank"
	"github.com/Moonyongjung/xpla.go/util"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// Query client for bank module.
func (i IXplaClient) QueryBank() (string, error) {
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
