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
		log.Println("query not empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fullURL := string(bodyData)
	_, err = url.ParseRequestURI(fullURL)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, fromCache, err := businesslayer.CreateID(s.Storage, s.Abbreviator, bodyData, s.Config.BaseShortenerAddress)
	if err != nil {
		log.Println(err)
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
		log.Println("empty query")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := s.Config.BaseShortenerAddress + "/" + r.URL.Path[1:]
	fullURL, ok := s.Storage.GetFullURL(id)
	if !ok {
		log.Printf("URL for %s not found", id)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("URL for %s found: %s", id, fullURL)
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
