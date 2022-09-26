package module

import (
	mevidence "github.com/Moonyongjung/xpla.go/core/evidence"
	"github.com/Moonyongjung/xpla.go/util"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
)

// Query client for evidence module.
func (i IXplaClient) QueryEvidence() (string, error) {
	queryClient := evidencetypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Query all evidences
	case i.Ixplac.GetMsgType() == mevidence.EvidenceQueryAllMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(*evidencetypes.QueryAllEvidenceRequest)
		res, err = queryClient.AllEvidence(
			i.Ixplac.GetContext(),
			convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Query evidence
	case i.Ixplac.GetMsgType() == mevidence.EvidenceQueryMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(*evidencetypes.QueryEvidenceRequest)
		res, err = queryClient.Evidence(
			i.Ixplac.GetContext(),
			convertMsg,
		)
		if err != nil {
			return "", err
		}

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err = printProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
