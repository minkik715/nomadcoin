package blockchain

import (
	"encoding/json"
	"github.com/rlaalsrl715/nomadcoin/db"
	"github.com/rlaalsrl715/nomadcoin/utils"
	"net/http"
	"sync"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
	m                 sync.Mutex
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

func (b *blockchain) AddBlock() *Block {
	b.m.Lock()
	Mempool.M.Lock()
	defer b.m.Unlock()
	defer Mempool.M.Unlock()
	block := createBLock(b.NewestHash, b.Height+1, Difficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = Difficulty(b)
	persistBlockchain(b)
	return block
}

func persistBlockchain(b *blockchain) {
	db.SaveCheckPoint(utils.ToBytes(b))
}

func Blocks(b *blockchain) []*Block {
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

func Difficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func recalculateDifficulty(b *blockchain) int {
	allBocks := Blocks(b)
	newestBlock := allBocks[0]
	lastRecalculatedBlock := allBocks[difficultyInterval-1]
	actualTime := (newestBlock.TimesStamp / 60) - (lastRecalculatedBlock.TimesStamp / 60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(input.TxID, Blockchain()).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Address == address {
					if _, ok := creatorTxs[tx.Id]; !ok {
						uTxOut := &UTxOut{tx.Id, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, txOUt := range txOuts {
		amount += txOUt.Amount
	}
	return amount
}

func Blockchain() *blockchain {
	once.Do(
		func() {
			b = &blockchain{Height: 0, CurrentDifficulty: defaultDifficulty}
			checkPoint := db.Blockchain()
			if checkPoint == nil {
				b.AddBlock()
			} else {
				b.restore(checkPoint)
			}
		},
	)
	return b
}

func (b *blockchain) Replace(blocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()
	b.CurrentDifficulty = blocks[0].Difficulty
	b.Height = len(blocks)
	b.NewestHash = blocks[0].Hash
	persistBlockchain(b)
	db.RemoveAllBlocks()
	for _, block := range blocks {
		block.persist()
	}
}

func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()
	encoder := json.NewEncoder(rw)
	encoder.Encode(b)
}
