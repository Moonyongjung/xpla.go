package bank

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// (Tx) make msg - bank send
func MakeBankSendMsg(bankSendMsg types.BankSendMsg, privKey key.PrivateKey) (*bank.MsgSend, error) {
	msg, err := parseBankSendArgs(bankSendMsg, privKey)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// (Query) make msg - all balances
func MakeBankAllBalancesMsg(bankBalancesMsg types.BankBalancesMsg) (*bank.QueryAllBalancesRequest, error) {
	if (types.BankBalancesMsg{}) == bankBalancesMsg {
		return &bank.QueryAllBalancesRequest{}, util.LogErr("Empty request or type of parameter is not correct")
	}

	msg, err := parseBankAllBalancesArgs(bankBalancesMsg)
	if err != nil {
		return &bank.QueryAllBalancesRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - balance
func MakeBankBalanceMsg(bankBalancesMsg types.BankBalancesMsg) (*bank.QueryBalanceRequest, error) {
	if (types.BankBalancesMsg{}) == bankBalancesMsg {
		return &bank.QueryBalanceRequest{}, util.LogErr("Empty request or type of parameter is not correct")
	}

	msg, err := parseBankBalanceArgs(bankBalancesMsg)
	if err != nil {
		return &bank.QueryBalanceRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - denominations metadata
func MakeDenomsMetaDataMsg() (bank.QueryDenomsMetadataRequest, error) {
	msg := parseDenomsMetaDataArgs()
	return msg, nil
}

// (Query) make msg - denomination metadata
func MakeDenomMetaDataMsg(denomMetadataMsg types.DenomMetadataMsg) (bank.QueryDenomMetadataRequest, error) {
	msg := parseDenomMetaDataArgs(denomMetadataMsg)
	return msg, nil
}

// (Query) make msg - total supply
func MakeTotalSupplyMsg() (bank.QueryTotalSupplyRequest, error) {
	msg, err := parseTotalArgs()
	if err != nil {
		return bank.QueryTotalSupplyRequest{}, err
	}
	return msg, nil
}

// (Query) make msg - supply of
func MakeSupplyOfMsg(totalMsg types.TotalMsg) (bank.QuerySupplyOfRequest, error) {
	msg, err := parseSupplyOfArgs(totalMsg)
	if err != nil {
		return bank.QuerySupplyOfRequest{}, err
	}
	return msg, nil
}
