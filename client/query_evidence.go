package client

import (
	mevidence "github.com/Moonyongjung/xpla.go/core/evidence"
	"github.com/Moonyongjung/xpla.go/util"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
)

// Query client for evidence module.
func queryEvidence(xplac *XplaClient) (string, error) {
	queryClient := evidencetypes.NewQueryClient(xplac.Grpc)

	switch {
	// Query all evidences
	case xplac.MsgType == mevidence.EvidenceQueryAllMsgType:
		convertMsg, _ := xplac.Msg.(*evidencetypes.QueryAllEvidenceRequest)
		res, err = queryClient.AllEvidence(
			xplac.Context,
			convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Query evidence
	case xplac.MsgType == mevidence.EvidenceQueryMsgType:
		convertMsg, _ := xplac.Msg.(*evidencetypes.QueryEvidenceRequest)
		res, err = queryClient.Evidence(
			xplac.Context,
			convertMsg,
		)
		if err != nil {
			return "", err
		}

	default:
		return "", util.LogErr("invalid msg type")
	}

	out, err = printProto(xplac, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
