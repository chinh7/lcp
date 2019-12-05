package crypto

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	cdc "github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"
)

const (
	// MethodNameByteLength is number of bytes preservered for method name
	MethodNameByteLength = 64
)

// TxData data for contract deploy/invoke
type TxData struct {
	Method string
	Params []byte
}

// TxSigner information about transaction signer
type TxSigner struct {
	PubKey    []byte
	Nonce     uint64
	Signature []byte
}

// Tx transaction
type Tx struct {
	From     TxSigner
	Data     []byte
	To       Address
	GasLimit uint64
}

// Address derived from TxSigner PubKey
func (txSigner *TxSigner) Address() Address {
	return AddressFromPubKey(txSigner.PubKey)
}

// CreateAddress create a new contract address based on pubkey and nonce
func (txSigner *TxSigner) CreateAddress() Address {
	// data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	cloned := &TxSigner{PubKey: txSigner.PubKey, Nonce: txSigner.Nonce}
	var res = sha3.Sum256(cloned.Serialize())
	return AddressFromPubKey(res[:])
}

// String TxSigner string presentation
func (txSigner TxSigner) String() string {
	return fmt.Sprintf("Nonce: %d Pubkey: %X Signature: %X", txSigner.Nonce, txSigner.PubKey, txSigner.Signature)
}

// String Tx string presentation
func (tx Tx) String() string {
	return fmt.Sprintf("Data: %s To: %s Signer: %s", hex.EncodeToString(tx.Data), tx.To, tx.From)
}

// GetSigHash get the transaction data used for signing
func (tx *Tx) GetSigHash() []byte {
	clone := *tx
	clone.From.Signature = nil
	bz, _ := cdc.EncodeToBytes(clone)
	return bz
}

func (tx *Tx) sigVerified() bool {
	signature := tx.From.Signature
	log.Printf("Signature %X\n", signature)
	return ed25519.Verify(tx.From.PubKey, tx.GetSigHash(), signature)
}

// Serialize a Tx to bytes
func (tx *Tx) Serialize() []byte {
	bytes, _ := cdc.EncodeToBytes(tx)
	return bytes
}

// Serialize a TxData to bytes
func (txData *TxData) Serialize() []byte {
	var bytes []byte
	nameBytes := make([]byte, MethodNameByteLength)
	copy(nameBytes[:], txData.Method)
	bytes = append(bytes, nameBytes...)
	bytes = append(bytes, txData.Params...)
	return bytes
}

// Serialize a TxData to bytes
func (txSigner *TxSigner) Serialize() []byte {
	bytes, _ := cdc.EncodeToBytes(txSigner)
	return bytes
}

// Deserialize converts bytes to Tx
func (tx *Tx) Deserialize(bz []byte) {
	cdc.DecodeBytes(bz, &tx)
}

// Deserialize converts bytes to TxData
func (txData *TxData) Deserialize(bz []byte) {
	txData.Method = string(bytes.Trim(bz[0:64], "\x00"))
	txData.Params = bz[64:]
}

// Deserialize converts bytes to txSigner
func (txSigner *TxSigner) Deserialize(bz []byte) {
	cdc.DecodeBytes(bz, &txSigner)
}
