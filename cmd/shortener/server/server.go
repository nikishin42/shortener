package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nikishin42/shortener/cmd/shortener/config"
	"github.com/nikishin42/shortener/cmd/shortener/pkg/abbreviator"
	"github.com/nikishin42/shortener/cmd/shortener/pkg/storage"
)

type Server struct {
	Config      *config.Config
	Storage     storage.StorageI
	Abbreviator abbreviator.AbbreviatorI
	Router      *mux.Router
}

func New(config *config.Config, storage storage.StorageI, abbreviator abbreviator.AbbreviatorI) *Server {
	app := &Server{
		Config:      config,
		Storage:     storage,
		Abbreviator: abbreviator,
		Router:      mux.NewRouter(),
	}
	app.Router.HandleFunc("/", app.Homepage).Methods(http.MethodPost)
	app.Router.HandleFunc("/{id}", app.Redirect).Methods(http.MethodGet)
	return app
}

func (s *Server) Start() {
	log.Fatal(http.ListenAndServe(s.Config.Address, s.Router))
}
