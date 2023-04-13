package scene

import (
	"github.com/ksusonic/alice-coffee/cloud/internal/ctx"
	"github.com/ksusonic/alice-coffee/cloud/internal/scenario/nlg"
	"github.com/ksusonic/alice-coffee/cloud/pkg/dialogs"
)

func WhatCanYouDo(
	_ *ctx.SceneCtx,
	_ *dialogs.Request,
	_ string,
	_ dialogs.Slots,
	resp *dialogs.Response) *dialogs.Response {
	return resp.Text(nlg.WhatCanYouDo)
}
