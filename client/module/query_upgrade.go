package module

import (
	"context"

	mupgrade "github.com/Moonyongjung/xpla.go/core/upgrade"
	"github.com/Moonyongjung/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// Query client for upgrade module.
func (i IXplaClient) QueryUpgrade() (string, error) {
	queryClient := upgradetypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Upgrade applied
	case i.Ixplac.GetMsgType() == mupgrade.UpgradeAppliedMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(upgradetypes.QueryAppliedPlanRequest)
		appliedPlanRes, err := queryClient.AppliedPlan(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

		if appliedPlanRes.Height == 0 {
			return "", err
		}
		headerData, err := appliedReturnBlockheader(appliedPlanRes, i.Ixplac.GetRpc(), i.Ixplac.GetContext())
		if err != nil {
			return "", err
		}
		return string(headerData), nil

	// Upgrade all module versions
	case i.Ixplac.GetMsgType() == mupgrade.UpgradeQueryAllModuleVersionsMsgType ||
		i.Ixplac.GetMsgType() == mupgrade.UpgradeQueryModuleVersionsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(upgradetypes.QueryModuleVersionsRequest)
		res, err = queryClient.ModuleVersions(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	// Upgrade plan
	case i.Ixplac.GetMsgType() == mupgrade.UpgradePlanMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(upgradetypes.QueryCurrentPlanRequest)
		res, err = queryClient.CurrentPlan(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", err
		}

	default:
		return "", util.LogErr("invalid msg type")

	}

	out, err = printProto(i, res)
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
