package businesslayer

import (
	"fmt"

	"github.com/nikishin42/shortener/cmd/shortener/interfaces"
)

func GetFullAddress(baseShortAddress, id string, storage interfaces.GetterFullURL) (string, error) {
	shortURL := baseShortAddress + "/" + id
	fullURL, ok := storage.GetFullURL(shortURL)
	if !ok {
		return "", fmt.Errorf("full URL for %s not found", shortURL)
	}
	return fullURL, nil
}
