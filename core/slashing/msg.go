package slashing

import (
	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/xpladev/xpla/app/params"
)

// (Tx) make msg - unjail
func MakeUnjailMsg(privKey key.PrivateKey) (slashingtypes.MsgUnjail, error) {
	return parseUnjailArgs(privKey)
}

// (Query) make msg - slahsing params
func MakeQuerySlashingParamsMsg() (slashingtypes.QueryParamsRequest, error) {
	return slashingtypes.QueryParamsRequest{}, nil
}

// (Query) make msg - signing infos
func MakeQuerySigningInfosMsg() (slashingtypes.QuerySigningInfosRequest, error) {
	return slashingtypes.QuerySigningInfosRequest{
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - signing info
func MakeQuerySigningInfoMsg(signingInfoMsg types.SigningInfoMsg, xplacEncodingConfig params.EncodingConfig) (slashingtypes.QuerySigningInfoRequest, error) {
	return parseQuerySigingInfoArgs(signingInfoMsg, xplacEncodingConfig)
}
