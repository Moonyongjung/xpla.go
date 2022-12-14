package bank

import (
	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// (Tx) make msg - bank send
func MakeBankSendMsg(bankSendMsg types.BankSendMsg, privKey key.PrivateKey) (banktypes.MsgSend, error) {
	msg, err := parseBankSendArgs(bankSendMsg, privKey)
	if err != nil {
		return banktypes.MsgSend{}, err
	}

	return msg, nil
}

// (Query) make msg - all balances
func MakeBankAllBalancesMsg(bankBalancesMsg types.BankBalancesMsg) (banktypes.QueryAllBalancesRequest, error) {
	if (types.BankBalancesMsg{}) == bankBalancesMsg {
		return banktypes.QueryAllBalancesRequest{}, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}

	msg, err := parseBankAllBalancesArgs(bankBalancesMsg)
	if err != nil {
		return banktypes.QueryAllBalancesRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - balance
func MakeBankBalanceMsg(bankBalancesMsg types.BankBalancesMsg) (banktypes.QueryBalanceRequest, error) {
	if (types.BankBalancesMsg{}) == bankBalancesMsg {
		return banktypes.QueryBalanceRequest{}, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}

	msg, err := parseBankBalanceArgs(bankBalancesMsg)
	if err != nil {
		return banktypes.QueryBalanceRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - denominations metadata
func MakeDenomsMetaDataMsg() (banktypes.QueryDenomsMetadataRequest, error) {
	return banktypes.QueryDenomsMetadataRequest{}, nil
}

// (Query) make msg - denomination metadata
func MakeDenomMetaDataMsg(denomMetadataMsg types.DenomMetadataMsg) (banktypes.QueryDenomMetadataRequest, error) {
	return banktypes.QueryDenomMetadataRequest{
		Denom: denomMetadataMsg.Denom,
	}, nil
}

// (Query) make msg - total supply
func MakeTotalSupplyMsg() (banktypes.QueryTotalSupplyRequest, error) {
	return banktypes.QueryTotalSupplyRequest{Pagination: core.PageRequest}, nil
}

// (Query) make msg - supply of
func MakeSupplyOfMsg(totalMsg types.TotalMsg) (banktypes.QuerySupplyOfRequest, error) {
	return banktypes.QuerySupplyOfRequest{Denom: totalMsg.Denom}, nil
}
