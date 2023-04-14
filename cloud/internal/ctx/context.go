package ctx

import (
	"context"

	"go.uber.org/zap"
)

type WebSocket interface {
	SendJSON(jsonData interface{}) error
	CheckActive() bool
}

type GlobalCtx struct {
	Socket WebSocket
}

type SceneCtx struct {
	Ctx       context.Context
	GlobalCtx *GlobalCtx
	Logger    *zap.SugaredLogger
	ReqId     string
}
