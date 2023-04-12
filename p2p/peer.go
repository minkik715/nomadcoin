package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type peers struct {
	v map[string]*peer
	m sync.Mutex
}

var Peers peers = peers{
	v: make(map[string]*peer),
}

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
	go p.write()
	Peers.v[key] = &p
	return &p
}
func (p *peer) read() {
	defer p.close()
	// delete peer in case of error
	for {
		m := Message{}
		err := p.conn.ReadJSON(&m)
		if err != nil {
			break
		}
		fmt.Println(m.Kind)
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
	Peers.m.Lock()
	defer Peers.m.Unlock()
	p.conn.Close()
	delete(Peers.v, p.key)
}

func AllPeers(p *peers) []string {
	p.m.Lock()
	defer p.m.Unlock()
	var keys []string
	for key := range p.v {
		keys = append(keys, key)
	}
	return keys
}
