package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	
	"math/big"

	"goblockchain/blockchain"
)

// Wallet manages cryptographic keys and transaction creation
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Address    string
}

// NewWallet generates a new wallet with ECDSA keys
func NewWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	publicKey := elliptic.Marshal(elliptic.P256(), privateKey.X, privateKey.Y)
	hash := sha256.Sum256(publicKey)
	address := hex.EncodeToString(hash[:])
	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}

// CreateTransaction generates and signs a new transaction
func (w *Wallet) CreateTransaction(recipient string, value float32) (*blockchain.Transaction, error) {
	tx := &blockchain.Transaction{
		SenderPublicKey:  w.PublicKey,
		RecipientAddress: recipient,
		Value:            value,
	}
	data, err := tx.MarshalForSigning()
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(data)
	signature, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}
	tx.Signature = signature
	return tx, nil
}

// SerializePrivateKey converts the private key to a hex string
func (w *Wallet) SerializePrivateKey() string {
	return hex.EncodeToString(w.PrivateKey.D.Bytes())
}

// DeserializePrivateKey reconstructs a private key from a hex string
func DeserializePrivateKey(hexStr string) (*ecdsa.PrivateKey, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	d := new(big.Int).SetBytes(bytes)
	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
		},
		D: d,
	}
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.Curve.ScalarBaseMult(d.Bytes())
	return privateKey, nil
}