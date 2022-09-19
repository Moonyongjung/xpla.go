package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/context/ctxhttp"

	"github.com/Moonyongjung/xpla.go/util"
	cmclient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	userInfoUrl  = "/cosmos/auth/v1beta1/accounts/"
	simulateUrl  = "/cosmos/tx/v1beta1/simulate"
	broadcastUrl = "/cosmos/tx/v1beta1/txs"
)

// LoadAccount simulates gas and fee for a transaction
func (xplac *XplaClient) LoadAccount(address sdk.AccAddress) (res authtypes.AccountI, err error) {
	out, err := ctxHttpClient("GET", xplac.Opts.LcdURL+userInfoUrl+address.String(), nil)
	if err != nil {
		return nil, err
	}

	var response authtypes.QueryAccountResponse
	err = xplac.EncodingConfig.Marshaler.UnmarshalJSON(out, &response)
	if err != nil {
		return nil, util.LogErr(err, "failed to unmarshal response")
	}

	return response.Account.GetCachedValue().(authtypes.AccountI), nil
}

// Simulate tx and get response
func (xplac *XplaClient) Simulate(txbuilder cmclient.TxBuilder) (*sdktx.SimulateResponse, error) {
	sig := signing.SignatureV2{
		PubKey: xplac.PrivateKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode: xplac.Opts.SignMode,
		},
		Sequence: util.FromStringToUint64(xplac.Opts.Sequence),
	}

	if err := txbuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	sdkTx := txbuilder.GetTx()
	txBytes, err := xplac.EncodingConfig.TxConfig.TxEncoder()(sdkTx)
	if err != nil {
		return nil, err
	}

	reqBytes, err := xplac.EncodingConfig.Marshaler.MarshalJSON(&sdktx.SimulateRequest{
		TxBytes: txBytes,
	})
	if err != nil {
		return nil, err
	}

	out, err := ctxHttpClient("POST", xplac.Opts.LcdURL+simulateUrl, reqBytes)
	if err != nil {
		return nil, err
	}

	var response sdktx.SimulateResponse
	err = xplac.EncodingConfig.Marshaler.UnmarshalJSON(out, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Make new http client for inquiring several information.
func ctxHttpClient(methodType string, url string, reqBody []byte) ([]byte, error) {
	var resp *http.Response
	var err error
	ctx := context.Background()
	httpClient := &http.Client{Timeout: 30 * time.Second}

	if methodType == "GET" {
		resp, err = ctxhttp.Get(ctx, httpClient, url)
		if err != nil {
			return nil, util.LogErr(err, "failed GET method")
		}
	} else if methodType == "POST" {
		resp, err = ctxhttp.Post(ctx, httpClient, url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, util.LogErr(err, "failed POST method")
		}
	} else {
		return nil, util.LogErr(err, "not correct method")
	}

	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, util.LogErr(err, "failed to read response")
	}

	if resp.StatusCode != 200 {
		return nil, util.LogErr(resp.StatusCode, ":", string(out))
	}

	return out, nil
}