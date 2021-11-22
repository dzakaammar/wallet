package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("register a user", func() {
	InitHandler()
	// repository
	Context("register", func() {
		When("input is not valid", func() {
			body := strings.NewReader(`{"username": ""}`)
			req := httptest.NewRequest(http.MethodPost, "/create_user", body)
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)

			It("gives 400", func() {
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("input is valid", func() {
			body := strings.NewReader(`{"username": "test"}`)
			req := httptest.NewRequest(http.MethodPost, "/create_user", body)
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)

			It("gives 201", func() {
				Expect(resp.Code).To(Equal(http.StatusCreated))
			})

			It("gives the user token", func() {
				var respBody struct {
					Token string `json:"token"`
				}
				err := json.NewDecoder(resp.Body).Decode(&respBody)
				Expect(err).To(BeNil())
				Expect(respBody.Token).To(Not(BeEmpty()))
			})
		})
	})

	Context("register with existing username", func() {
		body := strings.NewReader(`{"username": "test"}`)
		req := httptest.NewRequest(http.MethodPost, "/create_user", body)
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)

		It("gives 409", func() {
			Expect(resp.Code).To(Equal(http.StatusConflict))
		})
	})
})

var _ = Describe("check balance", func() {
	InitHandler()

	body := strings.NewReader(`{"username": "test"}`)
	req := httptest.NewRequest(http.MethodPost, "/create_user", body)
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	var respBody struct {
		Token string `json:"token"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	Context("unauthorized", func() {
		req := httptest.NewRequest(http.MethodGet, "/balance_read", nil)
		req.Header.Set("Authorization", "")
		resp := httptest.NewRecorder()

		h.ServeHTTP(resp, req)
		It("gives 404", func() {
			Expect(resp.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("authorized", func() {
		req := httptest.NewRequest(http.MethodGet, "/balance_read", nil)
		req.Header.Set("Authorization", respBody.Token)
		resp := httptest.NewRecorder()

		h.ServeHTTP(resp, req)

		It("gives 200", func() {
			Expect(resp.Code).To(Equal(http.StatusOK))
		})

		It("gives balance", func() {
			var respBody struct {
				Balance int `json:"balance"`
			}

			err := json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).To(BeNil())
			Expect(respBody.Balance).To(Equal(0))
		})
	})
})

var _ = Describe("topup balance", func() {
	InitHandler()

	body := strings.NewReader(`{"username": "test"}`)
	req := httptest.NewRequest(http.MethodPost, "/create_user", body)
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)

	var respBody struct {
		Token string `json:"token"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	Context("user unauthorized", func() {
		body := strings.NewReader(`{"amount": 1000}`)
		req := httptest.NewRequest(http.MethodPost, "/balance_topup", body)
		resp := httptest.NewRecorder()

		h.ServeHTTP(resp, req)

		It("gives 401", func() {
			Expect(resp.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("user authorized", func() {
		When("amount is minus", func() {
			body := strings.NewReader(`{"amount": -1}`)
			req := httptest.NewRequest(http.MethodPost, "/balance_topup", body)
			req.Header.Set("Authorization", respBody.Token)
			resp := httptest.NewRecorder()

			h.ServeHTTP(resp, req)

			It("gives 400", func() {
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("amount is 10000000", func() {
			body := strings.NewReader(`{"amount": 10000000}`)
			req := httptest.NewRequest(http.MethodPost, "/balance_topup", body)
			req.Header.Set("Authorization", respBody.Token)
			resp := httptest.NewRecorder()

			h.ServeHTTP(resp, req)

			It("gives 400", func() {
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("amount is valid", func() {
			body := strings.NewReader(`{"amount": 1000000}`)
			req := httptest.NewRequest(http.MethodPost, "/balance_topup", body)
			req.Header.Set("Authorization", respBody.Token)
			resp := httptest.NewRecorder()

			h.ServeHTTP(resp, req)

			It("gives 204", func() {
				Expect(resp.Code).To(Equal(http.StatusNoContent))
			})
		})
	})

	Context("check balance after topup", func() {
		req := httptest.NewRequest(http.MethodGet, "/balance_read", nil)
		req.Header.Set("Authorization", respBody.Token)
		resp := httptest.NewRecorder()

		h.ServeHTTP(resp, req)

		It("gives 200", func() {
			Expect(resp.Code).To(Equal(http.StatusOK))
		})

		It("gives balance", func() {
			var respBody struct {
				Balance int `json:"balance"`
			}

			err := json.NewDecoder(resp.Body).Decode(&respBody)
			Expect(err).To(BeNil())
			Expect(respBody.Balance).To(Equal(1000000))
		})
	})
})
