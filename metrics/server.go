package metrics

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type server struct {
	srv *http.Server
}

func New() *server {
	setup()

	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    ":2112",
		Handler: r,
	}

	return &server{
		srv: srv,
	}
}

func (s *server) Start() {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("failed to start metrics handler")
			os.Exit(1)
		}
	}()
}

func (s *server) Stop() {
	err := s.srv.Shutdown(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("failed to shut down metrics server")
	}
}
