package module

import (
	mevm "github.com/Moonyongjung/xpla.go/core/evm"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/ethereum/go-ethereum"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Query client for evm module.
func (i IXplaClient) QueryEvm() (string, error) {
	evmClient, err := util.NewEvmClient(i.Ixplac.GetEvmRpc(), i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	gasAdj := i.Ixplac.GetGasAdjustment()
	if i.Ixplac.GetGasAdjustment() == "" {
		gasAdj = types.DefaultGasAdjustment
	}

	gasLimit := i.Ixplac.GetGasLimit()
	if i.Ixplac.GetGasLimit() == "" {
		gasLimitAdjustment, err := util.GasLimitAdjustment(util.FromStringToUint64(util.DefaultEvmGasLimit), gasAdj)
		if err != nil {
			return "", err
		}
		gasLimit = gasLimitAdjustment
	}

	gasPrice := i.Ixplac.GetGasPrice()
	if i.Ixplac.GetGasPrice() == "" {
		gasPrice = types.DefaultGasPrice
	}

	gasPriceBigInt, err := util.FromStringToBigInt(gasPrice)
	if err != nil {
		return "", err
	}

	switch {
	// Evm call contract
	case i.Ixplac.GetMsgType() == mevm.EvmCallSolContractMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(types.CallSolContractMsg)

		callByteData, err := GetAbiPack(convertMsg.ContractFuncCallName, convertMsg.Args...)
		if err != nil {
			return "", err
		}

		fromAddr := util.FromStringToByte20Address(i.Ixplac.GetPrivateKey().PubKey().Address().String())
		toAddr := util.FromStringToByte20Address(convertMsg.ContractAddress)
		value, err := util.FromStringToBigInt("0")
		if err != nil {
			return "", err
		}

		msg := ethereum.CallMsg{
			From:     fromAddr,
			To:       &toAddr,
			Gas:      util.FromStringToUint64(gasLimit),
			GasPrice: gasPriceBigInt,
			Value:    value,
			Data:     callByteData,
		}

		res, err := evmClient.Client.CallContract(evmClient.Ctx, msg, nil)
		if err != nil {
			return "", err
		}

		result, err := GetAbiUnpack(convertMsg.ContractFuncCallName, res)
		if err != nil {
			return "", err
		}

		var callSolContractResponse types.CallSolContractResponse
		for _, res := range result {
			callSolContractResponse.ContractResponse = append(callSolContractResponse.ContractResponse, util.ToString(res, ""))
		}

		return jsonReturn(callSolContractResponse)

	// Evm transaction by hash
	case i.Ixplac.GetMsgType() == mevm.EvmGetTransactionByHashMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(types.GetTransactionByHashMsg)
		commonTxHash := util.FromStringHexToHash(convertMsg.TxHash)
		tx, isPending, err := evmClient.Client.TransactionByHash(evmClient.Ctx, commonTxHash)
		if isPending {
			return "", util.LogErr("tx is pending..")
		}
		if err != nil {
			return "", err
		}

		json, err := tx.MarshalJSON()
		if err != nil {
			return "", err
		}

		return string(json), nil

	// Evm block by hash or height
	case i.Ixplac.GetMsgType() == mevm.EvmGetBlockByHashHeightMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(types.GetBlockByHashHeightMsg)
		var block *ethtypes.Block
		var blockResponse types.BlockResponse

		if convertMsg.BlockHash != "" {
			commonBlockHash := util.FromStringHexToHash(convertMsg.BlockHash)
			block, err = evmClient.Client.BlockByHash(evmClient.Ctx, commonBlockHash)
			if err != nil {
				return "", err
			}
		} else {
			blockNumber, err := util.FromStringToBigInt(convertMsg.BlockHeight)
			if err != nil {
				return "", err
			}

			block, err = evmClient.Client.BlockByNumber(evmClient.Ctx, blockNumber)
			if err != nil {
				return "", err
			}
		}

		txs := block.Body().Transactions
		uncles := block.Body().Uncles

		blockResponse.BlockHeader = block.Header()
		blockResponse.Transactions = txs
		blockResponse.Uncles = uncles

		return jsonReturn(blockResponse)

	// Evm account information
	case i.Ixplac.GetMsgType() == mevm.EvmQueryAccountInfoMsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(types.AccountInfoMsg)
		account := util.FromStringToByte20Address(convertMsg.Account)

		balance, err := evmClient.Client.BalanceAt(evmClient.Ctx, account, nil)
		if err != nil {
			return "", err
		}
		currentNonce, err := evmClient.Client.NonceAt(evmClient.Ctx, account, nil)
		if err != nil {
			return "", err
		}
		storage, err := evmClient.Client.StorageAt(evmClient.Ctx, account, util.FromStringHexToHash("0"), nil)
		if err != nil {
			return "", err
		}
		code, err := evmClient.Client.CodeAt(evmClient.Ctx, account, nil)
		if err != nil {
			return "", err
		}
		pendingBalance, err := evmClient.Client.PendingBalanceAt(evmClient.Ctx, account)
		if err != nil {
			return "", err
		}
		pendingNonce, err := evmClient.Client.PendingNonceAt(evmClient.Ctx, account)
		if err != nil {
			return "", err
		}
		pendingStorage, err := evmClient.Client.PendingStorageAt(evmClient.Ctx, account, util.FromStringHexToHash("0"))
		if err != nil {
			return "", err
		}
		pendingCode, err := evmClient.Client.PendingCodeAt(evmClient.Ctx, account)
		if err != nil {
			return "", err
		}
		pendingTransactionCount, err := evmClient.Client.PendingTransactionCount(evmClient.Ctx)
		if err != nil {
			return "", err
		}

		bech32Addr, err := util.FromByte20AddressToCosmosAddr(account)
		if err != nil {
			return "", err
		}

		var accountInfoResponse types.AccountInfoResponse

		accountInfoResponse.Account = account.Hex()
		accountInfoResponse.Bech32Account = bech32Addr.String()
		accountInfoResponse.Balance = balance
		accountInfoResponse.Nonce = currentNonce
		accountInfoResponse.Storage = string(storage)
		accountInfoResponse.Code = string(code)
		accountInfoResponse.PendingBalance = pendingBalance
		accountInfoResponse.PendingNonce = pendingNonce
		accountInfoResponse.PendingStorage = string(pendingStorage)
		accountInfoResponse.PendingCode = string(pendingCode)
		accountInfoResponse.PendingTransactionCount = pendingTransactionCount

		return jsonReturn(accountInfoResponse)

	// Evm suggest gas price
	case i.Ixplac.GetMsgType() == mevm.EvmSuggestGasPriceMsgType:
		gasPrice, err := evmClient.Client.SuggestGasPrice(evmClient.Ctx)
		if err != nil {
			return "", err
		}

		gasTipCap, err := evmClient.Client.SuggestGasTipCap(evmClient.Ctx)
		if err != nil {
			return "", err
		}

		var suggestGasPriceResponse types.SuggestGasPriceResponse
		suggestGasPriceResponse.GasPrice = gasPrice
		suggestGasPriceResponse.GasTipCap = gasTipCap

		return jsonReturn(suggestGasPriceResponse)

	// Evm chain ID
	case i.Ixplac.GetMsgType() == mevm.EvmQueryChainIdMsgType:
		chainId, err := evmClient.Client.ChainID(evmClient.Ctx)
		if err != nil {
			return "", err
		}

		var ethChainIdResponse types.EthChainIdResponse
		ethChainIdResponse.ChainID = chainId

		return jsonReturn(ethChainIdResponse)

	// Evm latest block height
	case i.Ixplac.GetMsgType() == mevm.EvmQueryCurrentBlockNumberMsgType:
		blockNumber, err := evmClient.Client.BlockNumber(evmClient.Ctx)
		if err != nil {
			return "", err
		}

		var ethBlockNumberResponse types.EthBlockNumberResponse
		ethBlockNumberResponse.BlockNumber = blockNumber

		return jsonReturn(ethBlockNumberResponse)

	// Web3 client version
	case i.Ixplac.GetMsgType() == mevm.EvmWeb3ClientVersionMsgType:
		var result string
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "web3_clientVersion")
		if err != nil {
			return "", err
		}

		var web3ClientVersionResponse types.Web3ClientVersionResponse
		web3ClientVersionResponse.Web3ClientVersion = result

		return jsonReturn(web3ClientVersionResponse)

	// Web3 sha
	case i.Ixplac.GetMsgType() == mevm.EvmWeb3Sha3MsgType:
		convertMsg, _ := i.Ixplac.GetMsg().(types.Web3Sha3Msg)

		var result string
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "web3_sha3", convertMsg.InputParam)
		if err != nil {
			return "", err
		}

		var web3Sha3Response types.Web3Sha3Response
		web3Sha3Response.Web3Sha3 = result

		return jsonReturn(web3Sha3Response)

	// network ID
	case i.Ixplac.GetMsgType() == mevm.EvmNetVersionMsgType:
		var result string
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "net_version")
		if err != nil {
			return "", err
		}

		var netVersionResponse types.NetVersionResponse
		netVersionResponse.NetVersion = result

		return jsonReturn(netVersionResponse)

	// the number of peers
	case i.Ixplac.GetMsgType() == mevm.EvmNetPeerCountMsgType:
		var result int
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "net_peerCount")
		if err != nil {
			return "", err
		}

		var netPeerCountResponse types.NetPeerCountResponse
		netPeerCountResponse.NetPeerCount = result

		return jsonReturn(netPeerCountResponse)

	// actively listening for network connections
	case i.Ixplac.GetMsgType() == mevm.EvmNetListeningMsgType:
		var result bool
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "net_listening")
		if err != nil {
			return "", err
		}

		var netListeningResponse types.NetListeningResponse
		netListeningResponse.NetListening = result

		return jsonReturn(netListeningResponse)

	default:
		return "", util.LogErr("invalid evm msg type")
	}
}

func jsonReturn(value interface{}) (string, error) {
	json, err := util.JsonMarshalData(value)
	if err != nil {
		return "", err
	}

	return string(json), nil
}
