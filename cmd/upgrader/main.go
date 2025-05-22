package main

import (
	"context"
	"os"

	"github.com/LiquidCats/graceful"
	"github.com/LiquidCats/upgrader/configs"
	"github.com/LiquidCats/upgrader/internal/adapter/http"
	"github.com/LiquidCats/upgrader/internal/adapter/http/handlers"
	"github.com/LiquidCats/upgrader/internal/adapter/metrics/prometheus"
	"github.com/LiquidCats/upgrader/internal/app/services"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	_ "go.uber.org/automaxprocs"
)

const app = "upgrader"

func main() {
	logger := zerolog.New(os.Stdout).With().Caller().Timestamp().Stack().Logger()

	zerolog.DefaultContextLogger = &logger // nolint:reassign
	ctx := logger.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg, err := configs.Load(app)
	if err != nil {
		logger.Fatal().Err(err).Msg("app: failed to load configuration")
	}

	zerolog.SetGlobalLevel(cfg.App.LogLevel)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: string(cfg.Redis.Password),
		DB:       cfg.Redis.DB,
	})

	metrics := prometheus.NewMetrics(app)

	rootHandler := handlers.NewRootHandler()
	apiHandler := handlers.NewAPIHandler(cfg.Workers)

	router := http.NewRouter()

	router.Any("/", rootHandler.Handle)

	v1Group := router.Group("/v1")
	v1Group.Any("/", apiHandler.Handle)

	runners := []graceful.Runner{
		graceful.Signals,
		graceful.ServerRunner(router, cfg.HTTP),
		prometheus.GerHandler(),
	}

	for _, workerCfg := range cfg.Workers {
		pubSub := redisClient.Subscribe(ctx, workerCfg.FromTopic)

		service := services.NewWebSocketService(workerCfg, pubSub)
		handler := handlers.NewWsHandler(metrics, service)

		v1Group.Any(workerCfg.ToWebsocket, handler.Handle)

		runners = append(runners, service.SubscribeIncomingMessages, service.SubscribeOutgoingMessages)
	}

	logger.Info().Msg("starting up")

	if err = graceful.WaitContext(ctx, runners...); err != nil {
		logger.Fatal().Err(err).Msg("service terminated")
	}

	logger.Info().Msg("service shutdown")
}
