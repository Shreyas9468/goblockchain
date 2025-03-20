package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Constants for blockchain configuration
const (
	MINING_DIFFICULTY = 3
	MINING_REWARD     = 1.0
)

// BlockChain represents the blockchain with its chain and transaction pool
type BlockChain struct {
	TransactionPool   []*Transaction
	Chain             []*Block
	BlockchainAddress string
	Peers             []string
}

// NewBlockChain initializes a blockchain with a genesis block
func NewBlockChain(blockchainAddress string) *BlockChain {
	bc := &BlockChain{
		BlockchainAddress: blockchainAddress,
	}
	genesisBlock := NewBlock(0, [32]byte{}, []*Transaction{})
	bc.Chain = append(bc.Chain, genesisBlock)
	return bc
}

// AddTransaction adds a verified transaction to the pool
func (bc *BlockChain) AddTransaction(tx *Transaction) bool {
	if !VerifyTransaction(tx) {
		return false
	}
	bc.TransactionPool = append(bc.TransactionPool, tx)
	return true
}

// Mining performs proof-of-work and adds a new block
func (bc *BlockChain) Mining() bool {
	rewardTx := &Transaction{
		SenderPublicKey:  nil,
		RecipientAddress: bc.BlockchainAddress,
		Value:            MINING_REWARD,
		Signature:        nil,
	}
	bc.TransactionPool = append([]*Transaction{rewardTx}, bc.TransactionPool...)

	nonce := bc.ProofOfWork()
	previousHash := bc.Chain[len(bc.Chain)-1].Hash()
	bc.CreateBlock(nonce, previousHash)
	return true
}

// CreateBlock adds a new block to the chain and clears the transaction pool
func (bc *BlockChain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	block := NewBlock(nonce, previousHash, bc.TransactionPool)
	bc.Chain = append(bc.Chain, block)
	bc.TransactionPool = []*Transaction{}
	return block
}

// ProofOfWork finds a nonce that satisfies the difficulty requirement
func (bc *BlockChain) ProofOfWork() int {
	transactions := append([]*Transaction{}, bc.TransactionPool...)
	previousHash := bc.Chain[len(bc.Chain)-1].Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions) {
		nonce++
	}
	return nonce
}

// ValidProof checks if a nonce produces a valid hash
func (bc *BlockChain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction) bool {
	guessBlock := Block{
		Timestamp:    time.Now().UnixNano(),
		Nonce:        nonce,
		PreviousHash: previousHash,
		Transactions: transactions,
	}
	hash := fmt.Sprintf("%x", guessBlock.Hash())
	return hash[:MINING_DIFFICULTY] == strings.Repeat("0", MINING_DIFFICULTY)
}

// CalculateTotalAmount computes the balance for a given address
func (bc *BlockChain) CalculateTotalAmount(address string) float32 {
	var total float32
	for _, block := range bc.Chain {
		for _, tx := range block.Transactions {
			if tx.RecipientAddress == address {
				total += tx.Value
			}
			if tx.SenderPublicKey != nil {
				hash := sha256.Sum256(tx.SenderPublicKey)
				senderAddress := hex.EncodeToString(hash[:])
				if senderAddress == address {
					total -= tx.Value
				}
			}
		}
	}
	return total
}

// AddPeer registers a new peer node
func (bc *BlockChain) AddPeer(peerURL string) {
	bc.Peers = append(bc.Peers, peerURL)
}