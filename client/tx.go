package client

import (
	"encoding/base64"
	"math/big"

	"github.com/Moonyongjung/xpla.go/client/module"
	mevm "github.com/Moonyongjung/xpla.go/core/evm"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	kmultisig "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type deploySolTx struct {
	ChainId  *big.Int
	Nonce    *big.Int
	Value    *big.Int
	GasLimit uint64
	GasPrice *big.Int
}

// Create and sign a transaction before it is broadcasted to xpla chain.
// Options required for create and sign are stored in the xpla client and reflected when the values of those options exist.
// Create and sign transaction must be needed in order to send transaction to the chain.
func (xplac *XplaClient) CreateAndSignTx() ([]byte, error) {
	if xplac.Err != nil {
		return nil, xplac.Err
	}

	if xplac.Opts.AccountNumber == "" || xplac.Opts.Sequence == "" {
		if xplac.Opts.LcdURL == "" && xplac.Opts.GrpcURL == "" {
			xplac.WithAccountNumber(util.FromUint64ToString(types.DefaultAccNum))
			xplac.WithSequence(util.FromUint64ToString(types.DefaultAccSeq))
		} else {
			account, err := xplac.LoadAccount(sdk.AccAddress(xplac.Opts.PrivateKey.PubKey().Address()))
			if err != nil {
				return nil, err
			}
			xplac.WithAccountNumber(util.FromUint64ToString(account.GetAccountNumber()))
			xplac.WithSequence(util.FromUint64ToString(account.GetSequence()))
		}
	}

	if xplac.Opts.GasAdjustment == "" {
		xplac.WithGasAdjustment(types.DefaultGasAdjustment)
	}

	if xplac.Opts.GasPrice == "" {
		xplac.WithGasPrice(types.DefaultGasPrice)
	}

	if xplac.Module == mevm.EvmModule {
		return xplac.createAndSignEvmTx()

	} else {
		builder, err := setTxBuilderMsg(xplac)
		if err != nil {
			return nil, err
		}

		// Set default sign mode (DIRECT=1)
		if xplac.Opts.SignMode == signing.SignMode_SIGN_MODE_UNSPECIFIED {
			xplac.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
		}

		if xplac.Opts.GasLimit == "" {
			if xplac.Opts.LcdURL == "" && xplac.Opts.GrpcURL == "" {
				xplac.WithGasLimit(types.DefaultGasLimit)
			} else {
				simulate, err := xplac.Simulate(builder)
				if err != nil {
					return nil, err
				}
				gasLimitAdjustment, err := util.GasLimitAdjustment(simulate.GasInfo.GasUsed, xplac.Opts.GasAdjustment)
				if err != nil {
					return nil, err
				}
				xplac.WithGasLimit(gasLimitAdjustment)
			}
		}

		if xplac.Opts.FeeAmount == "" {
			gasLimitBigInt, err := util.FromStringToBigInt(xplac.Opts.GasLimit)
			if err != nil {
				return nil, err
			}

			gasPriceBigInt, err := util.FromStringToBigInt(xplac.Opts.GasPrice)
			if err != nil {
				return nil, err
			}

			feeAmount := util.MulBigInt(gasLimitBigInt, gasPriceBigInt)
			xplac.WithFeeAmount(util.FromBigIntToString(feeAmount))
		}

		builder = convertAndSetBuilder(xplac, builder)

		privs := []cryptotypes.PrivKey{xplac.Opts.PrivateKey}
		accNums := []uint64{util.FromStringToUint64(xplac.Opts.AccountNumber)}
		accSeqs := []uint64{util.FromStringToUint64(xplac.Opts.Sequence)}

		var sigsV2 []signing.SignatureV2

		err = txSignRound(xplac, sigsV2, privs, accSeqs, accNums, builder)
		if err != nil {
			return nil, err
		}

		sdkTx := builder.GetTx()
		txBytes, err := xplac.EncodingConfig.TxConfig.TxEncoder()(sdkTx)
		if err != nil {
			return nil, err
		}

		if xplac.Opts.OutputDocument != "" {
			jsonTx, err := xplac.EncodingConfig.TxConfig.TxJSONEncoder()(sdkTx)
			if err != nil {
				return nil, err
			}
			util.SaveJsonPretty(jsonTx, xplac.Opts.OutputDocument)

			return nil, nil
		}

		return txBytes, nil
	}

}

// Create transaction with unsigning.
// It returns txbytes of byte type when output document options is nil.
// If not, save the unsigned transaction file which name is "outputDocument"
func (xplac *XplaClient) CreateUnsignedTx() ([]byte, error) {
	if xplac.Err != nil {
		return nil, xplac.Err
	}
	builder, err := setTxBuilderMsg(xplac)
	if err != nil {
		return nil, err
	}

	sdkTx := builder.GetTx()
	txBytes, err := xplac.EncodingConfig.TxConfig.TxEncoder()(sdkTx)
	if err != nil {
		return nil, err
	}

	if xplac.Opts.OutputDocument != "" {
		jsonTx, err := xplac.EncodingConfig.TxConfig.TxJSONEncoder()(sdkTx)
		if err != nil {
			return nil, err
		}
		util.SaveJsonPretty(jsonTx, xplac.Opts.OutputDocument)

		return nil, nil
	}

	return txBytes, nil
}

// Sign created unsigned transaction.
func (xplac *XplaClient) SignTx(signTxMsg types.SignTxMsg) ([]byte, error) {
	if xplac.Err != nil {
		return nil, xplac.Err
	}
	var emptySignTxMsg types.SignTxMsg
	if signTxMsg == emptySignTxMsg {
		return nil, util.LogErr("need sign tx message of LCD client's option")
	}

	clientCtx, err := util.NewClient()
	if err != nil {
		return nil, err
	}
	err = clientCtx.Keyring.ImportPrivKey(types.XplaToolDefaultName, key.EncryptArmorPrivKey(xplac.Opts.PrivateKey, key.DefaultEncryptPassphrase), key.DefaultEncryptPassphrase)
	if err != nil {
		return nil, err
	}

	clientCtx.WithSignModeStr("direct")

	clientCtx, txFactory, newTx, err := readTxAndInitContexts(clientCtx, signTxMsg.UnsignedFileName)
	if err != nil {
		return nil, err
	}

	txCfg := clientCtx.TxConfig
	txBuilder, err := txCfg.WrapTxBuilder(newTx)
	if err != nil {
		return nil, err
	}

	signatureOnly := signTxMsg.SignatureOnly
	multisig := signTxMsg.MultisigAddress
	from := signTxMsg.FromAddress
	generateOnly := false
	offline := true

	_, fromName, _, err := cmclient.GetFromFields(txFactory.Keybase(), from, generateOnly)
	if err != nil {
		return nil, err
	}

	if multisig != "" {
		multisigAddr, err := sdk.AccAddressFromBech32(multisig)
		if err != nil {
			multisigAddr, _, _, err = cmclient.GetFromFields(txFactory.Keybase(), multisig, generateOnly)
			if err != nil {
				return nil, err
			}
		}
		err = authclient.SignTxWithSignerAddress(
			txFactory, clientCtx, multisigAddr, fromName, txBuilder, offline, signTxMsg.Overwrite,
		)
		if err != nil {
			return nil, err
		}
		signatureOnly = true
	} else {
		err = authclient.SignTx(txFactory, clientCtx, fromName, txBuilder, offline, signTxMsg.Overwrite)
	}
	if err != nil {
		return nil, err
	}

	var json []byte
	if signTxMsg.Amino {
		stdTx, err := tx.ConvertTxToStdTx(clientCtx.LegacyAmino, txBuilder.GetTx())
		if err != nil {
			return nil, err
		}
		req := authcli.BroadcastReq{
			Tx:   stdTx,
			Mode: "block|sync|async",
		}
		json, err = clientCtx.LegacyAmino.MarshalJSON(req)
		if err != nil {
			return nil, err
		}
	} else {
		json, err = marshalSignatureJSON(txCfg, txBuilder, signatureOnly)
		if err != nil {
			return nil, err
		}
	}

	if xplac.Opts.OutputDocument != "" {
		util.SaveJsonPretty(json, xplac.Opts.OutputDocument)
		return nil, nil
	}

	return json, nil
}

// Sign created unsigned transaction with multi signatures.
func (xplac *XplaClient) MultiSign(txMultiSignMsg types.TxMultiSignMsg) ([]byte, error) {

	clientCtx, err := util.NewClient()
	if err != nil {
		return nil, err
	}

	parseTx, err := authclient.ReadTxFromFile(clientCtx, txMultiSignMsg.FileName)
	if err != nil {
		return nil, err
	}

	txFactory := util.NewFactory(clientCtx)
	if txFactory.SignMode() == signing.SignMode_SIGN_MODE_UNSPECIFIED {
		txFactory = txFactory.WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON)
	}
	txFactory = txFactory.WithChainID(xplac.ChainId).
		WithAccountNumber(uint64(types.DefaultAccNum)).
		WithSequence(uint64(types.DefaultAccSeq))

	txCfg := clientCtx.TxConfig
	txBuilder, err := txCfg.WrapTxBuilder(parseTx)
	if err != nil {
		return nil, err
	}

	multisigInfo, err := getMultisigInfo(clientCtx, txMultiSignMsg.FromName)
	if err != nil {
		return nil, err
	}

	multisigPub := multisigInfo.GetPubKey().(*kmultisig.LegacyAminoPubKey)
	multisigSig := multisig.NewMultisig(len(multisigPub.PubKeys))
	clientCtx = clientCtx.WithOffline(txMultiSignMsg.Offline)
	if !clientCtx.Offline {
		accnum, seq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, multisigInfo.GetAddress())
		if err != nil {
			return nil, err
		}
		txFactory = txFactory.WithAccountNumber(accnum).WithSequence(seq)
	}

	for _, sigFile := range txMultiSignMsg.SignatureFiles {
		sigs, err := unmarshalSignatureJSON(clientCtx, sigFile)
		if err != nil {
			return nil, err
		}

		if txFactory.ChainID() == "" {
			return nil, util.LogErr("set the chain id with either the --chain-id flag or config file")
		}

		signingData := authsigning.SignerData{
			ChainID:       txFactory.ChainID(),
			AccountNumber: txFactory.AccountNumber(),
			Sequence:      txFactory.Sequence(),
		}

		for _, sig := range sigs {
			err = authsigning.VerifySignature(sig.PubKey, signingData, sig.Data, txCfg.SignModeHandler(), txBuilder.GetTx())
			if err != nil {
				addr, _ := sdk.AccAddressFromHex(sig.PubKey.Address().String())
				return nil, util.LogErr("couldn't verify signature for address", addr)
			}

			if err := multisig.AddSignatureV2(multisigSig, sig, multisigPub.GetPubKeys()); err != nil {
				return nil, err
			}
		}
	}

	sigV2 := signing.SignatureV2{
		PubKey:   multisigPub,
		Data:     multisigSig,
		Sequence: txFactory.Sequence(),
	}

	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		return nil, err
	}

	sigOnly := txMultiSignMsg.SignatureOnly
	aminoJson := txMultiSignMsg.Amino

	var json []byte
	if aminoJson {
		stdTx, err := tx.ConvertTxToStdTx(clientCtx.LegacyAmino, txBuilder.GetTx())
		if err != nil {
			return nil, err
		}

		req := authcli.BroadcastReq{
			Tx:   stdTx,
			Mode: "block|sync|async",
		}

		json, _ = clientCtx.LegacyAmino.MarshalJSON(req)
	} else {
		json, err = marshalSignatureJSON(txCfg, txBuilder, sigOnly)
		if err != nil {
			return nil, err
		}
	}

	if txMultiSignMsg.OutputDocument == "" {
		return json, nil
	}

	util.SaveJsonPretty(json, xplac.Opts.OutputDocument)
	return nil, nil
}

// Create and sign transaction of evm.
func (xplac *XplaClient) createAndSignEvmTx() ([]byte, error) {
	ethPrivKey, err := toECDSA(xplac.Opts.PrivateKey)
	if err != nil {
		return nil, err
	}

	chainId, err := util.ConvertEvmChainId(xplac.ChainId)
	if err != nil {
		return nil, err
	}

	if xplac.Opts.OutputDocument != "" {
		util.LogInfo("no create output document as tx of evm")
	}

	if xplac.Opts.GasLimit == "" {
		gasLimitAdjustment, err := util.GasLimitAdjustment(util.FromStringToUint64(util.DefaultEvmGasLimit), xplac.Opts.GasAdjustment)
		if err != nil {
			return nil, err
		}
		xplac.WithGasLimit(gasLimitAdjustment)
	}

	gasPrice, err := util.FromStringToBigInt(xplac.Opts.GasPrice)
	if err != nil {
		return nil, err
	}

	switch {
	case xplac.MsgType == mevm.EvmSendCoinMsgType:
		convertMsg, _ := xplac.Msg.(types.SendCoinMsg)
		toAddr := util.FromStringToByte20Address(convertMsg.ToAddress)
		amount, err := util.FromStringToBigInt(convertMsg.Amount)
		if err != nil {
			return nil, err
		}

		return evmTxSignRound(xplac, toAddr, gasPrice, amount, nil, chainId, ethPrivKey)

	case xplac.MsgType == mevm.EvmDeploySolContractMsgType:
		nonce, err := util.FromStringToBigInt(xplac.Opts.Sequence)
		if err != nil {
			return nil, err
		}

		value, err := util.FromStringToBigInt(util.DefaultSolidityValue)
		if err != nil {
			return nil, err
		}

		tx := deploySolTx{
			ChainId:  chainId,
			Nonce:    nonce,
			Value:    value,
			GasLimit: util.FromStringToUint64(xplac.Opts.GasLimit),
			GasPrice: gasPrice,
		}

		txbytes, err := util.JsonMarshalData(tx)
		if err != nil {
			return nil, err
		}

		return txbytes, nil

	case xplac.MsgType == mevm.EvmInvokeSolContractMsgType:
		convertMsg, _ := xplac.Msg.(types.InvokeSolContractMsg)
		invokeByteData, err := module.GetAbiPack(convertMsg.ContractFuncCallName, convertMsg.Args...)
		if err != nil {
			return nil, err
		}

		toAddr := util.FromStringToByte20Address(convertMsg.ContractAddress)
		amount, err := util.FromStringToBigInt("0")
		if err != nil {
			return nil, err
		}

		return evmTxSignRound(xplac, toAddr, gasPrice, amount, invokeByteData, chainId, ethPrivKey)

	default:
		return nil, util.LogErr("invalid evm msg type")
	}
}

// Encode transaction by using base64.
func (xplac *XplaClient) EncodeTx(encodeTxMsg types.EncodeTxMsg) (string, error) {
	if xplac.Err != nil {
		return "", xplac.Err
	}
	clientCtx, err := util.NewClient()
	if err != nil {
		return "", err
	}

	tx, err := authclient.ReadTxFromFile(clientCtx, encodeTxMsg.FileName)
	if err != nil {
		return "", err
	}

	txbytes, err := xplac.EncodingConfig.TxConfig.TxEncoder()(tx)
	if err != nil {
		return "", err
	}

	txbytesBase64 := base64.StdEncoding.EncodeToString(txbytes)

	return txbytesBase64, nil
}

// Decode transaction which encoded by base64
func (xplac *XplaClient) DecodeTx(decodeTxMsg types.DecodeTxMsg) (string, error) {
	if xplac.Err != nil {
		return "", xplac.Err
	}
	txbytes, err := base64.StdEncoding.DecodeString(decodeTxMsg.EncodedByteString)
	if err != nil {
		return "", err
	}

	tx, err := xplac.EncodingConfig.TxConfig.TxDecoder()(txbytes)
	if err != nil {
		return "", err
	}

	json, err := xplac.EncodingConfig.TxConfig.TxJSONEncoder()(tx)
	if err != nil {
		return "", err
	}

	return string(json), nil
}

// Validate signature
func (xplac *XplaClient) ValidateSignatures(validateSignaturesMsg types.ValidateSignaturesMsg) (string, error) {
	if xplac.Err != nil {
		return "", xplac.Err
	}
	resBool := true
	clientCtx, err := util.NewClient()
	if err != nil {
		return "", err
	}
	stdTx, err := authclient.ReadTxFromFile(clientCtx, validateSignaturesMsg.FileName)
	if err != nil {
		return "", err
	}

	sigTx := stdTx.(authsigning.SigVerifiableTx)
	signModeHandler := clientCtx.TxConfig.SignModeHandler()

	signers := sigTx.GetSigners()

	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return "", err
	}

	if len(sigs) != len(signers) {
		resBool = false
	}

	for i, sig := range sigs {
		var (
			PubKey         = sig.PubKey
			multisigHeader string
			multiSigMsg    string
			sigAddr        = sdk.AccAddress(PubKey.Address())
			sigSanity      = "OK"
		)

		if i >= len(signers) || !sigAddr.Equals(signers[i]) {
			sigSanity = "ERROR: signature does not match its respective signer"
			resBool = false
		}

		if !validateSignaturesMsg.Offline && resBool {
			accNum, accSeq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, sigAddr)
			if err != nil {
				return "", err
			}

			signingData := authsigning.SignerData{
				ChainID:       validateSignaturesMsg.ChainID,
				AccountNumber: accNum,
				Sequence:      accSeq,
			}
			err = authsigning.VerifySignature(PubKey, signingData, sig.Data, signModeHandler, sigTx)
			if err != nil {
				return "", err
			}
		}

		util.LogInfo(i, ":", sigAddr.String(), "[", sigSanity, "]", multisigHeader, multiSigMsg)
	}

	if resBool {
		return "success validate", nil
	} else {
		return "signature validation failed", nil
	}
}
