package main

import (
	"github.com/ksusonic/alice-coffee/pkg/dialogs"
)

func main() {
	config := dialogs.LoadConfig(dialogs.ParseArgs())
	updates := dialogs.StartServer("/hook", config)

	updates.Loop(func(k dialogs.Kit) *dialogs.Response {
		req, resp := k.Init()
		if req.IsNewSession() {
			return resp.Text("Hello World!")
		}
		return resp.Text(req.Text())
	})
}
