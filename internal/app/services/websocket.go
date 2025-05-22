package services

import (
	"context"
	"sync"

	"github.com/LiquidCats/upgrader/configs"
	"github.com/LiquidCats/upgrader/internal/app/domain/entities"
	"github.com/LiquidCats/upgrader/internal/app/port/bus"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type WebSocketService struct {
	mu sync.RWMutex

	wsClients map[*websocket.Conn]bool

	relay chan *entities.MessagePayload

	cfg        configs.WorkerConfig
	subscriber bus.Subscriber
}

func NewWebSocketService(
	cfg configs.WorkerConfig,
	subscriber bus.Subscriber,
) *WebSocketService {
	return &WebSocketService{
		cfg:        cfg,
		subscriber: subscriber,

		relay: make(chan *entities.MessagePayload, 256), //nolint:mnd

		wsClients: make(map[*websocket.Conn]bool),
	}
}

func (srv *WebSocketService) AddClient(ws *websocket.Conn) {
	srv.mu.RLock()
	defer srv.mu.RUnlock()
	srv.wsClients[ws] = true
}

func (srv *WebSocketService) RemoveClient(ws *websocket.Conn) {
	srv.mu.RLock()
	defer srv.mu.RUnlock()
	_ = ws.Close()
	delete(srv.wsClients, ws)
}

func (srv *WebSocketService) clientLen() int {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	return len(srv.wsClients)
}
func (srv *WebSocketService) SubscribeOutgoingMessages(ctx context.Context) error {
	logger := zerolog.Ctx(ctx).With().Str("topic", srv.cfg.FromTopic).Logger()

	logger.Info().Msg("subscribing to outgoing messages")
	defer logger.Info().Msg("stopped subscribing to outgoing messages")
	defer close(srv.relay)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-srv.relay:
			if !ok {
				logger.Info().Msg("relay channel is closed")
				return nil
			}

			srv.handleOutgoingMessage(&logger, msg)
		}
	}
}

func (srv *WebSocketService) SubscribeIncomingMessages(ctx context.Context) error {
	logger := zerolog.Ctx(ctx).With().Str("topic", srv.cfg.FromTopic).Logger()

	defer func() {
		_ = srv.subscriber.Unsubscribe(ctx, srv.cfg.FromTopic)
		_ = srv.subscriber.Close()
		logger.Info().Msg("stopped listening for incoming messages")
	}()

	logger.Info().Msg("listening for incoming messages")

	msgCh := srv.subscriber.Channel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgCh:
			if !ok {
				logger.Info().Msg("subscriber channel is closed")
				return nil
			}

			if srv.clientLen() > 0 {
				srv.mu.RLock()
				srv.relay <- entities.NewMessagePayloadFrom(msg)
				srv.mu.RUnlock()
			}
		}
	}
}

func (srv *WebSocketService) handleOutgoingMessage(logger *zerolog.Logger, msg *entities.MessagePayload) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	for conn, exists := range srv.wsClients {
		if !exists {
			srv.RemoveClient(conn)

			logger.Info().Any("remote_addr", conn.RemoteAddr()).Msg("client disconnected")

			return
		}

		if err := conn.WriteMessage(websocket.TextMessage, msg.Bytes()); err != nil {
			srv.RemoveClient(conn)

			logger.Error().Err(err).Any("remote_addr", conn.RemoteAddr()).Msg("failed to send message to client")
			continue
		}
	}
}
