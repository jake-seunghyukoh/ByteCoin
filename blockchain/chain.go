package blockchain

import (
	"fmt"
	"sync"

	"github.com/ohshyuk5/ByteCoin/db"
	"github.com/ohshyuk5/ByteCoin/utils"
)

var b *blockChain
var once sync.Once

type blockChain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

func (b *blockChain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b blockChain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockChain) AddBlock(data string) {
	block := createBlock(data)
	b.NewestHash = block.Hash
	b.Height = block.Height
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

func BlockChain() *blockChain {
	if b == nil {
		once.Do(func() {
			b = &blockChain{"", 0}
			checkpoint := db.Blockchain()
			if checkpoint == nil {
				fmt.Println("Initializing...")
				b.AddBlock("Genesis")
			} else {
				fmt.Println("Restoring...")
				b.restore(checkpoint)
			}
		})
	}
	return b
}
