package wsock

import "net"

type WebSocketChatWorker struct {
	conn  net.Conn
	msgCH chan string
	Idx   uint64
}
