package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rlaalsrl715/nomadcoin/utils"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	for {
		fmt.Println("wait catching msg")
		_, p, err := conn.ReadMessage()
		fmt.Println("receive msge")
		utils.HandleErr(err)
		fmt.Printf("%s", p)
	}
}
