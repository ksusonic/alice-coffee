package ws

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WebSocket struct {
	Router   chi.Router
	upgrader websocket.Upgrader

	Logger *zap.Logger
}

func NewWebsocket(logger *zap.Logger) *WebSocket {
	c := &WebSocket{
		Router: chi.NewRouter(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		Logger: logger,
	}

	c.Router.HandleFunc("/ws", c.wsEndpoint)
	return c
}

func (c *WebSocket) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.Logger.Error("conn upgrade failed: %v", zap.Error(err))
		return
	}
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			c.Logger.Info("read: %v", zap.Error(err))
			break
		}
		log.Printf("recv: %s", message)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			c.Logger.Info("write: %v", zap.Error(err))
			break
		}
	}
}
