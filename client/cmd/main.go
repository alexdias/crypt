package main

import (
	"github.com/alexnvdias/crypt/client"
)

func main() {
	c := client.NewHTTPClient()
	c.Store([]byte("1"), []byte("data_1"))
}
