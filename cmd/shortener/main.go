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
	fullUrl := string(bodyData)
	_, err = url.ParseRequestURI(fullUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if shortUrl, ok := cache[fullUrl]; ok {
		w.Write([]byte(shortUrl))
		return
	}
	hash, err := hasher.Write(bodyData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortUrl, err := sc.Encode([]uint64{uint64(hash)})
	if _, ok := cacheBack[shortUrl]; ok || err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shortUrl = "http://localhost:8080/" + shortUrl
	cache[fullUrl] = shortUrl
	cacheBack[shortUrl] = fullUrl
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortUrl))
}

func apiPage(w http.ResponseWriter, r *http.Request) {
	shortUrl := "http://localhost:8080" + r.URL.String()
	if fullUrl, ok := cacheBack[shortUrl]; ok {
		w.Header().Set("Location", fullUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	log.Println(shortUrl)
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
