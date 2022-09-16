package evidence

import (
	"encoding/hex"

	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
)

// Parsing - evidence
func parseQueryEvidenceArgs(queryEvidenceMsg types.QueryEvidenceMsg) (*evidencetypes.QueryEvidenceRequest, error) {
	decodedHash, err := hex.DecodeString(queryEvidenceMsg.Hash)
	if err != nil {
		return nil, err
	}

	return &evidencetypes.QueryEvidenceRequest{
		EvidenceHash: decodedHash,
	}, nil
}

// Parsing - all evidences
func parseQueryAllEvidenceArgs() (*evidencetypes.QueryAllEvidenceRequest, error) {
	return &evidencetypes.QueryAllEvidenceRequest{
		Pagination: core.PageRequest,
	}, nil
}
