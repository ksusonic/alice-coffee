package main

import (
	"time"

	"github.com/ksusonic/alice-coffee/cloud/config"
	"github.com/ksusonic/alice-coffee/cloud/internal/scenario"
	"github.com/ksusonic/alice-coffee/cloud/internal/scenario/nlg"
	"github.com/ksusonic/alice-coffee/cloud/internal/ws"
	"github.com/ksusonic/alice-coffee/cloud/pkg/dialogs"
	"go.uber.org/zap"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic("Could not load config: " + err.Error())
	}
	logger := getLogger(conf.Debug)

	c := ws.NewWebsocket(logger.Desugar().Named("websocket"))
	updates := dialogs.StartServer("/hook", dialogs.ServerConf{
		Address:  conf.Address,
		Debug:    conf.Debug,
		Timeout:  time.Duration(conf.Dialogs.Timeout),
		AutoPong: true,
	}, c.Router)

	ctx := scenario.Context{
		Logger: logger,
	}
	d := scenario.NewIntentDispatcher(&ctx)

	updates.Loop(func(k dialogs.Kit) *dialogs.Response {
		req, resp := k.Init()

		if req.IsNewSession() {
			return resp.Text(nlg.Greeting)
		}

		result := d.TryResponse(req, resp)
		if result != nil {
			return result
		}

		logger.Info("Could not make relevant response")
		return resp.TextWithTTS(nlg.IrrelevantPhrase())
	})
}

func getLogger(debug bool) *zap.SugaredLogger {
	getSuggared := func(logger *zap.Logger, err error) *zap.SugaredLogger {
		if err != nil {
			panic(err)
		}
		return logger.Sugar()
	}

	if debug {
		return getSuggared(zap.NewDevelopment())
	} else {
		return getSuggared(zap.NewProduction())
	}
}
