package distribution

import (
	"github.com/Moonyongjung/xpla.go/core"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"
	"github.com/gogo/protobuf/proto"

	distv1beta1 "cosmossdk.io/api/cosmos/distribution/v1beta1"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

var out []byte
var res proto.Message
var err error

// Query client for distribution module.
func QueryDistribution(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcDist(i)
	} else {
		return queryByLcdDist(i)
	}
}

func queryByGrpcDist(i core.QueryClient) (string, error) {
	queryClient := disttypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Distribution params
	case i.Ixplac.GetMsgType() == DistributionQueryDistributionParamsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Distribution validator outstanding rewards
	case i.Ixplac.GetMsgType() == DistributionValidatorOutstandingRewardsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryValidatorOutstandingRewardsRequest)
		res, err = queryClient.ValidatorOutstandingRewards(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Distribution commission
	case i.Ixplac.GetMsgType() == DistributionQueryDistCommissionMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryValidatorCommissionRequest)
		res, err = queryClient.ValidatorCommission(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Distribution slashes
	case i.Ixplac.GetMsgType() == DistributionQuerySlashesMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryValidatorSlashesRequest)
		res, err = queryClient.ValidatorSlashes(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Distribution rewards
	case i.Ixplac.GetMsgType() == DistributionQueryRewardsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryDelegationRewardsRequest)
		res, err = queryClient.DelegationRewards(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Distribution total rewards
	case i.Ixplac.GetMsgType() == DistributionQueryTotalRewardsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryDelegationTotalRewardsRequest)
		res, err = queryClient.DelegationTotalRewards(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Distribution community pool
	case i.Ixplac.GetMsgType() == DistributionQueryCommunityPoolMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryCommunityPoolRequest)
		res, err = queryClient.CommunityPool(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err = core.PrintProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

const (
	distParamsLabel             = "params"
	distValidatorLabel          = "validators"
	distDelegatorLabel          = "delegators"
	distOutstandingRewardsLabel = "outstanding_rewards"
	distCommissionLabel         = "commission"
	distSlashesLabel            = "slashes"
	distRewardsLabel            = "rewards"
	distCommunityPoolLabel      = "community_pool"
)

func queryByLcdDist(i core.QueryClient) (string, error) {
	url := util.MakeQueryLcdUrl(distv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Distribution params
	case i.Ixplac.GetMsgType() == DistributionQueryDistributionParamsMsgType:
		url = url + distParamsLabel

	// Distribution validator outstanding rewards
	case i.Ixplac.GetMsgType() == DistributionValidatorOutstandingRewardsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryValidatorOutstandingRewardsRequest)

		url = url + util.MakeQueryLabels(distValidatorLabel, convertMsg.ValidatorAddress, distOutstandingRewardsLabel)

	// Distribution commission
	case i.Ixplac.GetMsgType() == DistributionQueryDistCommissionMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryValidatorCommissionRequest)

		url = url + util.MakeQueryLabels(distValidatorLabel, convertMsg.ValidatorAddress, distCommissionLabel)

	// Distribution slashes
	case i.Ixplac.GetMsgType() == DistributionQuerySlashesMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryValidatorSlashesRequest)

		url = url + util.MakeQueryLabels(distValidatorLabel, convertMsg.ValidatorAddress, distSlashesLabel)

	// Distribution rewards
	case i.Ixplac.GetMsgType() == DistributionQueryRewardsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryDelegationRewardsRequest)

		url = url + util.MakeQueryLabels(distDelegatorLabel, convertMsg.DelegatorAddress, distRewardsLabel, convertMsg.ValidatorAddress)

	// Distribution total rewards
	case i.Ixplac.GetMsgType() == DistributionQueryTotalRewardsMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(disttypes.QueryDelegationTotalRewardsRequest)

		url = url + util.MakeQueryLabels(distDelegatorLabel, convertMsg.DelegatorAddress, distRewardsLabel)

	// Distribution community pool
	case i.Ixplac.GetMsgType() == DistributionQueryCommunityPoolMsgType:
		url = url + distCommunityPoolLabel

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	return string(out), nil
}
