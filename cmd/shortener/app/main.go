package main

import (
	"github.com/nikishin42/shortener/cmd/shortener/internal/server"
)

func main() {
	server.New().Start()
}
