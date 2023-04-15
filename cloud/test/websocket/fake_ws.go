package websocket

import (
	"errors"
	"math/rand"
	"net/http"

	"go.uber.org/zap"
)

type FakeWebsocket struct {
	Logger *zap.SugaredLogger
}

func (f FakeWebsocket) Router() http.Handler {
	return http.NotFoundHandler()
}

func (f FakeWebsocket) SendJSON(jsonData interface{}) error {
	if rand.Intn(100) < 10 { // 10% chance to fail
		return errors.New("got unexpected error")
	}

	f.Logger.Info("fake-sent: %v", jsonData)
	return nil
}

func (f FakeWebsocket) CheckActive() bool {
	return rand.Intn(100) < 10
}
