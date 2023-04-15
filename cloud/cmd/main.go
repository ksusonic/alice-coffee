package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ksusonic/alice-coffee/cloud/config"
	customCtx "github.com/ksusonic/alice-coffee/cloud/internal/ctx"
	"github.com/ksusonic/alice-coffee/cloud/internal/scenario"
	"github.com/ksusonic/alice-coffee/cloud/internal/server"
	"github.com/ksusonic/alice-coffee/cloud/internal/ws"
	"github.com/ksusonic/alice-coffee/cloud/pkg/dialogs"
	"github.com/ksusonic/alice-coffee/cloud/test/websocket"
	"go.uber.org/zap"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic("Could not load config: " + err.Error())
	}
	logger := getLogger(conf.Debug)

	logger.Info("setup stage", zap.Any("config", conf))

	webSocket := getWebsocket(conf.DemoMode, logger)
	dialogsRouter, updates := dialogs.Router(dialogs.ServerConf{
		Debug:    conf.Debug,
		Timeout:  time.Duration(conf.Dialogs.Timeout),
		AutoPong: true,
	})

	dispatcher := scenario.NewIntentDispatcher(&customCtx.GlobalCtx{
		Socket: webSocket,
	}, logger.Named("dispatcher").Sugar())

	s := server.NewServer(&conf.Server, logger.Named("server"),
		server.Route{
			Pattern: "/hook",
			Handler: dialogsRouter,
		},
		server.Route{
			Pattern: "/ws",
			Handler: webSocket.Router(),
		})
	srv := s.Run() // start server

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go updates.Loop(ctx, dispatcher.Handler)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	logger.Debug("caught " + (<-osSignal).String())

	toCtx, toCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer toCancel()

	if srvErr := srv.Shutdown(toCtx); srvErr != nil {
		logger.Fatal("server shutdown error", zap.Error(srvErr))
	}

	logger.Info("server stopped")
}

func getLogger(debug bool) *zap.Logger {
	retireLogger := func(logger *zap.Logger, err error) *zap.Logger {
		if err != nil {
			panic(err)
		}
		return logger
	}

	if debug {
		return retireLogger(zap.NewDevelopment())
	} else {
		return retireLogger(zap.NewProduction())
	}
}

func getWebsocket(demo bool, logger *zap.Logger) customCtx.WebSocket {
	namedLogger := logger.Named("websocket")
	if demo {
		return websocket.FakeWebsocket{Logger: namedLogger.Sugar()}
	} else {
		return ws.NewWebsocket(namedLogger)
	}
}
