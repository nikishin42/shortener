package servicelayer

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/nikishin42/shortener/cmd/shortener/businesslayer"
)

func (s *Server) Homepage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if len(r.URL.Path) > 1 {
		log.Print("query not empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fullURL := string(bodyData)
	_, err = url.ParseRequestURI(fullURL)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, fromCache, err := businesslayer.CreateID(s.Storage, s.Abbreviator, bodyData, s.Config.BaseShortenerAddress)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if fromCache {
		log.Printf("ID for %s found: %s", fullURL, id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(id))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func (s *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) < 2 {
		log.Print("empty query")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := r.URL.Path[1:]
	fullURL, err := businesslayer.GetFullAddress(s.Config.BaseShortenerAddress, id, s.Storage)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("URL for %s found: %s", id, fullURL)
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
