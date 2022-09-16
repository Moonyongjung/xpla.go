package client

import (
	mauth "github.com/Moonyongjung/xpla.go/core/auth"
	mauthz "github.com/Moonyongjung/xpla.go/core/authz"
	mbank "github.com/Moonyongjung/xpla.go/core/bank"
	mdist "github.com/Moonyongjung/xpla.go/core/distribution"
	mevidence "github.com/Moonyongjung/xpla.go/core/evidence"
	mevm "github.com/Moonyongjung/xpla.go/core/evm"
	mfeegrant "github.com/Moonyongjung/xpla.go/core/feegrant"
	mgov "github.com/Moonyongjung/xpla.go/core/gov"
	mmint "github.com/Moonyongjung/xpla.go/core/mint"
	mparams "github.com/Moonyongjung/xpla.go/core/params"
	mslashing "github.com/Moonyongjung/xpla.go/core/slashing"
	mstaking "github.com/Moonyongjung/xpla.go/core/staking"
	mupgrade "github.com/Moonyongjung/xpla.go/core/upgrade"
	mwasm "github.com/Moonyongjung/xpla.go/core/wasm"
	"github.com/Moonyongjung/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/gogo/protobuf/proto"
)

var out []byte
var res proto.Message
var err error

// Query transactions and xpla blockchain information.
// Execute a query of functions for all modules.
// After module query messages are generated, it receives query messages/information to the xpla client receiver and transmits a query message.
func (xplac *XplaClient) Query() (string, error) {
	if xplac.Err != nil {
		return "", xplac.Err
	}

	if xplac.Grpc == nil {
		return "", util.LogErr("error: Need GRPC URL when query methods")
	}

	if xplac.Module == mauth.AuthModule {
		return queryAuth(xplac)

	} else if xplac.Module == mauthz.AuthzModule {
		return queryAuthz(xplac)

	} else if xplac.Module == mbank.BankModule {
		return queryBank(xplac)

	} else if xplac.Module == mdist.DistributionModule {
		return queryDistribution(xplac)

	} else if xplac.Module == mevidence.EvidenceModule {
		return queryEvidence(xplac)

	} else if xplac.Module == mevm.EvmModule {
		return queryEvm(xplac)

	} else if xplac.Module == mfeegrant.FeegrantModule {
		return queryFeegrant(xplac)

	} else if xplac.Module == mgov.GovModule {
		return queryGov(xplac)

	} else if xplac.Module == mmint.MintModule {
		return queryMint(xplac)

	} else if xplac.Module == mslashing.SlashingModule {
		return querySlashing(xplac)

	} else if xplac.Module == mstaking.StakingModule {
		return queryStaking(xplac)

	} else if xplac.Module == mparams.ParamsModule {
		return queryParams(xplac)

	} else if xplac.Module == mupgrade.UpgradeModule {
		return queryUpgrade(xplac)

	} else if xplac.Module == mwasm.WasmModule {
		return queryWasm(xplac)
	} else {
		return "", util.LogErr("No module")
	}
}

// Print protobuf message by using cosmos sdk codec.
func printProto(xplac *XplaClient, toPrint proto.Message) ([]byte, error) {
	out, err := xplac.EncodingConfig.Marshaler.MarshalJSON(toPrint)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Print object by using cosmos sdk legacy amino.
func printObjectLegacy(xplac *XplaClient, toPrint interface{}) ([]byte, error) {
	out, err := xplac.EncodingConfig.Amino.MarshalJSON(toPrint)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// For auth module and gov module, make cosmos sdk client for querying.
func clientForQuery(xplac *XplaClient) (cmclient.Context, error) {
	client, err := cmclient.NewClientFromNode(xplac.Opts.RpcURL)
	if err != nil {
		return cmclient.Context{}, err
	}

	clientCtx, err := util.NewClient()
	if err != nil {
		return cmclient.Context{}, err
	}

	clientCtx = clientCtx.
		WithNodeURI(xplac.Opts.RpcURL).
		WithClient(client)

	return clientCtx, nil
}

// After call(as query) solidity contract, the response of chain is unpacked by ABI.
func getAbiUnpack(callName string, data []byte) ([]interface{}, error) {
	contractAbi, err := mevm.XplaSolContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	unpacked, err := contractAbi.Unpack(callName, data)
	if err != nil {
		return nil, err
	}

	return unpacked, nil
}
