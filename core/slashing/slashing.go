package slashing

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// (Tx) make msg - unjail
func MakeUnjailMsg(privKey key.PrivateKey) (*slashingtypes.MsgUnjail, error) {
	msg, err := parseUnjailArgs(privKey)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// (Query) make msg - slahsing params
func MakeQuerySlashingParamsMsg() (slashingtypes.QueryParamsRequest, error) {
	msg, err := parseQuerySlashingParamsArgs()
	if err != nil {
		return slashingtypes.QueryParamsRequest{}, nil
	}

	return msg, nil
}

// (Query) make msg - signing infos
func MakeQuerySigningInfosMsg() (slashingtypes.QuerySigningInfosRequest, error) {
	msg, err := parseQuerySigingInfosArgs()
	if err != nil {
		return slashingtypes.QuerySigningInfosRequest{}, nil
	}

	return msg, nil
}

// (Query) make msg - signing info
func MakeQuerySigningInfoMsg(signingInfoMsg types.SigningInfoMsg, xplacEncodingConfig params.EncodingConfig) (slashingtypes.QuerySigningInfoRequest, error) {
	msg, err := parseQuerySigingInfoArgs(signingInfoMsg, xplacEncodingConfig)
	if err != nil {
		return slashingtypes.QuerySigningInfoRequest{}, nil
	}

	return msg, nil
}
