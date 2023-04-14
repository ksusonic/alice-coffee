package ws

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WebSocket struct {
	Router   chi.Router
	upgrader websocket.Upgrader
	conn     *websocket.Conn

	Logger *zap.Logger
}

func (c *WebSocket) CheckActive() bool {
	return c.conn != nil
}

func (c *WebSocket) SendJSON(jsonData interface{}) error {
	if c.conn == nil {
		c.Logger.Debug("could not send data", zap.Any("data", jsonData))
		return errors.New("no connection with websocket")
	}
	err := c.conn.WriteJSON(jsonData)
	if err != nil {
		c.Logger.Debug("could not send data", zap.Any("data", jsonData))
		return err
	}
	return nil
}

func NewWebsocket(logger *zap.Logger) *WebSocket {
	c := &WebSocket{
		Router: chi.NewRouter(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				username, password, _ := r.BasicAuth()
				logger.Debug("check origin",
					zap.String("host", r.Host),
					zap.String("user-agent", r.UserAgent()),
					zap.String("username", username),
					zap.String("password", password),
				)
				return true
			},
		},
		Logger: logger,
	}

	c.Router.HandleFunc("/", c.wsEndpoint)
	return c
}

func (c *WebSocket) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	if c.conn != nil {
		c.Logger.Error("has already connected with", zap.String("host", c.conn.RemoteAddr().String()))
		w.WriteHeader(http.StatusGone)
		return
	}

	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.Logger.Error("conn upgrade failed: %v", zap.Error(err))
		return
	}
	c.conn = conn
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			c.Logger.Error("error closing connection", zap.Error(err))
		}
		c.conn = nil
	}(conn)
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			c.Logger.Info("read: %v", zap.Error(err))
			break
		}
		c.Logger.Info("received", zap.ByteString("msg", message))
		err = conn.WriteMessage(mt, message)
		if err != nil {
			c.Logger.Info("write: %v", zap.Error(err))
			break
		}
	}
}
