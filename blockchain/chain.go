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

func persistBlockchain(b *blockChain) {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockChain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
}

func Blocks(b *blockChain) []*Block {
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

func recalculateDifficulty(b *blockChain) int {
	allBlocks := Blocks(b)

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

func getDifficulty(b *blockChain) int {
	if b.Height == 0 {
		return defaultDifficulty
	}
	if b.Height%recalculateInterval == 0 {
		return recalculateDifficulty(b)
	}
	return b.CurrentDifficulty
}

func UTxOutsByAddress(b *blockChain, address string) []*UTxOut {
	var uTxOuts []*UTxOut               // Unspent Transaction Outputs
	creatorTxs := make(map[string]bool) // Map of transactions

	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if _, spent := creatorTxs[tx.ID]; !spent {
						uTxOut := &UTxOut{tx.ID, index, output.Amount}

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

func BalanceByAddress(b *blockChain, address string) int {
	var amount int

	txOuts := UTxOutsByAddress(b, address)
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}

	return amount
}

func BlockChain() *blockChain {
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
	return b
}
