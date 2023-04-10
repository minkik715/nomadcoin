package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
)

var Peer map[string]*peer = make(map[string]*peer)

type peer struct {
	conn *websocket.Conn
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	p := peer{conn}
	key := fmt.Sprintf("%s:%s", address, port)
	Peer[key] = &p
	return &p
}
