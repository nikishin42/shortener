package abbreviator

import (
	"crypto/md5"
	"hash"
	"log"

	"github.com/sqids/sqids-go"
)

type Abbreviator struct {
	hash hash.Hash
	sc   *sqids.Sqids
}

func New() *Abbreviator {
	cs, err := sqids.New()
	if err != nil {
		log.Fatalln(err)
	}
	return &Abbreviator{
		hash: md5.New(),
		sc:   cs,
	}
}

func (a *Abbreviator) CreateID(data []byte, base string) (string, error) {
	hash, err := a.hash.Write(data)
	if err != nil {
		return "", err
	}
	shortURL, err := a.sc.Encode([]uint64{uint64(hash)})
	if err != nil {
		return "", err
	}
	shortURL = base + "/" + shortURL
	return shortURL, nil
}
