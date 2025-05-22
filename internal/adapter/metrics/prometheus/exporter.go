package prometheus

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

type Exporter struct {
	server *http.Server
}

func NewServer() *Exporter {
	mux := gin.New()
	gin.SetMode(gin.ReleaseMode)

	mux.Any("/metrics", createHandler())

	srv := &http.Server{
		Addr:    "0.0.0.0:9100",
		Handler: mux,
	}

	return &Exporter{
		server: srv,
	}
}

func createHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (e *Exporter) Run(ctx context.Context) error {
	logger := zerolog.Ctx(ctx).With().Any("addr", e.server.Addr).Logger()

	go func() {
		logger.Info().Msg("starting metrics exporter")
		if err := e.server.ListenAndServe(); err != nil {
			logger.Fatal().Err(err).Msg("metrics export stopped")
		}
	}()

	<-ctx.Done()

	logger.Info().Msg("shutting down metrics exporter")

	if err := e.server.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("failed to shutdown metrics exporter")
	}

	return ctx.Err()
}
