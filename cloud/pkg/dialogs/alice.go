package dialogs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Kit структура для передачи данных в главный цикл.
type Kit struct {
	Req  *Request
	Resp *Response
	// Ctx позволяет получить сигнал об истечении периода ожидания ответа.
	Ctx context.Context

	c chan<- *Response
}

// Init получает входящий пакет и заготовку исходящего из данных запроса.
func (k Kit) Init() (*Request, *Response) {
	return k.Req, k.Resp
}

// Stream канал, передающий данные в основной цикл.
type Stream <-chan Kit

// Handler сигнатура функции, передаваемой методу Loop().
type Handler func(k Kit) *Response

// Loop отвечает за работу главного цикла.
func (updates Stream) Loop(f Handler) {
	for kit := range updates {
		go func(k Kit) {
			k.c <- f(k)
			close(k.c)
		}(kit)
	}
}

// StartServer регистрирует обработчик входящих пакетов.
func StartServer(hookPath string, config Config) Stream {

	stream := make(chan Kit, 1)
	http.HandleFunc(hookPath, webhook(config, stream))

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
	}()
	// TODO log start
	return stream
}

func webhook(conf Config, stream chan<- Kit) http.HandlerFunc {
	reqPool := sync.Pool{
		New: func() interface{} {
			return new(Request)
		},
	}

	respPool := sync.Pool{
		New: func() interface{} {
			return new(Response)
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(conf.Timeout)*time.Millisecond)
		defer cancel()

		if conf.Debug {
			requestDump, err := httputil.DumpRequest(r, true)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(string(requestDump))
		}

		req := reqPool.Get().(*Request)
		defer reqPool.Put(req)

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(req.clean()); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := respPool.Get().(*Response)
		resp.clean().prepareResponse(req)
		defer respPool.Put(resp)

		if conf.AutoPong {
			if req.Type() == SimpleUtterance && req.Text() == "ping" {
				if md, err := json.Marshal(resp.Text("pong")); err == nil {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write(md)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		req.Bearer = r.Header.Get("Authorization")

		back := make(chan *Response)
		stream <- Kit{
			Req:  req,
			Resp: resp,
			Ctx:  ctx,

			c: back,
		}

		var response *Response
		select {
		case <-ctx.Done():
			log.Println(ctx.Err())
			w.WriteHeader(http.StatusInternalServerError)
			return
		case response = <-back:
		}

		writer := io.Writer(w)

		if conf.Debug {
			var buf bytes.Buffer
			writer = io.MultiWriter(w, &buf)
			defer func() {
				fmt.Printf("\n%s\n\n", buf.String())
			}()
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(writer)
		if err := encoder.Encode(&response); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
