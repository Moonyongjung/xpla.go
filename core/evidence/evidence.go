package evidence

import (
	"github.com/Moonyongjung/xpla.go/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
)

// (Query) make msg - evidence
func MakeQueryEvidenceMsg(queryEvidenceMsg types.QueryEvidenceMsg) (*evidencetypes.QueryEvidenceRequest, error) {
	msg, err := parseQueryEvidenceArgs(queryEvidenceMsg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// (Query) make msg - all evidences
func MakeQueryAllEvidenceMsg() (*evidencetypes.QueryAllEvidenceRequest, error) {
	msg, err := parseQueryAllEvidenceArgs()
	if err != nil {
		return nil, err
	}

	return msg, nil
}
