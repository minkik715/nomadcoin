package p2p

import (
	"encoding/json"
	"github.com/rlaalsrl715/nomadcoin/blockchain"
	"github.com/rlaalsrl715/nomadcoin/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func sendNewestBlock(p *peer) {
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{kind, utils.ToJsonBytes(payload)}
	return utils.ToJsonBytes(m)
}

func handleMessage(msg *Message, p *peer) {
	switch msg.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(msg.Payload, &payload))

	}
}
