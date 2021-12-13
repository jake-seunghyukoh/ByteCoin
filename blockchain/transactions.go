package blockchain

import (
	"errors"
	"github.com/ohshyuk5/ByteCoin/wallet"
	"time"

	"github.com/ohshyuk5/ByteCoin/utils"
)

const (
	minerReward int = 50
)

var Mempool *mempool = &mempool{} // No need to restore from db
var ErrorNoMoney = errors.New("not enough money")
var ErrorNotValid = errors.New("transaction invalid")

type mempool struct {
	Txs []*Tx
}

type Tx struct {
	ID        string   `json:"id"`
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
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(wallet.Wallet(), t.ID)
	}
}

func validate(tx *Tx) bool {
	for _, txIn := range tx.TxIns {
		prevTx := FindTransaction(BlockChain(), txIn.TxID)
		if prevTx == nil {
			return false
		}

		address := prevTx.TxOuts[txIn.Index].Address
		valid := wallet.Verify(txIn.Signature, tx.ID, address)
		if !valid {
			return false
		}
	}
	return true
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

func makeCoinbaseTx(address string) *Tx {
	txIn := []*TxIn{
		{TxID: "", Index: -1, Signature: "COINBASE"},
	}
	txOut := []*TxOut{
		{Address: address, Amount: minerReward},
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
		return nil, ErrorNoMoney
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
			TxID:      uTxOut.TxID,
			Index:     uTxOut.Index,
			Signature: from,
		}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{Address: from, Amount: change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOut := &TxOut{Address: to, Amount: amount}
	txOuts = append(txOuts, txOut)

	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrorNotValid
	}
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
