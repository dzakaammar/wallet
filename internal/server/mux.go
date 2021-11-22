package server

import (
	"context"
	"fmt"
	"net/http"
	"wallet/internal/rest"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	http.Server
}

func NewMuxHTTPServer(h *rest.Handler, port int) *HTTPServer {
	r := mux.NewRouter()

	h.Register(r)

	return &HTTPServer{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		},
	}
}

func (h *HTTPServer) Start() error {
	fmt.Printf("HTTP server is starting at port %s\n", h.Addr)
	return h.ListenAndServe()
}

func (h *HTTPServer) Stop(ctx context.Context) {
	fmt.Println("HTTP server is shutting down...")
	h.Shutdown(ctx)
}
