package blockchain

import (
	"errors"
	"github.com/rlaalsrl715/nomadcoin/utils"
	"time"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID  string `json:"txID"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxID   string `json:"txID"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

func makeCoinBaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

func makeTx(from string, to string, amount int) (*Tx, error) {
	balance := Blockchain().BalanceByAddress(from)
	if balance < amount {
		return nil, errors.New("no money ")
	}
	uTxOuts := Blockchain().UTxOutsByAddress(from)

	var total int
	var txOuts []*TxOut
	var txIns []*TxIn
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		total += uTxOut.Amount
		txIns = append(txIns, &TxIn{uTxOut.TxID, uTxOut.Index, from})
	}
	if change := total - amount; change != 0 {
		txOuts = append(txOuts, &TxOut{from, change})
	}
	txOuts = append(txOuts, &TxOut{to, amount})

	tx := &Tx{TxOuts: txOuts, TxIns: txIns, Timestamp: int(time.Now().Unix())}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("minki", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) txToConfirm() []*Tx {
	coinbase := makeCoinBaseTx("minki")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
			}
		}
	}
	return exists

}
