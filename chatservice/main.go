package main

import (
	"log"
	"net/http"

	chat "./Service"
	"github.com/go-chi/chi"
)

func main() {
	router := chat.Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	log.Fatal(http.ListenAndServe(":8080", router)) // Note, the port is usually gotten from the environment.
}
