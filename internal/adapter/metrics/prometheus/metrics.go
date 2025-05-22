package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	connectedClients *prometheus.GaugeVec
}

func NewMetrics(name string) *Metrics {
	connectedClients := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "connected_clients_total",
		Name:      name,
	}, []string{"websocket"})

	prometheus.MustRegister(connectedClients)

	return &Metrics{
		connectedClients: connectedClients,
	}
}

func (e *Metrics) ConnectedClientInc(websocket string) {
	e.connectedClients.WithLabelValues(websocket).Inc()
}

func (e *Metrics) ConnectedClientDec(websocket string) {
	e.connectedClients.WithLabelValues(websocket).Dec()
}
