package exporter

type ConnectedClientsMetric interface {
	Inc(websocket string)
	Dec(websocket string)
}

type SentMessagesMetric interface {
	Inc(websocket string)
}

type ReceivedMessagesMetric interface {
	Inc(channel string)
}
