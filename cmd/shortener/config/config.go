package config

import (
	"flag"
	"log"
	"os"

	"github.com/nikishin42/shortener/cmd/shortener/constants"
)

type Config struct {
	Address              string
	BaseShortenerAddress string
}

func parseFlags() (string, string) {
	address := flag.String("a", constants.HostPrefix, "HTTP server start address")
	baseShortenerAddress := flag.String("b", constants.HTTPHostPrefix, "base address of the resulting shortened URL")
	flag.Parse()
	log.Printf("flags parsed: -a: %s, -b: %s", *address, *baseShortenerAddress)
	return *address, *baseShortenerAddress
}

func New() *Config {
	a, b := parseFlags()
	address, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		log.Printf("server address got from env: %s", address)
	} else {
		address = a
		log.Printf("server address got from flag: %s", address)
	}
	baseShortenerAddress, ok := os.LookupEnv("BASE_URL")
	if ok {
		log.Printf("base shortener address got from env: %s", baseShortenerAddress)
	} else {
		baseShortenerAddress = b
		log.Printf("base shortener address got from flag: %s", baseShortenerAddress)
	}
	return &Config{
		Address:              address,
		BaseShortenerAddress: baseShortenerAddress,
	}
}
