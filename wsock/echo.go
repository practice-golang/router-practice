package wsock

import (
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func SockEcho(r *http.Request, w http.ResponseWriter) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println("ws UpgradeHTTP:", err)
	}

	go func() {
		defer conn.Close()
		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				log.Println("ws ReadClientData:", err)
			}

			msg = append([]byte("From server: "), msg...)

			err = wsutil.WriteServerMessage(conn, op, msg)
			if err != nil {
				log.Println("ws WriteServerMessage:", err)
			}
		}
	}()
}
