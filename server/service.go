package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

func Run(httpAddr string, logger log.Logger) {

	var ctx context.Context
	{
		ctx = context.Background()
	}

	var st Storage
	{
		st = NewInMemoryStorage()
	}

	s := NewStoreService(st)

	var h http.Handler
	{
		h = SetUpHTTPHandlers(ctx, s, log.NewContext(logger).With("component", "HTTP"))
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
