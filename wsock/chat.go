package wsock

import (
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var WebSockerWorkers map[uint64]WebSocketChatWorker = map[uint64]WebSocketChatWorker{}
var Message chan string
var Idx uint64 = 0

func Publisher() {
	Message = make(chan string, 10)
	for {
		select {
		case msg := <-Message:
			for _, worker := range WebSockerWorkers {
				if worker.conn != nil {
					// log.Println("Send message to worker #:", worker.Idx)
					worker.msgCH <- msg
				}
			}
		}
	}
}

func WebSocketChat(r *http.Request, w http.ResponseWriter) {
	var err error
	worker := WebSocketChatWorker{}
	worker.msgCH = make(chan string, 10)
	worker.conn, _, _, err = ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println("ws UpgradeHTTP:", err)
	}

	idx := Idx
	worker.Idx = idx
	Idx++
	WebSockerWorkers[idx] = worker

	// Publisher or Broadcaster
	go Publisher()

	// Publish message to all workers
	go func() {
		defer func() {
			WebSockerWorkers[idx] = WebSocketChatWorker{}
			(worker.conn).Close()
		}()
		for {
			recv, _, err := wsutil.ReadClientData(worker.conn)
			if err != nil {
				log.Println("ws ReadClientData:", err)
				break
			}
			Message <- string(recv)
		}
	}()

	// Receives messages from the worker.
	go func() {
		for {
			select {
			case msg := <-worker.msgCH:
				// log.Println("#", worker.Idx, "Received:", msg)
				err = wsutil.WriteServerMessage(worker.conn, ws.OpText, []byte(msg))
				if err != nil {
					log.Println("ws WriteServerMessage:", err)
				}
			}
		}
	}()
}
