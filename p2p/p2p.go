package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rlaalsrl715/nomadcoin/blockchain"
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

func AddPeer(address, port, openPort string, broadcast bool) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	peer := initPeer(conn, address, port)
	Peers.m.Lock()
	defer Peers.m.Unlock()
	if broadcast {
		payload := fmt.Sprintf("%s:%s", peer.key, openPort)
		BroadcastToPeers(payload, MessageNewPeerNotify, peer)
	}
	sendNewestBlock(peer)
}

func BroadcastToPeers(payload interface{}, kind MessageKind, excludePeer *peer) {
	switch kind {
	case MessageNewBlockNotify:
		Peers.m.Lock()
		defer Peers.m.Unlock()
	case MessageNewTxNotify:
		blockchain.Mempool.M.Lock()
		defer blockchain.Mempool.M.Unlock()
	}
	for _, v := range Peers.v {
		if v.key != excludePeer.key {
			notify(payload, v, kind)
		}
	}
}
