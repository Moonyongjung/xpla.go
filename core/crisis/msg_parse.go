package crisis

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

// Parsing - invariant broken
func parseInvariantBrokenArgs(invariantBrokenMsg types.InvariantBrokenMsg, privKey key.PrivateKey) (crisistypes.MsgVerifyInvariant, error) {
	if invariantBrokenMsg.ModuleName == "" || invariantBrokenMsg.InvariantRoute == "" {
		return crisistypes.MsgVerifyInvariant{}, util.LogErr(errors.ErrInsufficientParams, "invalid module name or invariant route")
	}

	senderAddr, err := util.GetAddrByPrivKey(privKey)
	if err != nil {
		return crisistypes.MsgVerifyInvariant{}, err
	}
	msg := crisistypes.NewMsgVerifyInvariant(senderAddr, invariantBrokenMsg.ModuleName, invariantBrokenMsg.InvariantRoute)

	return *msg, nil
}
