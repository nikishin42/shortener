package server

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func (a *Server) Homepage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("wrong method:", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.URL.String() != "localhost:8080" && r.URL.String() != "localhost:8080/" {
		log.Println("query not empty:", r.URL.String())
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
	if id, ok := a.Storage.GetID(fullURL); ok {
		log.Printf("ID for %s found: %s", fullURL, id)
		w.Write([]byte(id))
		return
	}
	id, err := a.Abbreviator.CreateID(bodyData)
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
	if r.Method != http.MethodGet {
		log.Println("wrong method:", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.URL.String() == "localhost:8080" || r.URL.String() == "localhost:8080/" {
		log.Println("empty query")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := r.RequestURI
	fullURL, ok := a.Storage.GetFullURL(id)
	if !ok {
		log.Printf("URL for %s not found", id)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("URL for %s found: %s", id, fullURL)
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
