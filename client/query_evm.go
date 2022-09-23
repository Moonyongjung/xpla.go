package client

import (
	mevm "github.com/Moonyongjung/xpla.go/core/evm"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/ethereum/go-ethereum"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Query client for evm module.
func queryEvm(xplac *XplaClient) (string, error) {
	evmClient, err := NewEvmClient(xplac.Opts.EvmRpcURL, xplac.Context)
	if err != nil {
		return "", err
	}

	if xplac.Opts.GasAdjustment == "" {
		xplac.WithGasAdjustment(defaultGasAdjustment)
	}

	if xplac.Opts.GasLimit == "" {
		gasLimitAdjustment, err := util.GasLimitAdjustment(util.FromStringToUint64(defaultEvmGasLimit), xplac.Opts.GasAdjustment)
		if err != nil {
			return "", err
		}
		xplac.WithGasLimit(gasLimitAdjustment)
	}

	if xplac.Opts.GasPrice == "" {
		xplac.WithGasPrice(defaultGasPrice)
	}

	gasPrice, err := util.FromStringToBigInt(xplac.Opts.GasPrice)
	if err != nil {
		return "", err
	}

	switch {
	// Evm call contract
	case xplac.MsgType == mevm.EvmCallSolContractMsgType:
		convertMsg, _ := xplac.Msg.(types.CallSolContractMsg)

		callByteData, err := getAbiPack(convertMsg.ContractFuncCallName, convertMsg.Args...)
		if err != nil {
			return "", err
		}

		fromAddr := util.FromStringToByte20Address(xplac.PrivateKey.PubKey().Address().String())
		toAddr := util.FromStringToByte20Address(convertMsg.ContractAddress)
		value, err := util.FromStringToBigInt("0")
		if err != nil {
			return "", err
		}

		msg := ethereum.CallMsg{
			From:     fromAddr,
			To:       &toAddr,
			Gas:      util.FromStringToUint64(xplac.Opts.GasLimit),
			GasPrice: gasPrice,
			Value:    value,
			Data:     callByteData,
		}

		res, err := evmClient.Client.CallContract(evmClient.Ctx, msg, nil)
		if err != nil {
			return "", err
		}

		result, err := getAbiUnpack(convertMsg.ContractFuncCallName, res)
		if err != nil {
			return "", err
		}

		var callSolContractResponse types.CallSolContractResponse
		for _, res := range result {
			callSolContractResponse.ContractResponse = append(callSolContractResponse.ContractResponse, util.ToString(res, ""))
		}

		return jsonReturn(callSolContractResponse)

	// Evm transaction by hash
	case xplac.MsgType == mevm.EvmGetTransactionByHashMsgType:
		convertMsg, _ := xplac.Msg.(types.GetTransactionByHashMsg)
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
	case xplac.MsgType == mevm.EvmGetBlockByHashHeightMsgType:
		convertMsg, _ := xplac.Msg.(types.GetBlockByHashHeightMsg)
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
	case xplac.MsgType == mevm.EvmQueryAccountInfoMsgType:
		convertMsg, _ := xplac.Msg.(types.AccountInfoMsg)
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
	case xplac.MsgType == mevm.EvmSuggestGasPriceMsgType:
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
	case xplac.MsgType == mevm.EvmQueryChainIdMsgType:
		chainId, err := evmClient.Client.ChainID(evmClient.Ctx)
		if err != nil {
			return "", err
		}

		var ethChainIdResponse types.EthChainIdResponse
		ethChainIdResponse.ChainID = chainId

		return jsonReturn(ethChainIdResponse)

	// Evm latest block height
	case xplac.MsgType == mevm.EvmQueryCurrentBlockNumberMsgType:
		blockNumber, err := evmClient.Client.BlockNumber(evmClient.Ctx)
		if err != nil {
			return "", err
		}

		var ethBlockNumberResponse types.EthBlockNumberResponse
		ethBlockNumberResponse.BlockNumber = blockNumber

		return jsonReturn(ethBlockNumberResponse)

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
