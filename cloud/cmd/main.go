package main

import (
	"log"
	"time"

	"github.com/ksusonic/alice-coffee/config"
	"github.com/ksusonic/alice-coffee/internal/queue"
	"github.com/ksusonic/alice-coffee/pkg/dialogs"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func main() {
	conf, err := config.LoadConfig()

	// TODO
	_ = queue.NewMessageQueue(conf.QueueUrl, aws.Config{
		Region:                      queue.REGION,
		Credentials:                 queue.NewCredentialsProvider(conf.SqsAccessKey, conf.SqsSecretKey),
		EndpointResolverWithOptions: queue.NewEndpointResolverWithOptions(),
	})

	if err != nil {
		log.Fatal(err)
	}
	updates := dialogs.StartServer("/hook", dialogs.ServerConf{
		Address:  conf.Address,
		Debug:    conf.Debug,
		Timeout:  time.Duration(conf.Timeout),
		AutoPong: true,
	})

	updates.Loop(func(k dialogs.Kit) *dialogs.Response {
		req, resp := k.Init()
		if req.IsNewSession() {
			return resp.Text("Hello World!")
		}
		return resp.Text(req.Text())
	})
}
