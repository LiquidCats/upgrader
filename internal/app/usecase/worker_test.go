package usecase_test

import (
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/LiquidCats/upgrader/configs"
	"github.com/LiquidCats/upgrader/internal/adapter/http/server"
	"github.com/LiquidCats/upgrader/internal/adapter/ws/upgrader"
	"github.com/LiquidCats/upgrader/internal/app/usecase"
	"github.com/LiquidCats/upgrader/test/mocks"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWorker_Run(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()

	cfg := configs.WorkerConfig{
		FromTopic:   "from_test",
		ToWebsocket: "to_test",
	}

	upg := upgrader.NewUpgrader()
	sub := mocks.NewMockSubscriber(t)

	rt := server.NewRouter()
	srv := httptest.NewServer(rt)
	defer srv.Close()

	ch := make(chan *redis.Message)

	sub.On("Channel", mock.Anything).Return(func() <-chan *redis.Message { return ch }())

	w := usecase.NewWorker(cfg, upg, sub, rt)

	go w.Run(ctx)

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/" + cfg.ToWebsocket

	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	require.NoError(t, err)
	defer client.Close()

	ch <- &redis.Message{
		Payload: `{"foo":"bar"}`,
	}

	msgType, reader, err := client.NextReader()
	require.NoError(t, err)
	msg, err := io.ReadAll(reader)
	require.NoError(t, err)

	assert.Equal(t, websocket.TextMessage, msgType)
	assert.Equal(t, `{"foo":"bar"}`, string(msg))
}
