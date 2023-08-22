package config

import (
	"flag"
	"log"

	"github.com/nikishin42/shortener/cmd/shortener/constants"
)

type Config struct {
	Address              string
	BaseShortenerAddress string
}

func New() *Config {
	address := flag.String("a", constants.HostPrefix, "HTTP server start address")
	baseShortenerAddress := flag.String("b", constants.HTTPHostPrefix, "base address of the resulting shortened URL")
	flag.Parse()
	log.Printf("flags parsed: -a: %s, -b: %s", *address, *baseShortenerAddress)
	return &Config{
		Address:              *address,
		BaseShortenerAddress: *baseShortenerAddress,
	}
}
