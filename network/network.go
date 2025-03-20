package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"goblockchain/blockchain"
)

// ResolveConflicts adopts the longest valid chain from peers
func ResolveConflicts(bc *blockchain.BlockChain) bool {
	var longestChain []*blockchain.Block
	maxLength := len(bc.Chain)

	for _, peer := range bc.Peers {
		resp, err := http.Get(peer + "/chain")
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		var chain []*blockchain.Block
		if err := json.NewDecoder(resp.Body).Decode(&chain); err != nil {
			continue
		}
		if len(chain) > maxLength && ValidChain(chain) {
			maxLength = len(chain)
			longestChain = chain
		}
	}

	if longestChain != nil {
		bc.Chain = longestChain
		return true
	}
	return false
}

// ValidChain verifies the integrity of a chain
func ValidChain(chain []*blockchain.Block) bool {
	if len(chain) == 0 {
		return false
	}
	// Check genesis block
	if chain[0].PreviousHash != [32]byte{} {
		return false
	}
	for i := 1; i < len(chain); i++ {
		current := chain[i]
		previous := chain[i-1]
		if current.PreviousHash != previous.Hash() {
			return false
		}
		hash := fmt.Sprintf("%x", current.Hash())
		if hash[:blockchain.MINING_DIFFICULTY] != strings.Repeat("0", blockchain.MINING_DIFFICULTY) {
			return false
		}
		for _, tx := range current.Transactions {
			if !blockchain.VerifyTransaction(tx) {
				return false
			}
		}
	}
	return true
}