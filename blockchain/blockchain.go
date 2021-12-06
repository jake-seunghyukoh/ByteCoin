package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

type blockChain struct {
	blocks []*Block
}

var b *blockChain
var once sync.Once
var ErrNotFound error = errors.New("block not Found")

func getLastHash() string {
	totalBlocks := len(GetBlockChain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockChain().blocks[totalBlocks-1].Hash
}

func (b *Block) calculateHash() {
	data := b.Data
	prevHash := b.PrevHash

	hash := sha256.Sum256([]byte(data + prevHash))

	b.Hash = fmt.Sprintf("%x", hash)
}

func createBlock(data string) *Block {
	newBlock := Block{Data: data, Hash: "", PrevHash: getLastHash(), Height: len(GetBlockChain().blocks) + 1}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockChain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockChain() *blockChain {
	if b == nil {
		once.Do(func() {
			b = &blockChain{}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}

func (b *blockChain) AllBlocks() []*Block {
	return b.blocks
}

func (b *blockChain) GetBlock(height int) (*Block, error) {
	if height > len(b.blocks) {
		return nil, ErrNotFound
	}
	return b.blocks[height-1], nil
}
