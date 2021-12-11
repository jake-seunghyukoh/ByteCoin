package blockchain

import (
	"errors"
	"github.com/ohshyuk5/ByteCoin/utils"
	"time"
)

const (
	minerReward int = 50
)

var Mempool *mempool = &mempool{} // No need to restore from db

type mempool struct {
	Txs []*Tx
}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

func makeCoinbaseTx(address string) *Tx {
	txIn := []*TxIn{
		{Owner: "COINBASE", Amount: minerReward},
	}
	txOut := []*TxOut{
		{Owner: address, Amount: minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIn,
		TxOuts:    txOut,
	}
	tx.getId()

	return &tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	balance := BlockChain().BalanceByAddress(from)
	if balance < amount {
		return nil, errors.New("not enough money")
	}
	var txIns []*TxIn
	var txOuts []*TxOut

	total := 0
	oldTxOuts := BlockChain().TxOutsByAddress(from)
	for _, txOut := range oldTxOuts {
		if total > amount {
			break
		}
		txIn := &TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, txIn)
		total += txIn.Amount
	}

	change := total - amount
	if change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		Id:        "",
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
