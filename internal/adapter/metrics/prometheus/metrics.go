package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	connectedClients *prometheus.GaugeVec
	sentMessages     *prometheus.CounterVec
	receivedMessages *prometheus.CounterVec
}

func NewMetrics(name string) *Metrics {
	connectedClients := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: name,
		Name:      "connected_clients_total",
	}, []string{"websocket"})

	sentMessages := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: name,
		Name:      "sent_messages_total",
	}, []string{"websocket"})

	receivedMessages := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: name,
		Name:      "received_messages_total",
	}, []string{"channel"})

	prometheus.MustRegister(
		connectedClients,
		sentMessages,
		receivedMessages,
	)

	return &Metrics{
		connectedClients: connectedClients,
		sentMessages:     sentMessages,
		receivedMessages: receivedMessages,
	}
}

func (e *Metrics) ConnectedClientInc(websocket string) {
	e.connectedClients.WithLabelValues(websocket).Inc()
}

func (e *Metrics) ConnectedClientDec(websocket string) {
	e.connectedClients.WithLabelValues(websocket).Dec()
}

func (e *Metrics) SentMessagesInc(websocket string) {
	e.sentMessages.WithLabelValues(websocket).Inc()
}

func (e *Metrics) ReceivedMessages(channel string) {
	e.receivedMessages.WithLabelValues(channel).Inc()
}
