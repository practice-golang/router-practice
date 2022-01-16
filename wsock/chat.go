package wsock

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var WSockRooms map[uint64]WebSocketChatRoom = map[uint64]WebSocketChatRoom{}
var WSockWorkers map[uint64]WebSocketChatWorker = map[uint64]WebSocketChatWorker{}
var RoomMessage chan MsgShape
var BroadcastMessage chan MsgShape
var PersonIdx uint64 = 0
var RoomIdx uint64 = 0

func MessageSender() {
	for {
		select {
		case msg := <-BroadcastMessage:
			msg.Message = "(Broadcast) " + msg.Message
			for _, worker := range WSockWorkers {
				if worker.Conn != nil {
					worker.Message <- msg
				}
			}
		case msg := <-RoomMessage:
			msg.Message = "(Room) " + msg.Message
			for _, worker := range WSockRooms[msg.RoomIdx].Workers {
				if worker.Conn != nil {
					worker.Message <- msg
				}
			}
		}
	}
}

func WebSocketChat(w http.ResponseWriter, r *http.Request) {
	var err error
	worker := WebSocketChatWorker{}
	worker.Message = make(chan MsgShape)
	worker.Conn, _, _, err = ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println("ws UpgradeHTTP:", err)
		return
	}

	pidx := PersonIdx
	worker.Idx = pidx
	PersonIdx++
	WSockWorkers[pidx] = worker

	ridx := RoomIdx
	if WSockRooms[ridx].Workers == nil {
		WSockRooms[ridx] = WebSocketChatRoom{
			Name:    "Chat room title",
			Workers: map[uint64]WebSocketChatWorker{},
		}
	}
	WSockRooms[ridx].Workers[pidx] = worker
	// RoomIdx++

	// Publish message to all workers
	go func() {
		defer func() {
			WSockWorkers[pidx] = WebSocketChatWorker{}
			(worker.Conn).Close()
		}()
		for {
			recv, _, err := wsutil.ReadClientData(worker.Conn)
			if err != nil {
				log.Println("ws ReadClientData:", err)
				break
			}

			msg := MsgShape{}
			err = json.Unmarshal(recv, &msg)
			if err != nil {
				log.Println("ws ReadClientData:", err)
				break
			}

			switch msg.Target {
			case "room":
				RoomMessage <- msg
			case "broadcast":
				BroadcastMessage <- msg
			}
		}
	}()

	// Receives messages from the worker.
	go func() {
		for {
			// select {
			// case msg := <-worker.Message:
			msg := <-worker.Message
			// log.Println("#", worker.Idx, "Received:", msg)
			err = wsutil.WriteServerMessage(worker.Conn, ws.OpText, []byte(msg.Message))
			if err != nil {
				log.Println("ws WriteServerMessage:", err)
			}
			// }
		}
	}()
}

func InitWebSocketChat() {
	RoomMessage = make(chan MsgShape)
	BroadcastMessage = make(chan MsgShape)
	// Publisher or Broadcaster
	go MessageSender()
}
