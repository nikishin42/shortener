package businesslayer

import (
	"errors"
	"fmt"

	"github.com/nikishin42/shortener/cmd/shortener/interfaces"
)

func GetFullAddress(baseShortAddress, id string, storage interfaces.GetterFullURL) (string, error) {
	shortURL := baseShortAddress + "/" + id
	fullURL, ok := storage.GetFullURL(shortURL)
	if !ok {
		return "", errors.New(fmt.Sprintf("full URL for %s not found", shortURL))
	}
	return fullURL, nil
}
