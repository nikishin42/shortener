package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nikishin42/shortener/cmd/shortener/pkg/abbreviator"
	"github.com/nikishin42/shortener/cmd/shortener/pkg/storage"
)

type Server struct {
	Storage     storage.StorageI
	Abbreviator abbreviator.AbbreviatorI
	Router      *mux.Router
}

func New(storage storage.StorageI, abbreviator abbreviator.AbbreviatorI) *Server {
	app := &Server{
		Storage:     storage,
		Abbreviator: abbreviator,
		Router:      mux.NewRouter(),
	}
	app.Router.HandleFunc("/", app.Homepage).Methods(http.MethodPost)
	app.Router.HandleFunc("/{id}", app.Redirect).Methods(http.MethodGet)
	return app
}

func (a *Server) Start() {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}
