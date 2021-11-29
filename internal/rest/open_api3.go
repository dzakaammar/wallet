package rest

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

func NewOpenAPI3() openapi3.T {
	swagger := openapi3.T{}

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
