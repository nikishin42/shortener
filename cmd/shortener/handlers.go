package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func (a *application) homepage(w http.ResponseWriter, r *http.Request) {
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
	if shortURL, ok := a.cache.ToShort[fullURL]; ok {
		log.Printf("ID for %s found: %s", fullURL, shortURL)
		w.Write([]byte(shortURL))
		return
	}
	hash, err := a.hasher.Write(bodyData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shortURL, err := a.sc.Encode([]uint64{uint64(hash)})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shortURL = "http://localhost:8080/" + shortURL
	if collision, ok := a.cache.ToFull[shortURL]; ok {
		log.Printf("found collision: old URL %s, new URL %s, short URL %s", collision, fullURL, shortURL)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("ID for %s created: %s", fullURL, shortURL)
	a.cache.ToShort[fullURL] = shortURL
	a.cache.ToFull[shortURL] = fullURL
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (a *application) redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := "http://localhost:8080" + r.URL.String()
	if fullURL, ok := a.cache.ToFull[shortURL]; ok {
		log.Printf("URL for %s found: %s", shortURL, fullURL)
		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	log.Printf("URL for %s not found", shortURL)
	w.WriteHeader(400)
}
