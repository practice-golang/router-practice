package wsock

import "net"

type MsgShape struct {
	Target  string `json:"target"`
	RoomIdx uint64 `json:"-"`
	Message string `json:"message"`
}

type WebSocketChatRoom struct {
	Name    string
	Workers map[uint64]WebSocketChatWorker
}

type WebSocketChatWorker struct {
	Idx     uint64
	Conn    net.Conn
	Message chan MsgShape
}
