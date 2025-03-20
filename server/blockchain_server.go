package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"goblockchain/blockchain"
	"goblockchain/network"
)

// BlockchainServer handles HTTP requests for the blockchain
type BlockchainServer struct {
	bc *blockchain.BlockChain
}

// NewBlockchainServer creates a new server instance
func NewBlockchainServer(bc *blockchain.BlockChain) *BlockchainServer {
	return &BlockchainServer{bc: bc}
}

// Run starts the blockchain server with automated mining and consensus
func (bs *BlockchainServer) Run(port string) {
	http.HandleFunc("/transactions/new", bs.newTransaction)
	http.HandleFunc("/mine", bs.mine)
	http.HandleFunc("/chain", bs.getChain)
	http.HandleFunc("/balance", bs.getBalance)
	http.HandleFunc("/consensus", bs.consensus)

	// Automate mining
	go func() {
		for {
			time.Sleep(10 * time.Second)
			if len(bs.bc.TransactionPool) > 0 {
				bs.bc.Mining()
				fmt.Println("Automatically mined a block")
			}
		}
	}()

	// Automate consensus
	go func() {
		for {
			time.Sleep(30 * time.Second)
			if network.ResolveConflicts(bs.bc) {
				fmt.Println("Chain replaced with longer chain from peers")
			}
		}
	}()

	fmt.Printf("Blockchain Server running on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func (bs *BlockchainServer) newTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var tx blockchain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if !bs.bc.AddTransaction(&tx) {
		http.Error(w, "Invalid transaction", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction added to pool"})
}

func (bs *BlockchainServer) mine(w http.ResponseWriter, r *http.Request) {
	if bs.bc.Mining() {
		json.NewEncoder(w).Encode(map[string]string{"message": "New block mined"})
	} else {
		http.Error(w, "Mining failed", http.StatusInternalServerError)
	}
}

func (bs *BlockchainServer) getChain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(bs.bc.Chain)
}

func (bs *BlockchainServer) getBalance(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter required", http.StatusBadRequest)
		return
	}
	balance := bs.bc.CalculateTotalAmount(address)
	json.NewEncoder(w).Encode(map[string]float32{"balance": balance})
}

func (bs *BlockchainServer) consensus(w http.ResponseWriter, r *http.Request) {
	if network.ResolveConflicts(bs.bc) {
		json.NewEncoder(w).Encode(map[string]string{"message": "Chain replaced with longer chain"})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"message": "Chain is authoritative"})
	}
}