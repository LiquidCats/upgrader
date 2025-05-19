package bus

import (
	"context"
	"io"

	"github.com/redis/go-redis/v9"
)

type Subscriber interface {
	Channel(opts ...redis.ChannelOption) <-chan *redis.Message
	Unsubscribe(ctx context.Context, channels ...string) error
	io.Closer
}
