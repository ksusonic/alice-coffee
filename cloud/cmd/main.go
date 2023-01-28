package main

import (
	"log"
	"time"

	"github.com/ksusonic/alice-coffee/cloud/config"
	"github.com/ksusonic/alice-coffee/cloud/internal/nlg"
	"github.com/ksusonic/alice-coffee/cloud/internal/queue"
	"github.com/ksusonic/alice-coffee/cloud/internal/scenario"
	"github.com/ksusonic/alice-coffee/cloud/pkg/dialogs"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	updates := dialogs.StartServer("/hook", dialogs.ServerConf{
		Address:  conf.Address,
		Debug:    conf.Debug,
		Timeout:  time.Duration(conf.Timeout),
		AutoPong: true,
	})

	ctx := scenario.Context{
		MessageQueue: queue.NewMessageQueue(conf.QueueUrl, conf.SqsAccessKey, conf.SqsSecretKey),
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

		log.Println("Could not make relevant response")
		return resp.TextWithTTS(req.Text(), req.Text())
	})
}
