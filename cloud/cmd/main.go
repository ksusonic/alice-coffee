package main

import (
	"log"
	"time"

	"github.com/ksusonic/alice-coffee/config"
	"github.com/ksusonic/alice-coffee/pkg/dialogs"
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
		AutoPong: conf.AutoPong,
	})

	updates.Loop(func(k dialogs.Kit) *dialogs.Response {
		req, resp := k.Init()
		if req.IsNewSession() {
			return resp.Text("Hello World!")
		}
		return resp.Text(req.Text())
	})
}
