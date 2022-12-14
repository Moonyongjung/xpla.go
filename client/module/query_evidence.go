package module

import (
	mevidence "github.com/Moonyongjung/xpla.go/core/evidence"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	evidencev1beta1 "cosmossdk.io/api/cosmos/evidence/v1beta1"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
)

// Query client for evidence module.
func (i IXplaClient) QueryEvidence() (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcEvidence(i)
	} else {
		return queryByLcdEvidence(i)
	}
}

func queryByGrpcEvidence(i IXplaClient) (string, error) {
	queryClient := evidencetypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Query all evidences
	case i.Ixplac.GetMsgType() == mevidence.EvidenceQueryAllMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(evidencetypes.QueryAllEvidenceRequest)
		res, err = queryClient.AllEvidence(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Query evidence
	case i.Ixplac.GetMsgType() == mevidence.EvidenceQueryMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(evidencetypes.QueryEvidenceRequest)
		res, err = queryClient.Evidence(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err = printProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

const (
	evidenceEvidenceLabel = "evidence"
)

func queryByLcdEvidence(i IXplaClient) (string, error) {
	url := util.MakeQueryLcdUrl(evidencev1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Query all evidences
	case i.Ixplac.GetMsgType() == mevidence.EvidenceQueryAllMsgType:
		url = url + evidenceEvidenceLabel

	// Query evidence
	case i.Ixplac.GetMsgType() == mevidence.EvidenceQueryMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(evidencetypes.QueryEvidenceRequest)

		url = url + util.MakeQueryLabels(evidenceEvidenceLabel, convertMsg.EvidenceHash.String())

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil

}
