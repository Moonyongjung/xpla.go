package client

import (
	"context"

	mupgrade "github.com/Moonyongjung/xpla.go/core/upgrade"
	"github.com/Moonyongjung/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// Query client for upgrade module.
func queryUpgrade(xplac *XplaClient) (string, error) {
	queryClient := upgradetypes.NewQueryClient(xplac.Grpc)

	switch {
	// Upgrade applied
	case xplac.MsgType == mupgrade.UpgradeAppliedMsgType:
		convertMsg, _ := xplac.Msg.(upgradetypes.QueryAppliedPlanRequest)
		appliedPlanRes, err := queryClient.AppliedPlan(
			xplac.Context,
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

		if appliedPlanRes.Height == 0 {
			return "", err
		}
		headerData, err := appliedReturnBlockheader(appliedPlanRes, xplac.Opts.RpcURL, xplac.Context)
		if err != nil {
			return "", err
		}
		return string(headerData), nil

	// Upgrade all module versions
	case xplac.MsgType == mupgrade.UpgradeQueryAllModuleVersionsMsgType ||
		xplac.MsgType == mupgrade.UpgradeQueryModuleVersionsMsgType:
		convertMsg, _ := xplac.Msg.(upgradetypes.QueryModuleVersionsRequest)
		res, err = queryClient.ModuleVersions(
			xplac.Context,
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Upgrade plan
	case xplac.MsgType == mupgrade.UpgradePlanMsgType:
		convertMsg, _ := xplac.Msg.(upgradetypes.QueryCurrentPlanRequest)
		res, err = queryClient.CurrentPlan(
			xplac.Context,
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	default:
		return "", util.LogErr("invalid msg type")

	}

	out, err = printProto(xplac, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func appliedReturnBlockheader(res *upgradetypes.QueryAppliedPlanResponse, rpcUrl string, ctx context.Context) ([]byte, error) {
	if rpcUrl == "" {
		return nil, util.LogErr("need RPC URL")
	}
	clientCtx, err := util.NewClient()
	if err != nil {
		return nil, err
	}

	client, err := cmclient.NewClientFromNode(rpcUrl)
	if err != nil {
		return nil, err
	}
	clientCtx = clientCtx.WithClient(client)

	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	headers, err := node.BlockchainInfo(ctx, res.Height, res.Height)
	if err != nil {
		return nil, err
	}

	if len(headers.BlockMetas) == 0 {
		return nil, util.LogErr("no headers returns for height", res.Height)
	}

	bytes, err := clientCtx.LegacyAmino.MarshalJSONIndent(headers.BlockMetas[0], "", "  ")
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
