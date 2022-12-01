package auth

import (
	"fmt"
	"strings"

	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Parsing - auth params
func ParseAuthParamArgs() (authtypes.QueryParamsRequest, error) {
	return authtypes.QueryParamsRequest{}, nil
}

// Parsing - auth account
func parseQueryAccAddressArgs(queryAccAddresMsg types.QueryAccAddressMsg) (authtypes.QueryAccountRequest, error) {
	return authtypes.QueryAccountRequest{
		Address: queryAccAddresMsg.Address,
	}, nil
}

// Parsing - auth accounts
func parseQueryAccountsArgs() (authtypes.QueryAccountsRequest, error) {
	return authtypes.QueryAccountsRequest{
		Pagination: core.PageRequest,
	}, nil
}

// Parsing - transaction by evnets
func parseTxsByEventsArgs(txsByEventsMsg types.QueryTxsByEventsMsg) ([]string, error) {
	eventFormat := "{eventType}.{eventAttribute}={value}"
	eventsRaw := txsByEventsMsg.Events
	eventsStr := strings.Trim(eventsRaw, "'")

	if txsByEventsMsg.Page == "" {
		txsByEventsMsg.Page = util.FromIntToString(rest.DefaultPage)
	}
	if txsByEventsMsg.Limit == "" {
		txsByEventsMsg.Limit = util.FromIntToString(rest.DefaultLimit)
	}

	var events []string
	if strings.Contains(eventsStr, "&") {
		events = strings.Split(eventsStr, "&")
	} else {
		events = append(events, eventsStr)
	}

	var tmEvents []string

	for _, event := range events {
		if !strings.Contains(event, "=") {
			return []string{}, util.LogErr("invalid event; event", event, "should be of the format:", eventFormat)
		} else if strings.Count(event, "=") > 1 {
			return []string{}, util.LogErr("invalid event; event", event, "should be of the format:", eventFormat)
		}

		tokens := strings.Split(event, "=")
		if tokens[0] == tmtypes.TxHeightKey {
			event = fmt.Sprintf("%s=%s", tokens[0], tokens[1])
		} else {
			event = fmt.Sprintf("%s='%s'", tokens[0], tokens[1])
		}

		tmEvents = append(tmEvents, event)
	}

	tmEvents = append(tmEvents, txsByEventsMsg.Page)
	tmEvents = append(tmEvents, txsByEventsMsg.Limit)

	return tmEvents, nil
}

// Parsing - transaction
func parseQueryTxArgs(queryTxMsg types.QueryTxMsg) ([]string, error) {
	if queryTxMsg.Type == "" || queryTxMsg.Type == "hash" {
		if queryTxMsg.Value == "" {
			return []string{}, util.LogErr("argument should be a tx hash")
		}
		return []string{queryTxMsg.Value, "hash"}, nil

	} else if queryTxMsg.Type == "signature" {
		if queryTxMsg.Value == "" {
			return nil, fmt.Errorf("argument should be comma-separated signatures")
		}
		sigParts := strings.Split(queryTxMsg.Value, ",")

		tmEvents := make([]string, len(sigParts))
		for i, sig := range sigParts {
			tmEvents[i] = fmt.Sprintf("%s.%s='%s'", sdk.EventTypeTx, sdk.AttributeKeySignature, sig)
		}

		tmEvents = append(tmEvents, queryTxMsg.Type)
		return tmEvents, nil

	} else if queryTxMsg.Type == "acc_seq" {
		if queryTxMsg.Value == "" {
			return []string{}, util.LogErr("`acc_seq` type takes an argument '<addr>/<seq>'")
		}

		tmEvents := []string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeTx, sdk.AttributeKeyAccountSequence, queryTxMsg.Value),
		}
		tmEvents = append(tmEvents, queryTxMsg.Type)
		return tmEvents, nil

	} else {
		return []string{}, util.LogErr("Unknown type")
	}
}
