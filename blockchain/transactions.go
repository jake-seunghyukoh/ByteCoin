package blockchain

import (
	"errors"
	"time"

	"github.com/ohshyuk5/ByteCoin/utils"
)

const (
	minerReward int = 50
)

var Mempool *mempool = &mempool{} // No need to restore from db

type mempool struct {
	Txs []*Tx
}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
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

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false

Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
}

func makeCoinbaseTx(address string) *Tx {
	txIn := []*TxIn{
		{TxID: "", Index: -1, Owner: "COINBASE"},
	}
	txOut := []*TxOut{
		{Owner: address, Amount: minerReward},
	}
	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIn,
		TxOuts:    txOut,
	}
	tx.getId()

	return &tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(BlockChain(), from) < amount {
		return nil, errors.New("not enough money")
	}
	var txOuts []*TxOut
	var txIns []*TxIn

	total := 0
	uTxOuts := UTxOutsByAddress(BlockChain(), from)

	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}

		txIn := &TxIn{
			TxID:  uTxOut.TxID,
			Index: uTxOut.Index,
			Owner: from,
		}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{Owner: from, Amount: change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOut := &TxOut{Owner: to, Amount: amount}
	txOuts = append(txOuts, txOut)

	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()

	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("me", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("me")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
