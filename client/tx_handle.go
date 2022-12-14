package client

import (
	"crypto/ecdsa"
	"math/big"
	"os"

	mauthz "github.com/Moonyongjung/xpla.go/core/authz"
	mbank "github.com/Moonyongjung/xpla.go/core/bank"
	mcrisis "github.com/Moonyongjung/xpla.go/core/crisis"
	mdist "github.com/Moonyongjung/xpla.go/core/distribution"
	mfeegrant "github.com/Moonyongjung/xpla.go/core/feegrant"
	mgov "github.com/Moonyongjung/xpla.go/core/gov"
	mparams "github.com/Moonyongjung/xpla.go/core/params"
	mreward "github.com/Moonyongjung/xpla.go/core/reward"
	mslashing "github.com/Moonyongjung/xpla.go/core/slashing"
	mstaking "github.com/Moonyongjung/xpla.go/core/staking"
	mupgrade "github.com/Moonyongjung/xpla.go/core/upgrade"
	mwasm "github.com/Moonyongjung/xpla.go/core/wasm"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/Moonyongjung/xpla.go/types"
	"github.com/Moonyongjung/xpla.go/types/errors"
	"github.com/Moonyongjung/xpla.go/util"

	"github.com/CosmWasm/wasmd/x/wasm"
	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	rewardtypes "github.com/xpladev/xpla/x/reward/types"
)

// Set message for transaction builder.
// Interface type messages are converted to correct type.
func setTxBuilderMsg(xplac *XplaClient) (cmclient.TxBuilder, error) {
	if xplac.Err != nil {
		return nil, xplac.Err
	}

	builder := xplac.EncodingConfig.TxConfig.NewTxBuilder()

	// Authz module
	if xplac.MsgType == mauthz.AuthzGrantMsgType {
		convertMsg, _ := xplac.Msg.(authz.MsgGrant)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mauthz.AuthzRevokeMsgType {
		convertMsg, _ := xplac.Msg.(authz.MsgRevoke)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mauthz.AuthzExecMsgType {
		convertMsg, _ := xplac.Msg.(authz.MsgExec)
		builder.SetMsgs(&convertMsg)

		// Bank module
	} else if xplac.MsgType == mbank.BankSendMsgType {
		convertMsg, _ := xplac.Msg.(banktypes.MsgSend)
		builder.SetMsgs(&convertMsg)

		// Crisis module
	} else if xplac.MsgType == mcrisis.CrisisInvariantBrokenMsgType {
		convertMsg, _ := xplac.Msg.(crisistypes.MsgVerifyInvariant)
		builder.SetMsgs(&convertMsg)

		// Distribution module
	} else if xplac.MsgType == mdist.DistributionFundCommunityPoolMsgType {
		convertMsg, _ := xplac.Msg.(disttypes.MsgFundCommunityPool)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mdist.DistributionProposalCommunityPoolSpendMsgType {
		convertMsg, _ := xplac.Msg.(govtypes.MsgSubmitProposal)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mdist.DistributionWithdrawRewardsMsgType {
		convertMsg, _ := xplac.Msg.([]sdk.Msg)
		builder.SetMsgs(convertMsg...)

	} else if xplac.MsgType == mdist.DistributionWithdrawAllRewardsMsgType {
		convertMsg, _ := xplac.Msg.([]sdk.Msg)
		builder.SetMsgs(convertMsg...)

	} else if xplac.MsgType == mdist.DistributionSetWithdrawAddrMsgType {
		convertMsg, _ := xplac.Msg.(disttypes.MsgSetWithdrawAddress)
		builder.SetMsgs(&convertMsg)

		// Feegrant module
	} else if xplac.MsgType == mfeegrant.FeegrantGrantMsgType {
		convertMsg, _ := xplac.Msg.(feegrant.MsgGrantAllowance)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mfeegrant.FeegrantRevokeGrantMsgType {
		convertMsg, _ := xplac.Msg.(feegrant.MsgRevokeAllowance)
		builder.SetMsgs(&convertMsg)

		// Gov module
	} else if xplac.MsgType == mgov.GovSubmitProposalMsgType {
		convertMsg, _ := xplac.Msg.(govtypes.MsgSubmitProposal)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mgov.GovDepositMsgType {
		convertMsg, _ := xplac.Msg.(govtypes.MsgDeposit)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mgov.GovVoteMsgType {
		convertMsg, _ := xplac.Msg.(govtypes.MsgVote)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mgov.GovWeightedVoteMsgType {
		convertMsg, _ := xplac.Msg.(govtypes.MsgVoteWeighted)
		builder.SetMsgs(&convertMsg)

		// Params module
	} else if xplac.MsgType == mparams.ParamsProposalParamChangeMsgType {
		convertMsg, _ := xplac.Msg.(govtypes.MsgSubmitProposal)
		builder.SetMsgs(&convertMsg)

		// Reward module
	} else if xplac.MsgType == mreward.RewardFundFeeCollectorMsgType {
		convertMsg, _ := xplac.Msg.(rewardtypes.MsgFundFeeCollector)
		builder.SetMsgs(&convertMsg)

		// slashing module
	} else if xplac.MsgType == mslashing.SlahsingUnjailMsgType {
		convertMsg, _ := xplac.Msg.(slashingtypes.MsgUnjail)
		builder.SetMsgs(&convertMsg)

		// Staking module
	} else if xplac.MsgType == mstaking.StakingCreateValidatorMsgType {
		convertMsg, _ := xplac.Msg.(sdk.Msg)
		builder.SetMsgs(convertMsg)

	} else if xplac.MsgType == mstaking.StakingEditValidatorMsgType {
		convertMsg, _ := xplac.Msg.(stakingtypes.MsgEditValidator)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mstaking.StakingDelegateMsgType {
		convertMsg, _ := xplac.Msg.(stakingtypes.MsgDelegate)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mstaking.StakingUnbondMsgType {
		convertMsg, _ := xplac.Msg.(stakingtypes.MsgUndelegate)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mstaking.StakingRedelegateMsgType {
		convertMsg, _ := xplac.Msg.(stakingtypes.MsgBeginRedelegate)
		builder.SetMsgs(&convertMsg)

		// Upgrade module
	} else if xplac.MsgType == mupgrade.UpgradeProposalSoftwareUpgradeMsgType {
		convertMsg, _ := xplac.Msg.(govtypes.MsgSubmitProposal)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mupgrade.UpgradeCancelSoftwareUpgradeMsgType {
		convertMsg, _ := xplac.Msg.(govtypes.MsgSubmitProposal)
		builder.SetMsgs(&convertMsg)

		// Wasm module
	} else if xplac.MsgType == mwasm.WasmStoreMsgType {
		convertMsg, _ := xplac.Msg.(wasm.MsgStoreCode)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mwasm.WasmInstantiateMsgType {
		convertMsg, _ := xplac.Msg.(wasm.MsgInstantiateContract)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mwasm.WasmExecuteMsgType {
		convertMsg, _ := xplac.Msg.(wasm.MsgExecuteContract)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mwasm.WasmClearContractAdminMsgType {
		convertMsg, _ := xplac.Msg.(wasm.MsgClearAdmin)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mwasm.WasmSetContractAdminMsgType {
		convertMsg, _ := xplac.Msg.(wasm.MsgUpdateAdmin)
		builder.SetMsgs(&convertMsg)

	} else if xplac.MsgType == mwasm.WasmMigrateMsgType {
		convertMsg, _ := xplac.Msg.(wasm.MsgMigrateContract)
		builder.SetMsgs(&convertMsg)

	} else {
		return nil, util.LogErr(errors.ErrInvalidMsgType, xplac.MsgType)
	}

	return builder, nil
}

// Set information for transaction builder.
func convertAndSetBuilder(xplac *XplaClient, builder cmclient.TxBuilder) cmclient.TxBuilder {
	if xplac.Opts.FeeAmount != "" {
		feeAmountDenomRemove, _ := util.FromStringToBigInt(util.DenomRemove(xplac.Opts.FeeAmount))
		feeAmountCoin := sdk.Coin{
			Amount: sdk.NewIntFromBigInt(feeAmountDenomRemove),
			Denom:  types.XplaDenom,
		}
		feeAmountCoins := sdk.NewCoins(feeAmountCoin)
		builder.SetFeeAmount(feeAmountCoins)
	}

	if xplac.Opts.GasLimit != "" {
		builder.SetGasLimit(util.FromStringToUint64(xplac.Opts.GasLimit))
	}

	if xplac.Opts.TimeoutHeight != "" {
		builder.SetTimeoutHeight(util.FromStringToUint64(xplac.Opts.TimeoutHeight))
	}

	builder.SetFeeGranter(xplac.Opts.FeeGranter)

	if types.Memo != "" {
		builder.SetMemo(types.Memo)
	}

	return builder
}

// Sign transaction by using given private key.
func txSignRound(xplac *XplaClient,
	sigsV2 []signing.SignatureV2,
	privs []cryptotypes.PrivKey,
	accSeqs []uint64,
	accNums []uint64,
	builder cmclient.TxBuilder) error {

	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  xplac.Opts.SignMode,
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}
		sigsV2 = append(sigsV2, sigV2)
	}

	err := builder.SetSignatures(sigsV2...)
	if err != nil {
		return util.LogErr(errors.ErrParse, err)
	}

	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			ChainID:       xplac.ChainId,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		sigV2, err := tx.SignWithPrivKey(
			xplac.Opts.SignMode,
			signerData,
			builder,
			priv,
			xplac.EncodingConfig.TxConfig,
			accSeqs[i],
		)
		if err != nil {
			return util.LogErr(errors.ErrParse, err)
		}

		sigsV2 = append(sigsV2, sigV2)
	}

	err = builder.SetSignatures(sigsV2...)
	if err != nil {
		return util.LogErr(errors.ErrParse, err)
	}

	return nil
}

// Sign evm transaction by using given private key.
func evmTxSignRound(xplac *XplaClient,
	toAddr common.Address,
	gasPrice *big.Int,
	amount *big.Int,
	invokeByteData []byte,
	chainId *big.Int,
	ethPrivKey *ecdsa.PrivateKey) ([]byte, error) {

	tx := evmtypes.NewTransaction(
		util.FromStringToUint64(xplac.Opts.Sequence),
		toAddr,
		amount,
		util.FromStringToUint64(xplac.Opts.GasLimit),
		gasPrice,
		invokeByteData,
	)

	signer := evmtypes.NewEIP155Signer(chainId)

	signedTx, err := evmtypes.SignTx(tx, signer, ethPrivKey)
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}
	txbytes, err := signedTx.MarshalJSON()
	if err != nil {
		return nil, util.LogErr(errors.ErrFailedToMarshal, err)
	}

	return txbytes, nil
}

// Read transaction file and make standard transaction.
func readTxAndInitContexts(clientCtx cmclient.Context, filename string) (cmclient.Context, tx.Factory, sdk.Tx, error) {
	stdTx, err := authclient.ReadTxFromFile(clientCtx, filename)
	if err != nil {
		return clientCtx, tx.Factory{}, nil, util.LogErr(errors.ErrParse, err)
	}

	txFactory := util.NewFactory(clientCtx)

	return clientCtx, txFactory, stdTx, nil
}

// Marshal signature type JSON.
func marshalSignatureJSON(txConfig cmclient.TxConfig, txBldr cmclient.TxBuilder, signatureOnly bool) ([]byte, error) {
	parsedTx := txBldr.GetTx()
	if signatureOnly {
		sigs, err := parsedTx.GetSignaturesV2()
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}
		return txConfig.MarshalSignatureJSON(sigs)
	}

	return txConfig.TxJSONEncoder()(parsedTx)
}

// Unmarshal signature type JSON.
func unmarshalSignatureJSON(clientCtx cmclient.Context, filename string) (sigs []signing.SignatureV2, err error) {
	var bytes []byte
	if bytes, err = os.ReadFile(filename); err != nil {
		return
	}
	return clientCtx.TxConfig.UnmarshalSignatureJSON(bytes)
}

// The secp-256k1 private key converts ECDSA privatkey for using evm module.
func toECDSA(privKey key.PrivateKey) (*ecdsa.PrivateKey, error) {
	return ethcrypto.ToECDSA(privKey.Bytes())
}

// Get multiple signatures information. It returns keyring of cosmos sdk.
func getMultisigInfo(clientCtx cmclient.Context, name string) (keyring.Info, error) {
	kb := clientCtx.Keyring
	multisigInfo, err := kb.Key(name)
	if err != nil {
		return nil, util.LogErr(errors.ErrKeyNotFound, "error getting keybase multisig account", err)
	}
	if multisigInfo.GetType() != keyring.TypeMulti {
		return nil, util.LogErr(errors.ErrInvalidMsgType, name, "must be of type", keyring.TypeMulti, ":", multisigInfo.GetType())
	}

	return multisigInfo, nil
}
