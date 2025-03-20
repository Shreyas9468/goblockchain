package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/json"
)

// Transaction represents a cryptocurrency transaction
type Transaction struct {
	SenderPublicKey  []byte  `json:"sender_public_key"`
	RecipientAddress string  `json:"recipient_address"`
	Value            float32 `json:"value"`
	Signature        []byte  `json:"signature"`
}

// MarshalForSigning serializes the transaction data for signing
func (t *Transaction) MarshalForSigning() ([]byte, error) {
	type txForSigning struct {
		SenderPublicKey  []byte  `json:"sender_public_key"`
		RecipientAddress string  `json:"recipient_address"`
		Value            float32 `json:"value"`
	}
	return json.Marshal(txForSigning{
		SenderPublicKey:  t.SenderPublicKey,
		RecipientAddress: t.RecipientAddress,
		Value:            t.Value,
	})
}

// VerifyTransaction checks the validity of a transaction's signature
func VerifyTransaction(tx *Transaction) bool {
	if tx.SenderPublicKey == nil && tx.Signature == nil {
		return true // Mining reward transaction
	}
	data, err := tx.MarshalForSigning()
	if err != nil {
		return false
	}
	hash := sha256.Sum256(data)
	x, y := elliptic.Unmarshal(elliptic.P256(), tx.SenderPublicKey)
	if x == nil || y == nil {
		return false
	}
	pubKey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
	return ecdsa.VerifyASN1(pubKey, hash[:], tx.Signature)
}