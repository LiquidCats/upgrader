package exporter

type ConnectedClientsMetrics interface {
	ConnectedClientInc(websocket string)
	ConnectedClientDec(websocket string)
}
