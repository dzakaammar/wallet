package rest

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

//go:generate go run ../../main.go openapi-gen --path .

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

	swagger.Components.SecuritySchemes = openapi3.SecuritySchemes{
		"bearerAuth": &openapi3.SecuritySchemeRef{
			Value: openapi3.NewSecurityScheme().WithType("http").WithScheme("basic"),
		},
	}

	swagger.Components.RequestBodies = openapi3.RequestBodies{
		"CreateUserRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Request used for registering a user").
				WithRequired(true).
				WithJSONSchema(openapi3.NewSchema().
					WithProperty("username", openapi3.NewStringSchema().
						WithMinLength(1))),
		},
		"TopupRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Request used for topup money").
				WithRequired(true).
				WithJSONSchema(openapi3.NewSchema().
					WithProperty("amount", openapi3.NewInt32Schema().
						WithMin(1))),
		},
	}

	swagger.Components.Responses = openapi3.Responses{
		"EmptyResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse(),
		},
		"CreateUserResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().WithDescription("create user success response").
				WithJSONSchema(openapi3.NewSchema().WithProperty("token", openapi3.NewStringSchema())),
		},
	}

	swagger.Paths = openapi3.Paths{
		"/create_user": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "CreateUser",
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/CreateUserRequest",
				},
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/CreateUserResponse",
					},
					"401": &openapi3.ResponseRef{
						Value: openapi3.NewResponse().WithDescription("username already exists"),
					},
				},
			},
		},
		"/balance_topup": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "BalanceTopup",
				Security: &openapi3.SecurityRequirements{
					{
						"bearerAuth": []string{},
					},
				},
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/TopupRequest",
				},
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/EmptyResponse",
					},
					"401": &openapi3.ResponseRef{
						Value: openapi3.NewResponse().WithDescription("username already exists"),
					},
				},
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
