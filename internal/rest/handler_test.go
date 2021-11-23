package rest_test

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

var _ = Describe("transfer", func() {
	InitHandler()

	body := strings.NewReader(`{"username": "alice"}`)
	req := httptest.NewRequest(http.MethodPost, "/create_user", body)
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)

	var alice struct {
		Token string `json:"token"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&alice)

	body = strings.NewReader(`{"username": "bob"}`)
	req = httptest.NewRequest(http.MethodPost, "/create_user", body)
	resp = httptest.NewRecorder()
	h.ServeHTTP(resp, req)

	var bob struct {
		Token string `json:"token"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&bob)

	Context("user unauthorized", func() {
		body := strings.NewReader(`{"to_username": "bob", "amount": 100}`)
		req := httptest.NewRequest(http.MethodPost, "/transfer", body)
		req.Header.Set("Authorization", "test")
		resp := httptest.NewRecorder()

		h.ServeHTTP(resp, req)

		It("gives 401", func() {
			Expect(resp.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("user authorized", func() {
		When("destination user not found", func() {
			body := strings.NewReader(`{"to_username": "other_user", "amount": 100}`)
			req := httptest.NewRequest(http.MethodPost, "/transfer", body)
			req.Header.Set("Authorization", alice.Token)
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)

			It("gives 404", func() {
				Expect(resp.Code).To(Equal(http.StatusNotFound))
			})
		})

		When("insufficient balance", func() {
			body := strings.NewReader(`{"to_username": "bob", "amount": 100}`)
			req := httptest.NewRequest(http.MethodPost, "/transfer", body)
			req.Header.Set("Authorization", alice.Token)
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)

			It("gives 400", func() {
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("destination user is valid and balance is sufficient", func() {
			body := strings.NewReader(`{"amount": 10000}`)
			req := httptest.NewRequest(http.MethodPost, "/balance_topup", body)
			req.Header.Set("Authorization", alice.Token)
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)

			body = strings.NewReader(`{"to_username": "bob", "amount": 100}`)
			req = httptest.NewRequest(http.MethodPost, "/transfer", body)
			req.Header.Set("Authorization", alice.Token)
			resp = httptest.NewRecorder()
			h.ServeHTTP(resp, req)

			It("gives 204", func() {
				Expect(resp.Code).To(Equal(http.StatusNoContent))
			})

			It("decreases the alice balance", func() {
				req := httptest.NewRequest(http.MethodGet, "/balance_read", nil)
				req.Header.Set("Authorization", alice.Token)
				resp := httptest.NewRecorder()

				h.ServeHTTP(resp, req)
				var respBody struct {
					Balance int `json:"balance"`
				}

				_ = json.NewDecoder(resp.Body).Decode(&respBody)

				Expect(respBody.Balance).To(Equal(9900))
			})

			It("increases the bob balance", func() {
				req := httptest.NewRequest(http.MethodGet, "/balance_read", nil)
				req.Header.Set("Authorization", bob.Token)
				resp := httptest.NewRecorder()
				h.ServeHTTP(resp, req)

				var respBody struct {
					Balance int `json:"balance"`
				}

				_ = json.NewDecoder(resp.Body).Decode(&respBody)

				Expect(respBody.Balance).To(Equal(100))
			})

			It("records the transaction to the alice", func() {
				req := httptest.NewRequest(http.MethodGet, "/top_transactions_per_user", nil)
				req.Header.Set("Authorization", alice.Token)
				resp := httptest.NewRecorder()
				h.ServeHTTP(resp, req)

				var respBody []struct {
					Username string `json:"username"`
					Amount   int    `json:"amount"`
				}

				_ = json.NewDecoder(resp.Body).Decode(&respBody)

				Expect(respBody).ToNot(BeNil())
				Expect(respBody).To(HaveLen(1))
				Expect(respBody[0].Username).To(Equal("alice"))
				Expect(respBody[0].Amount).To(Equal(-100))
			})

			It("records the transaction to the bob", func() {
				req := httptest.NewRequest(http.MethodGet, "/top_transactions_per_user", nil)
				req.Header.Set("Authorization", bob.Token)
				resp := httptest.NewRecorder()
				h.ServeHTTP(resp, req)

				var respBody []struct {
					Username string `json:"username"`
					Amount   int    `json:"amount"`
				}

				_ = json.NewDecoder(resp.Body).Decode(&respBody)

				Expect(respBody).ToNot(BeNil())
				Expect(len(respBody)).To(Equal(1))
				Expect(respBody[0].Username).To(Equal("bob"))
				Expect(respBody[0].Amount).To(Equal(100))
			})

			It("records the debit transaction", func() {
				req := httptest.NewRequest(http.MethodGet, "/top_users", nil)
				req.Header.Set("Authorization", alice.Token)
				resp := httptest.NewRecorder()
				h.ServeHTTP(resp, req)

				var respBody []struct {
					Username        string `json:"username"`
					TransactedValue int    `json:"transacted_value"`
				}

				_ = json.NewDecoder(resp.Body).Decode(&respBody)

				Expect(respBody).ToNot(BeNil())
				Expect(len(respBody)).To(Equal(1))
				Expect(respBody[0].Username).To(Equal("alice"))
				Expect(respBody[0].TransactedValue).To(Equal(100))
			})
		})
	})
})
