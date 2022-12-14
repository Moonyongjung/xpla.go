package crisis

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

// Parsing - invariant broken
func parseInvariantBrokenArgs(invariantBrokenMsg types.InvariantBrokenMsg, privKey key.PrivateKey) (crisistypes.MsgVerifyInvariant, error) {
	if invariantBrokenMsg.ModuleName == "" || invariantBrokenMsg.InvariantRoute == "" {
		return crisistypes.MsgVerifyInvariant{}, util.LogErr("invalid module name or invariant route")
	}

	senderAddr := util.GetAddrByPrivKey(privKey)
	msg := crisistypes.NewMsgVerifyInvariant(senderAddr, invariantBrokenMsg.ModuleName, invariantBrokenMsg.InvariantRoute)

	return *msg, nil
}
