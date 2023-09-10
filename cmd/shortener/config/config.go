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
	FileStoragePath      string
}

func parseFlags() (string, string, string) {
	address := flag.String("a", constants.DefaultHost, "HTTP server start address")
	baseShortenerAddress := flag.String("b", constants.HTTPHostPrefix, "base address of the resulting shortened URL")
	defaultFilePath := flag.String("f", constants.DefaultFilePath, "base file path")
	flag.Parse()
	log.Printf("flags parsed: -a: %s, -b: %s, -f: %s", *address, *baseShortenerAddress, *defaultFilePath)
	return *address, *baseShortenerAddress, *defaultFilePath
}

func New() *Config {
	a, b, f := parseFlags()
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
	baseFilePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		log.Printf("base file path got from env: %s", baseFilePath)
	} else {
		baseFilePath = f
		log.Printf("base file path got from flag: %s", baseFilePath)
	}
	return &Config{
		Address:              address,
		BaseShortenerAddress: baseShortenerAddress,
		FileStoragePath:      baseFilePath,
	}
}
