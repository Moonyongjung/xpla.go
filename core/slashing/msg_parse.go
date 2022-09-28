package slashing

import (
	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/xpladev/xpla/app/params"
)

// Parsing - unjail
func parseUnjailArgs(privKey key.PrivateKey) (*slashingtypes.MsgUnjail, error) {
	addr := util.GetAddrByPrivKey(privKey)
	msg := slashingtypes.NewMsgUnjail(sdk.ValAddress(addr))

	return msg, nil
}

// Parsing - slashing params
func parseQuerySlashingParamsArgs() (slashingtypes.QueryParamsRequest, error) {
	return slashingtypes.QueryParamsRequest{}, nil
}

// Parsing - signing infos
func parseQuerySigingInfosArgs() (slashingtypes.QuerySigningInfosRequest, error) {
	return slashingtypes.QuerySigningInfosRequest{
		Pagination: core.PageRequest,
	}, nil
}

// Parsing - signing info
func parseQuerySigingInfoArgs(signingInfoMsg types.SigningInfoMsg, xplacEncodingConfig params.EncodingConfig) (slashingtypes.QuerySigningInfoRequest, error) {
	var pk cryptotypes.PubKey
	err := xplacEncodingConfig.Marshaler.UnmarshalInterfaceJSON([]byte(signingInfoMsg.ConsPubKey), &pk)
	if err != nil {
		return slashingtypes.QuerySigningInfoRequest{}, err
	}

	return slashingtypes.QuerySigningInfoRequest{
		ConsAddress: sdk.ConsAddress(pk.Address()).String(),
	}, nil
}
