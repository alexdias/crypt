.PHONY: default
default: build

.PHONY: get
get:
	go get github.com/go-kit/kit/log
	go get github.com/gorilla/mux

.PHONY: build
build: get
	go build -o server cmd/main.go
