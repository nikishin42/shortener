package servicelayer

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/nikishin42/shortener/cmd/shortener/config"
	"github.com/nikishin42/shortener/cmd/shortener/interfaces"
)

type Server struct {
	Config      *config.Config
	Storage     interfaces.Storage
	Abbreviator interfaces.CreatorID
	Router      *mux.Router
	Logger      *zap.SugaredLogger
}

func New(config *config.Config, storage interfaces.Storage, abbreviator interfaces.CreatorID) *Server {
	app := &Server{
		Config:      config,
		Storage:     storage,
		Abbreviator: abbreviator,
		Router:      mux.NewRouter(),
		Logger:      zap.NewExample().Sugar(),
	}
	app.Router.Handle("/", app.WithLogging(app.Homepage)).Methods(http.MethodPost)
	app.Router.Handle("/{id}", app.WithLogging(app.Redirect)).Methods(http.MethodGet)
	return app
}

func (s *Server) Start() {
	defer s.Logger.Sync()
	err := http.ListenAndServe(s.Config.Address, s.Router)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Print("Error: ", err)
	}
}
