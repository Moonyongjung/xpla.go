package client

import (
	"github.com/Moonyongjung/xpla.go/core"

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
	queryClient := core.NewIXplaClient(xplac, qt)

	switch {
	case xplac.Module == mauth.AuthModule:
		return mauth.QueryAuth(*queryClient)

	case xplac.Module == mauthz.AuthzModule:
		return mauthz.QueryAuthz(*queryClient)

	case xplac.Module == mbank.BankModule:
		return mbank.QueryBank(*queryClient)

	case xplac.Module == mbase.Base:
		return mbase.QueryBase(*queryClient)

	case xplac.Module == mdist.DistributionModule:
		return mdist.QueryDistribution(*queryClient)

	case xplac.Module == mevidence.EvidenceModule:
		return mevidence.QueryEvidence(*queryClient)

	case xplac.Module == mevm.EvmModule:
		return mevm.QueryEvm(*queryClient)

	case xplac.Module == mfeegrant.FeegrantModule:
		return mfeegrant.QueryFeegrant(*queryClient)

	case xplac.Module == mgov.GovModule:
		return mgov.QueryGov(*queryClient)

	case xplac.Module == mibc.IbcModule:
		return mibc.QueryIbc(*queryClient)

	case xplac.Module == mmint.MintModule:
		return mmint.QueryMint(*queryClient)

	case xplac.Module == mparams.ParamsModule:
		return mparams.QueryParams(*queryClient)

	case xplac.Module == mreward.RewardModule:
		return mreward.QueryReward(*queryClient)

	case xplac.Module == mslashing.SlashingModule:
		return mslashing.QuerySlashing(*queryClient)

	case xplac.Module == mstaking.StakingModule:
		return mstaking.QueryStaking(*queryClient)

	case xplac.Module == mupgrade.UpgradeModule:
		return mupgrade.QueryUpgrade(*queryClient)

	case xplac.Module == mwasm.WasmModule:
		return mwasm.QueryWasm(*queryClient)

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
