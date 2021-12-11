package blockchain

import (
	"fmt"
	"sync"

	"github.com/ohshyuk5/ByteCoin/db"
	"github.com/ohshyuk5/ByteCoin/utils"
)

const (
	defaultDifficulty   int = 2
	recalculateInterval int = 5
	blockInterval       int = 2
	tolerance           int = 2
)

var b *blockChain
var once sync.Once

type blockChain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

func (b *blockChain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b blockChain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockChain) AddBlock() {
	block := createBlock()
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockChain) Blocks() []*Block {
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
func (b *blockChain) recalculateDifficulty() int {
	allBlocks := b.Blocks()

	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[recalculateInterval-1]

	actualTime := (newestBlock.Timestamp - lastRecalculatedBlock.Timestamp) / 60
	expectedTime := blockInterval * recalculateInterval

	if actualTime < expectedTime-tolerance {
		return b.CurrentDifficulty + 1
	}
	if actualTime > expectedTime+tolerance {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func (b *blockChain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	}
	if b.Height%recalculateInterval == 0 {
		return b.recalculateDifficulty()
	}
	return b.CurrentDifficulty
}

func (b *blockChain) txOuts() []*TxOut {
	var txOuts []*TxOut

	blocks := b.Blocks()
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txOuts = append(txOuts, tx.TxOuts...)
		}
	}

	return txOuts
}

func (b *blockChain) TxOutsByAddress(address string) []*TxOut {
	var ownedTxOuts []*TxOut
	txOuts := b.txOuts()

	for _, txOut := range txOuts {
		if txOut.Owner == address {
			ownedTxOuts = append(ownedTxOuts, txOut)
		}
	}

	return ownedTxOuts
}

func (b *blockChain) BalanceByAddress(address string) int {
	var amount int

	txOuts := b.TxOutsByAddress(address)
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}

	return amount
}

func BlockChain() *blockChain {
	if b == nil {
		once.Do(func() {
			b = &blockChain{Height: 0, CurrentDifficulty: defaultDifficulty}
			checkpoint := db.Blockchain()
			if checkpoint == nil {
				fmt.Println("Initializing...")
				b.AddBlock()
			} else {
				fmt.Println("Restoring...")
				b.restore(checkpoint)
			}
		})
	}
	return b
}
