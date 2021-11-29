package rest

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

//go:generate go run ../../main.go openapi-gen --path .
//go:generate oapi-codegen -package openapi3 -generate types  -o ../../pkg/openapi3/task_types.gen.go openapi3.yaml
//go:generate oapi-codegen -package openapi3 -generate client -o ../../pkg/openapi3/client.gen.go     openapi3.yaml

func NewOpenAPI3() openapi3.T {
	swagger := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "Simple Wallet",
			Description: "REST APIs used for interacting with the Simple Walleto Service",
			Version:     "0.0.0",
			License: &openapi3.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			Contact: &openapi3.Contact{
				URL: "https://github.com/dzakaammar/wallet",
			},
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				Description: "Local development",
				URL:         "http://127.0.0.1:8080",
			},
		},
	}

	return swagger
}

func RegisterOpenAPI3(router *mux.Router) {
	swagger := NewOpenAPI3()

	router.HandleFunc("/openapi3.json", func(w http.ResponseWriter, r *http.Request) {
		writeJSONResponse(w, http.StatusOK, &swagger)
	}).Methods(http.MethodGet)

	router.HandleFunc("/openapi3.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")

		data, _ := yaml.Marshal(&swagger)

		_, _ = w.Write(data)

		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
}
