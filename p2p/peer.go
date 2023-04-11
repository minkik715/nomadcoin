package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
)

var Peer map[string]*peer = make(map[string]*peer)

type peer struct {
	conn  *websocket.Conn
	inbox chan []byte
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	p := peer{conn, make(chan []byte)}
	key := fmt.Sprintf("%s:%s", address, port)
	go p.read()
	Peer[key] = &p
	return &p
}
func (p *peer) read() {
	// delete peer in case of error
	for {
		_, m, err := p.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("%s", m)
	}
}

func (p *peer) write() {
	// delete peer in case of error
	for {
		msg := <-p.inbox
		p.conn.WriteMessage(websocket.TextMessage, msg)
	}
}
