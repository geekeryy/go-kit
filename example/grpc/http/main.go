package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./example/grpc/http/"))
	http.Handle("/", http.StripPrefix("/", fs))
	log.Println("Running...")
	http.ListenAndServe("localhost:8081", nil)
}
