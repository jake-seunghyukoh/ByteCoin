package blockchain

import (
	"errors"
	"github.com/ohshyuk5/ByteCoin/db"
	"github.com/ohshyuk5/ByteCoin/utils"
	"strings"
	"time"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

var ErrNotFound = errors.New("block not found")

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

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
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock() *Block {
	block := Block{
		Hash:       "",
		PrevHash:   b.NewestHash,
		Height:     b.Height + 1,
		Difficulty: BlockChain().difficulty(),
		Nonce:      0,
		Timestamp:  0,
	}
	block.mine()
	block.Transactions = Mempool.TxToConfirm()
	block.persist()
	return &block
}
