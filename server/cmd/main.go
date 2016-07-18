package main

import (
	"os"

	"github.com/go-kit/kit/log"

	"github.com/alexnvdias/crypt/server"
)

func main() {
	httpAddr := ":8080"

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}

	server.Run(httpAddr, logger)
}
