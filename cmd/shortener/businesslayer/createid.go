package businesslayer

import (
	"log"

	"github.com/nikishin42/shortener/cmd/shortener/interfaces"
)

func GetOrCreateID(storage interfaces.GetOrSetID, abbreviator interfaces.CreatorID, bodyData []byte, baseShortAddress string) (string, bool, error) {
	fullURL := string(bodyData)
	if id, ok := storage.GetID(fullURL); ok {
		return id, ok, nil
	}
	id, err := abbreviator.CreateID(bodyData, baseShortAddress)
	if err != nil {
		log.Print(err)
		return "", false, err
	}
	err = storage.SetPair(id, fullURL)
	if err != nil {
		log.Print(err)
		return "", false, err
	}
	log.Printf("ID for %s created: %s", fullURL, id)
	return id, false, nil
}
