package usecase

import (
	"context"
	"net/http"
	"sync"

	"github.com/LiquidCats/upgrader/configs"
	"github.com/LiquidCats/upgrader/internal/app/domain/entities"
	"github.com/LiquidCats/upgrader/internal/app/port/bus"
	"github.com/LiquidCats/upgrader/internal/app/port/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type Worker struct {
	mu sync.RWMutex

	wsClients map[*websocket.Conn]bool

	relay chan *entities.MessagePayload

	cfg        configs.WorkerConfig
	subscriber bus.Subscriber
	router     gin.IRouter
	upgrader   ws.Upgrader
}

func NewWorker(cfg configs.WorkerConfig, upgrader ws.Upgrader, subscriber bus.Subscriber, router gin.IRouter) *Worker {
	return &Worker{
		cfg:        cfg,
		subscriber: subscriber,
		router:     router,
		upgrader:   upgrader,

		relay: make(chan *entities.MessagePayload, 256),

		wsClients: make(map[*websocket.Conn]bool),
	}
}

func (uc *Worker) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return uc.subscribeIncomingMessages(ctx)
	})

	g.Go(func() error {
		return uc.subscribeOutgoingMessages(ctx)
	})

	uc.router.GET(uc.cfg.ToWebsocket, func(c *gin.Context) {
		w, r := c.Writer, c.Request

		conn, err := uc.upgrader.Upgrade(w, r, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer conn.Close()

		uc.mu.Lock()
		uc.wsClients[conn] = true
		uc.mu.Unlock()

		// Listen for client closure
		defer func() {
			uc.mu.Lock()
			delete(uc.wsClients, conn)
			uc.mu.Unlock()
		}()

		for {
			if _, _, err = conn.NextReader(); err != nil {
				break
			}
		}
	})

	return g.Wait()
}

func (uc *Worker) subscribeOutgoingMessages(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	defer close(uc.relay)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-uc.relay:
			if !ok {
				logger.Info().Msg("relay channel closed")

				return nil
			}

			uc.handleOutgoingMessage(logger, msg)
		}
	}
}

func (uc *Worker) subscribeIncomingMessages(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	defer func() {
		_ = uc.subscriber.Close()
	}()
	defer func() {
		_ = uc.subscriber.Unsubscribe(ctx, uc.cfg.FromTopic)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-uc.subscriber.Channel():
			if !ok {
				logger.Info().Msg("subscriber channel closed")
				return nil
			}

			uc.handleIncomingMessage(msg)
		}
	}
}

func (uc *Worker) handleIncomingMessage(msg *redis.Message) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	uc.relay <- entities.NewMessagePayloadFrom(msg)
}

func (uc *Worker) handleOutgoingMessage(logger *zerolog.Logger, msg *entities.MessagePayload) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	for conn, b := range uc.wsClients {
		if !b {
			conn.Close()
			delete(uc.wsClients, conn)

			logger.Info().Any("remote_addr", conn.RemoteAddr()).Msg("client disconnected")

			return
		}

		if err := conn.WriteMessage(websocket.TextMessage, msg.Bytes()); err != nil {
			conn.Close()
			delete(uc.wsClients, conn)

			logger.Error().Err(err).Any("remote_addr", conn.RemoteAddr()).Msg("failed to send message to client")
			continue
		}
	}
}
