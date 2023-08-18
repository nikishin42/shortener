package abbreviator

import (
	"crypto/md5"
	"hash"
	"log"

	"github.com/sqids/sqids-go"
)

//go:generate mockgen --build_flags=--mod=mod -package=abbreviator -destination=abbreviator_mock.go . AbbreviatorI
type AbbreviatorI interface {
	CreateID(data []byte) (string, error)
}

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

func (s *Abbreviator) CreateID(data []byte) (string, error) {
	hash, err := s.hash.Write(data)
	if err != nil {
		return "", err
	}
	shortURL, err := s.sc.Encode([]uint64{uint64(hash)})
	if err != nil {
		return "", err
	}
	shortURL = "localhost:8080/" + shortURL
	return shortURL, nil
}
