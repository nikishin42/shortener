package main

import (
	"github.com/nikishin42/shortener/cmd/shortener/config"
	"github.com/nikishin42/shortener/cmd/shortener/pkg/abbreviator"
	"github.com/nikishin42/shortener/cmd/shortener/pkg/storage"
	"github.com/nikishin42/shortener/cmd/shortener/servicelayer"
)

func main() {
	cfg := config.New()
	servicelayer.New(cfg, storage.New(*cfg), abbreviator.New()).Start()
}
