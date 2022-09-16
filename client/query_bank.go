package client

import (
	"context"

	mbank "github.com/Moonyongjung/xpla.go/core/bank"
	"github.com/Moonyongjung/xpla.go/util"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// Query client for bank module.
func queryBank(xplac *XplaClient) (string, error) {
	queryClient := banktypes.NewQueryClient(xplac.Grpc)

	switch {
	// Bank balances
	case xplac.MsgType == mbank.BankAllBalancesMsgType:
		convertMsg, _ := xplac.Msg.(*banktypes.QueryAllBalancesRequest)
		res, err = queryClient.AllBalances(
			context.Background(),
			convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Bank balance
	case xplac.MsgType == mbank.BankBalanceMsgType:
		convertMsg, _ := xplac.Msg.(*banktypes.QueryBalanceRequest)
		res, err = queryClient.Balance(
			context.Background(),
			convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Bank denominations metadata
	case xplac.MsgType == mbank.BankDenomsMetadataMsgType:
		convertMsg, _ := xplac.Msg.(banktypes.QueryDenomsMetadataRequest)
		res, err = queryClient.DenomsMetadata(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Bank denomination metadata
	case xplac.MsgType == mbank.BankDenomMetadataMsgType:
		convertMsg, _ := xplac.Msg.(banktypes.QueryDenomMetadataRequest)
		res, err = queryClient.DenomMetadata(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Bank total
	case xplac.MsgType == mbank.BankTotalMsgType:
		convertMsg, _ := xplac.Msg.(banktypes.QueryTotalSupplyRequest)
		res, err = queryClient.TotalSupply(
			context.Background(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Bank total supply
	case xplac.MsgType == mbank.BankTotalSupplyOfMsgType:
		convertMsg, _ := xplac.Msg.(banktypes.QuerySupplyOfRequest)
		res, err = queryClient.SupplyOf(
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
