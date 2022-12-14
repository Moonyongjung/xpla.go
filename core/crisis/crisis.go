package crisis

import (
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"

	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

// (Tx) make msg - invariant broken
func MakeInvariantRouteMsg(invariantBrokenMsg types.InvariantBrokenMsg, privKey key.PrivateKey) (crisistypes.MsgVerifyInvariant, error) {
	msg, err := parseInvariantBrokenArgs(invariantBrokenMsg, privKey)
	if err != nil {
		return crisistypes.MsgVerifyInvariant{}, err
	}

	return msg, nil
}
