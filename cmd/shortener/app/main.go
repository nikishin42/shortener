package main

import (
	"github.com/nikishin42/shortener/cmd/shortener/server"
)

func main() {
	server.New().Start()
}
