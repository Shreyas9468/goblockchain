package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"time"
)

// Block represents a single block in the blockchain
type Block struct {
	Timestamp    int64         `json:"timestamp"`
	Nonce        int           `json:"nonce"`
	PreviousHash [32]byte      `json:"previous_hash"`
	Transactions []*Transaction `json:"transactions"`
}

// NewBlock creates a new block instance
func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		Timestamp:    time.Now().UnixNano(),
		Nonce:        nonce,
		PreviousHash: previousHash,
		Transactions: transactions,
	}
}

// Hash computes the SHA-256 hash of the block
func (b *Block) Hash() [32]byte {
	data, _ := json.Marshal(b)
	return sha256.Sum256(data)
}