package p2p

import (
	"fmt"
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
	utils.HandleErr(err)
	initPeer(conn, "xx", "xx")
}

func AddPeer(address, port string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws:%s:%s/ws", address, port), nil)
	utils.HandleErr(err)
	initPeer(conn, address, port)
}
