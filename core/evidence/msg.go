package evidence

import (
	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
)

// (Query) make msg - evidence
func MakeQueryEvidenceMsg(queryEvidenceMsg types.QueryEvidenceMsg) (evidencetypes.QueryEvidenceRequest, error) {
	return parseQueryEvidenceArgs(queryEvidenceMsg)
}

// (Query) make msg - all evidences
func MakeQueryAllEvidenceMsg() (evidencetypes.QueryAllEvidenceRequest, error) {
	return evidencetypes.QueryAllEvidenceRequest{
		Pagination: core.PageRequest,
	}, nil
}
