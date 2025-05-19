package upgrader

import (
	"time"

	"github.com/gorilla/websocket"
)

func NewUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		HandshakeTimeout:  10 * time.Second, //nolint:mnd
		EnableCompression: true,
	}
}
