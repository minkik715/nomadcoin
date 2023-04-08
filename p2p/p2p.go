package p2p

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rlaalsrl715/nomadcoin/utils"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			err = conn.Close()
			utils.HandleErr(err)
			break
		}
		fmt.Printf("Just got %s\n\n", p)
		time.Sleep(time.Second * 5)
		message := fmt.Sprintf("we also think thant %s", p)
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		utils.HandleErr(err)
	}
}
