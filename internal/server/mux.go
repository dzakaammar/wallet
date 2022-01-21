package server

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"wallet/internal/rest"

	"github.com/gorilla/mux"
)

//go:embed swagger
var swagger embed.FS

type HTTPServer struct {
	http.Server
}

func NewMuxHTTPServer(h *rest.Handler) *HTTPServer {
	r := mux.NewRouter()

	h.Register(r)

	rest.RegisterOpenAPI3(r)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(swagger))))

	return &HTTPServer{
		Server: http.Server{
			Handler: r,
		},
	}
}

func (h *HTTPServer) Start(port int) error {
	fmt.Printf("HTTP server is starting at port %s\n", h.Addr)
	h.Addr = fmt.Sprintf(":%d", port)
	return h.ListenAndServe()
}

func (h *HTTPServer) Stop(ctx context.Context) error {
	fmt.Println("HTTP server is shutting down...")
	return h.Shutdown(ctx)
}
