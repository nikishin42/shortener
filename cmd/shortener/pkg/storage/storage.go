package storage

import (
	"errors"
	"fmt"
)

//go:generate mockgen --build_flags=--mod=mod -package=storage -destination=storage_mock.go . StorageI
type StorageI interface {
	GetID(fullURL string) (string, bool)
	GetFullURL(shortURL string) (string, bool)
	SetPair(shortURL, fullURL string) error
}

type Storage struct {
	toShort map[string]string
	toFull  map[string]string
}

func (c *Storage) GetID(fullURL string) (string, bool) {
	shortURL, ok := c.toShort[fullURL]
	return shortURL, ok
}

func (c *Storage) GetFullURL(id string) (string, bool) {
	fullURL, ok := c.toFull[id]
	return fullURL, ok
}

func (c *Storage) SetPair(id, fullURL string) error {
	if collision, ok := c.GetFullURL(id); ok {
		return errors.New(fmt.Sprintf("found collision: old URL %s, new URL %s, ID %s", collision, fullURL, id))
	}
	c.toShort[fullURL] = id
	c.toFull[id] = fullURL
	return nil
}

func New() *Storage {
	return &Storage{
		toShort: make(map[string]string),
		toFull:  make(map[string]string),
	}
}
