package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/LiquidCats/upgrader/configs"
	"github.com/rs/zerolog"
)

type Srv struct {
	http *http.Server
}

func NewServer(cfg configs.AppConfig, router http.Handler) *Srv {
	server := &http.Server{
		Addr:           net.JoinHostPort("0.0.0.0", cfg.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second, // nolint:mnd
		WriteTimeout:   10 * time.Second, // nolint:mnd
		MaxHeaderBytes: 1 << 20,          // nolint:mnd
	}

	return &Srv{
		http: server,
	}
}

func (s *Srv) Run(ctx context.Context) error {
	logger := zerolog.Ctx(ctx).With().Any("addr", s.http.Addr).Logger()

	go func() {
		logger.Info().Msg("server: starting server")

		if err := s.http.ListenAndServe(); nil != err && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("app: cant start server")
		}
	}()

	<-ctx.Done()

	logger.Info().Msg("server: stopping server")

	if err := s.http.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("server: server shutdown failed")
	}

	return ctx.Err()
}
