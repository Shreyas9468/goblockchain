package server

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"goblockchain/wallet"

	"crypto/elliptic"
	"crypto/sha256"
)

// WalletServer handles wallet-related HTTP requests
type WalletServer struct {
	bcURL string
}

// NewWalletServer creates a new wallet server instance
func NewWalletServer(bcURL string) *WalletServer {
	return &WalletServer{bcURL: bcURL}
}

// Run starts the wallet server
func (ws *WalletServer) Run(port string) {
	http.HandleFunc("/wallet/new", ws.newWallet)
	http.HandleFunc("/wallet/balance", ws.getBalance)
	http.HandleFunc("/wallet/transaction", ws.createTransaction)
	fmt.Printf("Wallet Server running on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func (ws *WalletServer) newWallet(w http.ResponseWriter, r *http.Request) {
	walletInstance, err := wallet.NewWallet()
	if err != nil {
		http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
		return
	}
	response := map[string]string{
		"private_key": walletInstance.SerializePrivateKey(),
		"public_key":  hex.EncodeToString(walletInstance.PublicKey),
		"address":     walletInstance.Address,
	}
	json.NewEncoder(w).Encode(response)
}

func (ws *WalletServer) getBalance(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter required", http.StatusBadRequest)
		return
	}
	resp, err := http.Get(ws.bcURL + "/balance?address=" + url.QueryEscape(address))
	if err != nil {
		http.Error(w, "Failed to fetch balance", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body into a byte slice
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read balance response", http.StatusInternalServerError)
		return
	}

	// Set headers and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(body)
	if err != nil {
		// Log this error if needed; it's rare since Write rarely fails after headers are set
		fmt.Println("Error writing response:", err)
	}
}

func (ws *WalletServer) createTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		PrivateKey string  `json:"private_key"`
		Recipient  string  `json:"recipient"`
		Value      float32 `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	privateKey, err := wallet.DeserializePrivateKey(req.PrivateKey)
	if err != nil {
		http.Error(w, "Invalid private key", http.StatusBadRequest)
		return
	}
	walletInstance := &wallet.Wallet{
		PrivateKey: privateKey,
		PublicKey:  elliptic.Marshal(elliptic.P256(), privateKey.PublicKey.X, privateKey.PublicKey.Y),
	}
	hash := sha256.Sum256(walletInstance.PublicKey)
	walletInstance.Address = hex.EncodeToString(hash[:])

	tx, err := walletInstance.CreateTransaction(req.Recipient, req.Value)
	if err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}
	txData, err := json.Marshal(tx)
	if err != nil {
		http.Error(w, "Failed to marshal transaction", http.StatusInternalServerError)
		return
	}
	resp, err := http.Post(ws.bcURL+"/transactions/new", "application/json", bytes.NewBuffer(txData))
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to submit transaction", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction sent to blockchain"})
}