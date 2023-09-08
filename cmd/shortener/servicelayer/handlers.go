package servicelayer

import (
	"io"
	"net/http"
	"net/url"

	"github.com/nikishin42/shortener/cmd/shortener/businesslayer"
)

func (s *Server) Homepage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if len(r.URL.Path) > 1 {
		s.Logger.Infof("path not empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		s.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fullURL := string(bodyData)
	_, err = url.ParseRequestURI(fullURL)
	if err != nil {
		s.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, fromCache, err := businesslayer.GetOrCreateID(s.Storage, s.Abbreviator, bodyData, s.Config.BaseShortenerAddress)
	if err != nil {
		s.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if fromCache {
		s.Logger.Infof("ID for %s found: %s", fullURL, id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(id))
		return
	}
	s.Logger.Infof("ID for %s created: %s", fullURL, id)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func (s *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) < 2 {
		s.Logger.Infof("empty path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := r.URL.Path[1:]
	fullURL, err := businesslayer.GetFullAddress(s.Config.BaseShortenerAddress, id, s.Storage)
	if err != nil {
		s.Logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.Logger.Infof("URL for %s found: %s", id, fullURL)
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}