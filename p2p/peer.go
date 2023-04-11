package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
)

var Peers = make(map[string]*peer)

type peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	key := fmt.Sprintf("%s:%s", address, port)
	p := peer{
		key,
		address,
		port,
		conn,
		make(chan []byte),
	}
	go p.read()
	Peers[key] = &p
	return &p
}
func (p *peer) read() {
	defer p.close()
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
	defer p.close()
	for {
		msg, ok := <-p.inbox
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func (p *peer) close() {
	p.conn.Close()
	delete(Peers, p.key)
}
