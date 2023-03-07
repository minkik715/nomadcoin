package blockchain

import (
	"fmt"
	"github.com/rlaalsrl715/nomadcoin/db"
	"github.com/rlaalsrl715/nomadcoin/utils"
	"sync"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

const (
	defaultDifficulty  = 2
	difficultyInterval = 5
	blockInterval      = 2
	allowedRange       = 2
)

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) AddBlock(data string) {
	block := createBLock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = b.Difficulty()
	b.persist()
}

func (b *blockchain) persist() {
	db.SaveCheckPoint(utils.ToBytes(b))
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{Height: 0}
			checkPoint := db.Blockchain()
			if checkPoint == nil {
				b.AddBlock("Genesis Block")
			} else {
				b.restore(checkPoint)
			}
		})
		println(b.NewestHash)
	}
	return b
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func (b *blockchain) Difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

func (b *blockchain) recalculateDifficulty() int {
	allBocks := b.Blocks()
	newestBlock := allBocks[0]
	lastRecalculatedBlock := allBocks[difficultyInterval-1]
	actualTime := (newestBlock.TimesStamp / 60) - (lastRecalculatedBlock.TimesStamp / 60)
	expectedTime := difficultyInterval * blockInterval

	if actualTime <= (expectedTime - allowedRange) {
		fmt.Println(actualTime)
		fmt.Println(expectedTime)
		fmt.Println("UP")

		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		fmt.Println(actualTime)
		fmt.Println(expectedTime)
		fmt.Println("DOWN")

		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}
