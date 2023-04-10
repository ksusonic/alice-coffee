package dialogs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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
func StartServer(hookPath string, conf ServerConf, routers ...chi.Router) Stream {

	stream := make(chan Kit, 1)
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	for _, r := range routers {
		router.Mount("/", r)
	}

	router.HandleFunc(hookPath, webhook(conf, stream))

	go func() {
		err := http.ListenAndServe(conf.Address, router)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Started alice-dialogs server on", conf.Address)
	return stream
}

type ServerConf struct {
	Address  string
	Debug    bool
	Timeout  time.Duration
	AutoPong bool
}

func webhook(conf ServerConf, stream chan<- Kit) http.HandlerFunc {
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
