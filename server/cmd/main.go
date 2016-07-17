package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"

	"github.com/alexnvdias/crypt/server"
)

func main() {
	httpAddr := ":8080"

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}

	var ctx context.Context
	{
		ctx = context.Background()
	}

	var s server.StoreService
	{
		s = server.NewInMemoryStoreService()
	}

	var h http.Handler
	{
		h = server.SetUpHTTPHandlers(ctx, s, log.NewContext(logger).With("component", "HTTP"))
	}

	errors := make(chan error)
	go func() {
		channel := make(chan os.Signal)
		signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM) // handle these signals
		errors <- fmt.Errorf("%s", <-channel)
	}()

	go func() {
		logger.Log("INFO", "Listening on", httpAddr)
		errors <- http.ListenAndServe(httpAddr, h)
	}()

	logger.Log("Exiting", <-errors)
}
