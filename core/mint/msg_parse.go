package mint

import (
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// Parsing - mint params
func parseQueryMintParamsArgs() (minttypes.QueryParamsRequest, error) {
	return minttypes.QueryParamsRequest{}, nil
}

// Parsing - inflation
func parseQueryInflationArgs() (minttypes.QueryInflationRequest, error) {
	return minttypes.QueryInflationRequest{}, nil
}

// Parsing - annual provisions
func parseQueryAnnualProvisionsArgs() (minttypes.QueryAnnualProvisionsRequest, error) {
	return minttypes.QueryAnnualProvisionsRequest{}, nil
}
