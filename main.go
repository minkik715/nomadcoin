package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockchain struct {
	blocks []block
}

func (bc blockchain) getLastHash() string {
	size := len(bc.blocks)
	if size > 0 {
		return bc.blocks[size-1].hash
	}
	return ""
}

func (bc *blockchain) addBlock(data string) {
	prevHash := bc.getLastHash()
	newHash := fmt.Sprintf("%x", sha256.Sum256([]byte(data+prevHash)))
	bc.blocks = append(bc.blocks, block{data: data, hash: newHash, prevHash: prevHash})
}

func (bc blockchain) listBlocks() {
	for _, b := range bc.blocks {
		fmt.Printf("data: %s  hash: %s  preveHash: %s \n", b.data, b.hash, b.prevHash)
	}
}

func main() {
	chain := blockchain{}
	chain.addBlock("minki")
	chain.addBlock("yeonhi")
	chain.addBlock("nico")
	chain.listBlocks()
}
