package handlers

import (
	"net/http"
	"time"

	"github.com/LiquidCats/upgrader/internal/adapter/http/dto"
	"github.com/LiquidCats/upgrader/internal/app/port/exporter"
	"github.com/LiquidCats/upgrader/internal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/go-faster/errors"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var upgrader = websocket.Upgrader{ //nolint:gochecknoglobals
	HandshakeTimeout:  30 * time.Second, //nolint:mnd
	EnableCompression: true,
}

type WsHandler struct {
	metrics exporter.ConnectedClientsMetrics
	srv     *services.WebSocketService
}

func NewWsHandler(metrics exporter.ConnectedClientsMetrics, srv *services.WebSocketService) *WsHandler {
	return &WsHandler{metrics: metrics, srv: srv}
}

func (h *WsHandler) Handle(c *gin.Context) {
	logger := zerolog.Ctx(c.Request.Context()).With().Str("websocket", c.Request.URL.Path).Logger()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusExpectationFailed,
			dto.ErrorResponse(
				errors.Wrap(err, "connection cant be upgraded to websocket"),
			),
		)
		return
	}
	defer func() { _ = conn.Close() }()

	logger = logger.With().Any("client_addr", conn.RemoteAddr()).Logger()

	h.srv.AddClient(conn)
	defer h.srv.RemoveClient(conn)

	h.metrics.ConnectedClientInc(c.Request.URL.Path)
	defer h.metrics.ConnectedClientDec(c.Request.URL.Path)

	logger.Info().Msg("new client connected")
	defer logger.Info().Msg("client disconnected")

	for {
		if _, _, err = conn.NextReader(); err != nil {
			break
		}
	}
}
