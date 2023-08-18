package server

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func (a *Server) Homepage(w http.ResponseWriter, r *http.Request) {
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fullURL := string(bodyData)
	_, err = url.ParseRequestURI(fullURL)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if id, ok := a.Storage.GetID(fullURL); ok {
		log.Printf("ID for %s found: %s", fullURL, id)
		w.Write([]byte(id))
		return
	}
	id, err := a.Shortener.CreateID(bodyData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = a.Storage.SetPair(id, fullURL)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("ID for %s created: %s", fullURL, id)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func (a *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	id := "http://localhost:8080" + r.URL.String()
	if fullURL, ok := a.Storage.GetFullURL(id); ok {
		log.Printf("URL for %s found: %s", id, fullURL)
		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	log.Printf("URL for %s not found", id)
	w.WriteHeader(400)
}
