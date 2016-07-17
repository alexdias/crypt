# crypt

A small go microservice that provides functionality to store data in an encrypted format, along with future retrieval of such data.

## Go-kit

This service uses go-kit for its design, please have a look at the following links for more information:

-[Microservices in Go using Go-kit] (https://blog.heroku.com/microservices_in_go_using_go_kit)

-[Go kit: Go in the modern enterprise] (https://peter.bourgon.org/go-kit/)

-[Go kit FAQ] (http://gokit.io/faq/)

## Building and running the server

`mkdir -p $GOPATH/src/github.com/alexnvdias/crypt`

`cd $GOPATH/src/github.com/alexnvdias/crypt/server`

`make`

After running the above, you should have a `server` binary, which you can then run.

## Performing requests

Using cURL:

`curl -X POST -d '{"id": "data_id", "plaintext": "data_to_encrypt"}' http://localhost:8080/` to insert some plaintext data;

`curl 'http://localhost:8080?id=data_id&key=retrieval_key'` to retrieve the originally inserted data using the resulting key.
