package p2p

import (
	"github.com/gorilla/websocket"
	"github.com/rlaalsrl715/nomadcoin/utils"
	"net/http"
)

var upgrader = websocket.Upgrader{}

var conns []*websocket.Conn

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	conns = append(conns, conn)
	utils.HandleErr(err)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {

			err = conn.Close()
			utils.HandleErr(err)
			break
		}
		for _, aConn := range conns {
			if aConn != conn {
				err = aConn.WriteMessage(websocket.TextMessage, p)
				utils.HandleErr(err)
			}
		}
	}
}
