# crypt

A small go microservice that provides functionality to store data in an encrypted format, along with future retrieval of such data.

## Go-kit

This service uses go-kit for its design, please have a look at the following links for more information:

-[Microservices in Go using Go-kit] (https://blog.heroku.com/microservices_in_go_using_go_kit)

-[Go kit: Go in the modern enterprise] (https://peter.bourgon.org/go-kit/)

-[Go kit FAQ] (http://gokit.io/faq/)

## Building and running the server

    mkdir -p $GOPATH/src/github.com/alexnvdias/crypt
    cd $GOPATH/src/github.com/alexnvdias/crypt/server
    make

After running the above, you should have a `server` binary, which you can then run.

## Performing requests

Using cURL:

`curl -X POST -d '{"id": "data_id", "plaintext": "data_to_encrypt"}' http://localhost:8080/store` to insert some plaintext data;

`curl 'http://localhost:8080/retrieve?id=data_id&key=retrieval_key'` to retrieve the originally inserted data using the resulting key.
Please note that the retrieval key will need to be URL encoded as it is returned in base64 format.

## Running the client

The client can be used to perform requests to the server, like so:

    package main
    
    import (
        "github.com/alexnvdias/crypt/client"
    )
    
    func main() {
        c := client.NewHTTPClient("http://localhost:8080/")
        key, err_store := c.Store([]byte("1"), []byte("data_1"))
        if err_store != nil {
            panic(err_store)
        }
        data, err_retrieve := c.Retrieve([]byte("1"), key)
        if err_retrieve != nil {
            panic(err_retrieve)
        }
        println(string(data))
    }
