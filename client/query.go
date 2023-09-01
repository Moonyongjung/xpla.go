package client

import (
	"github.com/Moonyongjung/xpla.go/client/queries"
	qauth "github.com/Moonyongjung/xpla.go/client/queries/auth"
	qauthz "github.com/Moonyongjung/xpla.go/client/queries/authz"
	qbank "github.com/Moonyongjung/xpla.go/client/queries/bank"
	qbase "github.com/Moonyongjung/xpla.go/client/queries/base"
	qdist "github.com/Moonyongjung/xpla.go/client/queries/distribution"
	qevidence "github.com/Moonyongjung/xpla.go/client/queries/evidence"
	qevm "github.com/Moonyongjung/xpla.go/client/queries/evm"
	qfeegrant "github.com/Moonyongjung/xpla.go/client/queries/feegrant"
	qgov "github.com/Moonyongjung/xpla.go/client/queries/gov"
	qibc "github.com/Moonyongjung/xpla.go/client/queries/ibc"
	qmint "github.com/Moonyongjung/xpla.go/client/queries/mint"
	qparams "github.com/Moonyongjung/xpla.go/client/queries/params"
	qreward "github.com/Moonyongjung/xpla.go/client/queries/reward"
	qslashing "github.com/Moonyongjung/xpla.go/client/queries/slashing"
	qstaking "github.com/Moonyongjung/xpla.go/client/queries/staking"
	qupgrade "github.com/Moonyongjung/xpla.go/client/queries/upgrade"
	qwasm "github.com/Moonyongjung/xpla.go/client/queries/wasm"

	mauth "github.com/Moonyongjung/xpla.go/core/auth"
	mauthz "github.com/Moonyongjung/xpla.go/core/authz"
	mbank "github.com/Moonyongjung/xpla.go/core/bank"
	mbase "github.com/Moonyongjung/xpla.go/core/base"
	mdist "github.com/Moonyongjung/xpla.go/core/distribution"
	mevidence "github.com/Moonyongjung/xpla.go/core/evidence"
	mevm "github.com/Moonyongjung/xpla.go/core/evm"
	mfeegrant "github.com/Moonyongjung/xpla.go/core/feegrant"
	mgov "github.com/Moonyongjung/xpla.go/core/gov"
	mibc "github.com/Moonyongjung/xpla.go/core/ibc"
	mmint "github.com/Moonyongjung/xpla.go/core/mint"
	mparams "github.com/Moonyongjung/xpla.go/core/params"
	mreward "github.com/Moonyongjung/xpla.go/core/reward"
	mslashing "github.com/Moonyongjung/xpla.go/core/slashing"
	mstaking "github.com/Moonyongjung/xpla.go/core/staking"
	mupgrade "github.com/Moonyongjung/xpla.go/core/upgrade"
	mwasm "github.com/Moonyongjung/xpla.go/core/wasm"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"
)

// Query transactions and xpla blockchain information.
// Execute a query of functions for all modules.
// After module query messages are generated, it receives query messages/information to the xpla client receiver and transmits a query message.
func (xplac *XplaClient) Query() (string, error) {
	if xplac.Err != nil {
		return "", xplac.Err
	}

	if xplac.GetGrpcUrl() == "" && xplac.GetLcdURL() == "" {
		if xplac.Module == mevm.EvmModule {
			if xplac.GetEvmRpc() == "" {
				return "", util.LogErr(errors.ErrNotSatisfiedOptions, "evm JSON-RPC URL must exist")
			}

		} else {
			return "", util.LogErr(errors.ErrNotSatisfiedOptions, "at least one of the gRPC URL or LCD URL must exist for query")
		}
	}

	qt := setQueryType(xplac)
	ixplaClient := queries.NewIXplaClient(xplac, qt)

	switch {
	case xplac.Module == mauth.AuthModule:
		return qauth.QueryAuth(*ixplaClient)

	case xplac.Module == mauthz.AuthzModule:
		return qauthz.QueryAuthz(*ixplaClient)

	case xplac.Module == mbank.BankModule:
		return qbank.QueryBank(*ixplaClient)

	case xplac.Module == mbase.Base:
		return qbase.QueryBase(*ixplaClient)

	case xplac.Module == mdist.DistributionModule:
		return qdist.QueryDistribution(*ixplaClient)

	case xplac.Module == mevidence.EvidenceModule:
		return qevidence.QueryEvidence(*ixplaClient)

	case xplac.Module == mevm.EvmModule:
		return qevm.QueryEvm(*ixplaClient)

	case xplac.Module == mfeegrant.FeegrantModule:
		return qfeegrant.QueryFeegrant(*ixplaClient)

	case xplac.Module == mgov.GovModule:
		return qgov.QueryGov(*ixplaClient)

	case xplac.Module == mibc.IbcModule:
		return qibc.QueryIbc(*ixplaClient)

	case xplac.Module == mmint.MintModule:
		return qmint.QueryMint(*ixplaClient)

	case xplac.Module == mparams.ParamsModule:
		return qparams.QueryParams(*ixplaClient)

	case xplac.Module == mreward.RewardModule:
		return qreward.QueryReward(*ixplaClient)

	case xplac.Module == mslashing.SlashingModule:
		return qslashing.QuerySlashing(*ixplaClient)

	case xplac.Module == mstaking.StakingModule:
		return qstaking.QueryStaking(*ixplaClient)

	case xplac.Module == mupgrade.UpgradeModule:
		return qupgrade.QueryUpgrade(*ixplaClient)

	case xplac.Module == mwasm.WasmModule:
		return qwasm.QueryWasm(*ixplaClient)

	default:
		return "", util.LogErr(errors.ErrInvalidRequest, "invalid module")
	}
}

func setQueryType(xplac *XplaClient) uint8 {
	// Default query type is gRPC, not LCD.
	if xplac.Opts.GrpcURL != "" {
		return types.QueryGrpc
	} else {
		return types.QueryLcd
	}
}
