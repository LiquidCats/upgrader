package main

import (
	"context"
	"os"

	"github.com/LiquidCats/graceful"
	"github.com/LiquidCats/upgrader/configs"
	"github.com/LiquidCats/upgrader/internal/adapter/http/server"
	"github.com/LiquidCats/upgrader/internal/adapter/ws/upgrader"
	"github.com/LiquidCats/upgrader/internal/app/usecase"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
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
		logger.Error().Err(err).Msg("app: failed to load configuration")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: string(cfg.Redis.Password),
		DB:       cfg.Redis.DB,
	})

	router := server.NewRouter()
	wsUpgrader := upgrader.NewUpgrader()

	srv := server.NewServer(cfg.App, router)

	v1Group := router.Group("v1")

	runners := []graceful.Runner{
		graceful.Signals,
		srv.Run,
	}

	for _, workerCfg := range cfg.Workers {
		pubSub := redisClient.Subscribe(ctx, workerCfg.FromTopic)

		worker := usecase.NewWorker(workerCfg, wsUpgrader, pubSub, v1Group)

		runners = append(runners, worker.Run)
	}

	logger.Info().Msg("starting up")

	if err = graceful.WaitContext(ctx, runners...); err != nil {
		logger.Fatal().Err(err).Msg("service terminated")
	}

	logger.Info().Msg("service shutdown")
}
