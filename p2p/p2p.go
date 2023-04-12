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
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	openPort := r.URL.Query().Get("openPort")
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	initPeer(conn, ip, openPort)
}

func AddPeer(address, port, openPort string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:]), nil)
	utils.HandleErr(err)
	peer := initPeer(conn, address, port)
	sendNewestBlock(peer)
}
