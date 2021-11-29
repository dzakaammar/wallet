package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wallet/internal"
	"wallet/internal/account"
	"wallet/internal/inmemory"
	"wallet/internal/rest"
	"wallet/internal/server"
	"wallet/internal/transaction"
	"wallet/internal/user"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	RunE: runHTTPServer,
}

func init() {
	// pass all the env configs through flags, just for the sake of simplicity
	rootCmd.PersistentFlags().IntP("port", "p", 8080, "http server port")
	rootCmd.PersistentFlags().StringP("secret", "s", "secret", "token secret")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runHTTPServer(cmd *cobra.Command, _ []string) error {
	port, _ := cmd.PersistentFlags().GetInt("port")
	secret, _ := cmd.PersistentFlags().GetString("secret")

	// repository
	userRepo := inmemory.NewUserRepository()
	accountRepo := inmemory.NewAccountRepository()
	trxRepo := inmemory.NewTransactionRepository()

	accountLocker := inmemory.NewAccountLock()

	// service
	trxSvc := transaction.NewService(trxRepo, userRepo)
	accountSvc := account.NewService(accountRepo, userRepo, account.NewEventHandler(trxSvc), accountLocker)
	userSvc := user.NewService(userRepo, user.NewEventHandler(accountSvc))

	// http handler
	handler := rest.NewHandler(accountSvc, userSvc, trxSvc, internal.NewAuthToken(secret))

	// http server
	srv := server.NewMuxHTTPServer(handler, port)

	go func() {
		//- start service
		if err := srv.Start(); err != nil {
			fmt.Println(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Stop(ctx)

	return nil
}
