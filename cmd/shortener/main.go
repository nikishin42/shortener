package main

import (
	"github.com/nikishin42/shortener/cmd/shortener/pkg/abbreviator"
	"github.com/nikishin42/shortener/cmd/shortener/pkg/storage"
	"github.com/nikishin42/shortener/cmd/shortener/server"
)

func main() {
	server.New(storage.New(), abbreviator.New()).Start()
}
