package main

import (
	"crypto/md5"
	"hash"
	"log"

	"github.com/gorilla/mux"
	"github.com/sqids/sqids-go"
)

type application struct {
	cache  cache
	sc     *sqids.Sqids
	hasher hash.Hash
	router *mux.Router
}

type cache struct {
	ToShort map[string]string
	ToFull  map[string]string
}

func newApplication() *application {
	sc, err := sqids.New()
	if err != nil {
		log.Fatal(err)
	}
	return &application{
		cache: cache{
			ToShort: make(map[string]string),
			ToFull:  make(map[string]string),
		},
		sc:     sc,
		hasher: md5.New(),
		router: mux.NewRouter(),
	}
}
