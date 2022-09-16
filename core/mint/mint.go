package mint

import (
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// (Query) make msg - mint params
func MakeQueryMintParamsMsg() (minttypes.QueryParamsRequest, error) {
	msg, err := parseQueryMintParamsArgs()
	if err != nil {
		return minttypes.QueryParamsRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - inflation
func MakeQueryInflationMsg() (minttypes.QueryInflationRequest, error) {
	msg, err := parseQueryInflationArgs()
	if err != nil {
		return minttypes.QueryInflationRequest{}, err
	}

	return msg, nil
}

// (Query) make msg - annual provisions
func MakeQueryAnnualProvisionsMsg() (minttypes.QueryAnnualProvisionsRequest, error) {
	msg, err := parseQueryAnnualProvisionsArgs()
	if err != nil {
		return minttypes.QueryAnnualProvisionsRequest{}, err
	}

	return msg, nil
}
