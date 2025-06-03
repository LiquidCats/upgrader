package exporter

type ConnectedClientsMetrics interface {
	ConnectedClientInc(websocket string)
	ConnectedClientDec(websocket string)
}

type MessagesMetrics interface {
	SentMessagesInc(websocket string)
	ReceivedMessages(channel string)
}
