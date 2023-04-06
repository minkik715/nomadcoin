package p2p

import (
	"github.com/gorilla/websocket"
	"github.com/rlaalsrl715/nomadcoin/utils"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	_, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
}
