package main

import (
	"log"
	"net/http"
)

func main() {
	app := newApplication()
	app.router.HandleFunc("/", app.homepage).Methods(http.MethodPost)
	app.router.HandleFunc("/{id}", app.redirect).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", app.router))
}
