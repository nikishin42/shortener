package servicelayer

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nikishin42/shortener/cmd/shortener/config"
	"github.com/nikishin42/shortener/cmd/shortener/interfaces"
)

type Server struct {
	Config      *config.Config
	Storage     interfaces.Storage
	Abbreviator interfaces.CreatorID
	Router      *mux.Router
}

func New(config *config.Config, storage interfaces.Storage, abbreviator interfaces.CreatorID) *Server {
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
	err := http.ListenAndServe(s.Config.Address, s.Router)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Print("Error: ", err)
	}
}
