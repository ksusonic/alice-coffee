package scenario

import (
	"go.uber.org/zap"
)

type WebSocket interface {
	SendMessage(data []byte) error
}

type Context struct {
	Logger *zap.SugaredLogger
	Socket WebSocket
}
