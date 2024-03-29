package blockchain

import (
	"errors"
	"github.com/rlaalsrl715/nomadcoin/utils"
	"github.com/rlaalsrl715/nomadcoin/wallet"
	"sync"
	"time"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs map[string]*Tx
	M   sync.Mutex
}

var Mempool *mempool = &mempool{}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID      string `json:"txID"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
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
	balance := BalanceByAddress(from, Blockchain())
	if balance < amount {
		return nil, errors.New("no money ")
	}
	uTxOuts := UTxOutsByAddress(from, Blockchain())

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
	sign(tx)
	valid := validate(tx)
	if !valid {
		return nil, errors.New("tx invaild")
	}
	return tx, nil
}

func sign(tx *Tx) {
	for _, txIn := range tx.TxIns {
		txIn.Signature = wallet.Sign(tx.Id, wallet.Wallet())
	}
}

func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(txIn.TxID, Blockchain())
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.Id, address)
		if !valid {
			break
		}
	}
	return valid
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return nil, err
	}
	m.Txs[tx.Id] = tx
	return tx, nil
}

func (m *mempool) txToConfirm() []*Tx {
	coinbase := makeCoinBaseTx(wallet.Wallet().Address)
	var txs []*Tx
	m.Txs[coinbase.Id] = coinbase
	for _, v := range m.Txs {
		txs = append(txs, v)
	}
	m.Txs = make(map[string]*Tx)
	return txs
}

func isOnMempool(uTxOut *UTxOut) bool {
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				return true

			}
		}
	}
	return false
}

func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, b := range Blocks(b) {
		txs = append(txs, b.Transactions...)
	}
	return txs
}

func FindTx(txId string, b *blockchain) *Tx {
	txs := Txs(b)
	for _, tx := range txs {
		if tx.Id == txId {
			return tx
		}
	}
	return nil
}

func clearMempoolByPeerBlock(confirmTx []*Tx) {
	for _, tx := range confirmTx {
		_, ok := Mempool.Txs[tx.Id]
		if ok {
			delete(Mempool.Txs, tx.Id)
		}
	}
}
