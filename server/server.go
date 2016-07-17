package main

import "net/http"

func main() {
	http.HandleFunc("/", base_route)
	http.ListenAndServe(":8080", nil)
}

func base_route(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
