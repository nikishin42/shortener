package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nikishin42/shortener/cmd/shortener/pkg/shortener"
	"github.com/nikishin42/shortener/cmd/shortener/pkg/storage"
)

type Server struct {
	Storage   storage.StorageI
	Shortener shortener.ShortenerI
	Router    *mux.Router
}

func New(storage storage.StorageI, Shortener shortener.ShortenerI) *Server {
	app := &Server{
		Storage:   storage,
		Shortener: Shortener,
		Router:    mux.NewRouter(),
	}
	app.Router.HandleFunc("/", app.Homepage).Methods(http.MethodPost)
	app.Router.HandleFunc("/{id}", app.Redirect).Methods(http.MethodGet)
	return app
}

func (a *Server) Start() {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}