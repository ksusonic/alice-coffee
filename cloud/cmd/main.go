package main

import (
	"github.com/ksusonic/alice-coffee/config"
	"github.com/ksusonic/alice-coffee/pkg/dialogs"
	"log"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	updates := dialogs.StartServer("/hook", conf)

	updates.Loop(func(k dialogs.Kit) *dialogs.Response {
		req, resp := k.Init()
		if req.IsNewSession() {
			return resp.Text("Hello World!")
		}
		return resp.Text(req.Text())
	})
}
