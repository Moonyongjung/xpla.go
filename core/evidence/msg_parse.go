package evidence

import (
	"encoding/hex"

	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
)

// Parsing - evidence
func parseQueryEvidenceArgs(queryEvidenceMsg types.QueryEvidenceMsg) (evidencetypes.QueryEvidenceRequest, error) {
	decodedHash, err := hex.DecodeString(queryEvidenceMsg.Hash)
	if err != nil {
		return evidencetypes.QueryEvidenceRequest{}, util.LogErr(errors.ErrParse, err)
	}

	return evidencetypes.QueryEvidenceRequest{
		EvidenceHash: decodedHash,
	}, nil
}
