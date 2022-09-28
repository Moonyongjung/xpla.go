package params

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/xpladev/xpla/app/params"
)

// (Tx) make msg - param change
func MakeProposalParamChangeMsg(paramChangeMsg types.ParamChangeMsg, privKey key.PrivateKey, encodingConfig params.EncodingConfig) (*govtypes.MsgSubmitProposal, error) {
	msg, err := parseProposalParamChangeArgs(paramChangeMsg, privKey, encodingConfig)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// (Query) make msg - subspace
func MakeQueryParamsSubspaceMsg(subspaceMsg types.SubspaceMsg) (proposal.QueryParamsRequest, error) {
	msg, err := parseQueryParamsSubspaceArgs(subspaceMsg)
	if err != nil {
		return proposal.QueryParamsRequest{}, err
	}

	return msg, nil
}
