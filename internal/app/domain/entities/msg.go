package entities

import "github.com/redis/go-redis/v9"

type MessagePayload struct {
	V string
}

func (m *MessagePayload) Bytes() []byte {
	return []byte(m.V)
}

func NewMessagePayloadFrom(message *redis.Message) *MessagePayload {
	return &MessagePayload{
		V: message.Payload,
	}
}
