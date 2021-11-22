package rest_test

import (
	"net/http"
	"testing"
	"wallet/internal"
	"wallet/internal/account"
	"wallet/internal/inmemory"
	"wallet/internal/rest"
	"wallet/internal/transaction"
	"wallet/internal/user"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var h http.Handler

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store Suite")
}

func InitHandler() {
	userRepo := inmemory.NewUserRepository()
	accountRepo := inmemory.NewAccountRepository()
	trxRepo := inmemory.NewTransactionRepository()

	// service
	trxSvc := transaction.NewService(trxRepo, userRepo)
	accountSvc := account.NewService(accountRepo, userRepo, account.NewEventHandler(trxSvc))
	userSvc := user.NewService(userRepo, user.NewEventHandler(accountSvc))

	hd := rest.NewHandler(accountSvc, userSvc, trxSvc, internal.NewAuthToken("test"))
	r := mux.NewRouter()
	hd.Register(r)

	h = r
}

func Teardown() {
	h = nil
}
