package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ConnectedClient struct {
	internal *prometheus.GaugeVec
}

func NewConnectedClient(name string) *ConnectedClient {
	return &ConnectedClient{
		internal: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: name,
			Name:      "connected_clients_total",
		}, []string{"websocket"}),
	}
}

func (e *ConnectedClient) Inc(websocket string) {
	e.internal.WithLabelValues(websocket).Inc()
}
func (e *ConnectedClient) Dec(websocket string) {
	e.internal.WithLabelValues(websocket).Dec()
}

type SentMessages struct {
	internal *prometheus.CounterVec
}

func NewSentMessages(name string) *SentMessages {
	return &SentMessages{
		internal: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: name,
			Name:      "sent_messages_total",
		}, []string{"websocket"}),
	}
}

func (e *SentMessages) Inc(websocket string) {
	e.internal.WithLabelValues(websocket).Inc()
}

type ReceivedMessages struct {
	internal *prometheus.CounterVec
}

func NewReceivedMessages(name string) *ReceivedMessages {
	return &ReceivedMessages{
		internal: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: name,
			Name:      "received_messages_total",
		}, []string{"websocket"}),
	}
}

func (e *ReceivedMessages) Inc(websocket string) {
	e.internal.WithLabelValues(websocket).Inc()
}
