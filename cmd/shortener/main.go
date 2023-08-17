package main

import (
	"crypto/md5"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/sqids/sqids-go"
)

var (
	cache     = make(map[string]string)
	cacheBack = make(map[string]string)
	sc        *sqids.Sqids
	hasher    = md5.New()
)

func home(w http.ResponseWriter, r *http.Request) {
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fullURL := string(bodyData)
	_, err = url.ParseRequestURI(fullURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if shortURL, ok := cache[fullURL]; ok {
		w.Write([]byte(shortURL))
		return
	}
	hash, err := hasher.Write(bodyData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortURL, err := sc.Encode([]uint64{uint64(hash)})
	if _, ok := cacheBack[shortURL]; ok || err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shortURL = "http://localhost:8080/" + shortURL
	cache[fullURL] = shortURL
	cacheBack[shortURL] = fullURL
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func apiPage(w http.ResponseWriter, r *http.Request) {
	shortURL := "http://localhost:8080" + r.URL.String()
	if fullURL, ok := cacheBack[shortURL]; ok {
		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	log.Println(shortURL)
	w.WriteHeader(400)
}

func main() {
	var err error
	sc, err = sqids.New()
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods(http.MethodPost)
	r.HandleFunc("/{id}", apiPage).Methods(http.MethodGet)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
