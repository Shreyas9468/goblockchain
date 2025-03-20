package main

import (
	"goblockchain/blockchain"
	"goblockchain/server"
)



func main() {
	// Initialize blockchain node
	nodeAddress := "node1_address"
	bc := blockchain.NewBlockChain(nodeAddress)

	// Optionally add peers (uncomment to use)
	// bc.AddPeer("http://localhost:8081")
	// bc.AddPeer("http://localhost:8082")

	// Start blockchain server
	bs := server.NewBlockchainServer(bc)
	go bs.Run("8080")

	// Start wallet server
	ws := server.NewWalletServer("http://localhost:8080")
	go ws.Run("8081")

	// Keep the main goroutine running
	select {}
}