package storage

import (
	"fmt"

	"github.com/nikishin42/shortener/cmd/shortener/config"
	"github.com/nikishin42/shortener/cmd/shortener/pkg/filereader"
)

type Storage struct {
	toShort     map[string]string
	toFull      map[string]string
	fileService *filereader.FileReaderWriter
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
		return fmt.Errorf("found collision: old URL %s, new URL %s, ID %s", collision, fullURL, id)
	}
	c.toShort[fullURL] = id
	c.toFull[id] = fullURL
	err := c.fileService.WriteFileEvent(id, fullURL)
	if err != nil {
		return err
	}
	return nil
}

func (c *Storage) setValue() error {
	events, err := c.fileService.ReadFileEvents()
	if err != nil {
		return err
	}
	for _, event := range events {
		c.toShort[event.OriginURL] = event.ShortURL
		c.toFull[event.ShortURL] = event.OriginURL
		if err != nil {
			return err
		}
	}
	return nil
}

func New(cfg config.Config) *Storage {
	writer, err := filereader.New(cfg.FileStoragePath)
	if err != nil {
		panic(err)
	}
	s := &Storage{
		toShort:     make(map[string]string),
		toFull:      make(map[string]string),
		fileService: writer,
	}
	err = s.setValue()
	if err != nil {
		panic(err)
	}
	return s
}
