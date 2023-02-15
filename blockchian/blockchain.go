package blockchian

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data     string
	Hash     string
	prevHash string
}

type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once

func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].Hash
}

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.prevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash()}
	newBlock.calculateHash()
	return &newBlock
}

// 언제 receiver function 사용하고 언제 ...하는지
// 선언 후에 하는 function이라면 receiver fucnition으로 구현한디.
// 선언을 안하거나 선언 전에하는 function 이라면 그냥 fucntion
func (b *blockchain) AppendBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AppendBlock("Genesis Block")
		})
	}
	return b
}

func (b *blockchain) AllBLocks() []*Block {
	return b.blocks
}
