package qtest

import (
	"github.com/Moonyongjung/xpla.go/client"
	"github.com/Moonyongjung/xpla.go/util/testutil"
)

func NewTestXplaClient() *client.XplaClient {
	return client.NewXplaClient(testutil.TestChainId)
}

func ResetXplac(xplac *client.XplaClient) *client.XplaClient {
	xplac.WithOptions(client.Options{})
	xplac.Module = ""
	xplac.MsgType = ""
	xplac.Msg = nil
	xplac.Err = nil

	return xplac
}
