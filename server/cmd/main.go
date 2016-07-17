package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	logger.Log("INFO", "Handling requests on port 8080")
	http.HandleFunc("/", base_route)
	http.ListenAndServe(":8080", nil)
}

func base_route(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
