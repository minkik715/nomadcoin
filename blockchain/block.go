package blockchain

import (
	"errors"
	"github.com/rlaalsrl715/nomadcoin/db"
	"github.com/rlaalsrl715/nomadcoin/utils"
	"strings"
	"time"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	TimesStamp   int    `json:"timesStamp"`
	Transactions []*Tx  `json:"transactions"`
}

const difficulty int = 2

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

var ErrNotFound = errors.New("block not found")

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.TimesStamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBLock(prevHash string, height int, diff int) *Block {
	block := Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: diff,
		Nonce:      0,
	}
	block.mine()
	block.Transactions = Mempool.txToConfirm()
	block.persist()
	return &block
}

func AddPeerBlock(block *Block) {
	b.m.Lock()
	defer b.m.Unlock()
	b.NewestHash = block.Hash
	b.CurrentDifficulty = block.Difficulty
	b.Height = block.Height
	persistBlockchain(b)
	block.persist()
	clearMempoolByPeerBlock(block.Transactions)
}
