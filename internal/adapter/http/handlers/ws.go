package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/LiquidCats/upgrader/internal/adapter/http/dto"
	"github.com/LiquidCats/upgrader/internal/app/port/exporter"
	"github.com/LiquidCats/upgrader/internal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/go-faster/errors"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout:  30 * time.Second,
	EnableCompression: true,
}

type WsHandler struct {
	mu sync.RWMutex

	metrics exporter.ConnectedClientsMetrics
	srv     *services.WebSocketService
}

func NewWsHandler(metrics exporter.ConnectedClientsMetrics, srv *services.WebSocketService) *WsHandler {
	return &WsHandler{metrics: metrics, srv: srv}
}

func (h *WsHandler) Handle(c *gin.Context) {
	logger := zerolog.Ctx(c).With().Str("websocket", c.Request.URL.Path).Logger()

	logger.Debug().Msg("handling websocket connection")
	defer logger.Debug().Msg("stopped handling websocket connection")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusExpectationFailed,
			dto.ErrorResponse(
				errors.Wrap(err, "connection cant be upgraded to websocket")))
		return
	}

	logger = logger.With().Any("client_addr", conn.RemoteAddr()).Logger()

	h.srv.AddClient(conn)
	h.metrics.ConnectedClientInc(c.Request.URL.Path)
	logger.Info().Msg("new client connected")

	defer func() {
		_ = conn.Close()

		h.srv.RemoveClient(conn)
		h.metrics.ConnectedClientDec(c.Request.URL.Path)
		logger.Info().Msg("client disconnected")
	}()
	//
	//if err = conn.SetReadDeadline(time.Now().Add(15 * time.Second)); err != nil {
	//	return errors.Wrap(err, "could not set read deadline")
	//}
	//if err = conn.SetWriteDeadline(time.Now().Add(15 * time.Second)); err != nil {
	//	return errors.Wrap(err, "failed to set connection deadline")
	//}

	for {
		if _, _, err = conn.NextReader(); err != nil {
			break
		}
	}
}
