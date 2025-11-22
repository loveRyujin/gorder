package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("listening on :8080...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v", r.RequestURI)
		fmt.Fprintln(w, "<h1>Welcome to Home Page<h1>")
	})
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
	log.Println("server shutdown")
}
