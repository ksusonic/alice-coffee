package ctx

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type WebSocket interface {
	SendJSON(jsonData interface{}) error
	CheckActive() bool
	Router() http.Handler
}

type GlobalCtx struct {
	Socket WebSocket
}

type SceneCtx struct {
	Ctx       context.Context
	GlobalCtx *GlobalCtx
	Logger    *zap.SugaredLogger
	ReqID     string
}
