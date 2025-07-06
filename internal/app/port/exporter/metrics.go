package exporter

type ConnectedClientsMetric interface {
	Inc(websocket string)
	Dec(websocket string)
}

type MessageCounterMetric interface {
	Inc(channel string)
}
