package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Config struct {
	Address string        `yaml:"address" env:"LISTEN_ADDRESS" env-default:"localhost:8080"`
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
}

type Server struct {
	server *http.Server
}

func NewServer(cfg Config, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:         cfg.Address,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  30 * time.Second,
			Handler:      handler,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = s.server.Shutdown(ctxShutdown)
	}()

	err := s.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
