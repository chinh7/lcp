package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"

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
	GasPrice uint64
}

// Address derived from TxSigner PubKey
func (txSigner *TxSigner) Address() Address {
	return AddressFromPubKey(txSigner.PubKey)
}

// CreateAddress create a new contract address based on pubkey and nonce
func (txSigner *TxSigner) CreateAddress() Address {
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
	return fmt.Sprintf("Data: %s To: %s Signer: %s", hex.EncodeToString(tx.Data), tx.To.String(), tx.From.String())
}

// GetSigHash get the transaction data used for signing
func (tx *Tx) GetSigHash() ([]byte, error) {
	clone := *tx
	clone.From.Signature = nil
	return cdc.EncodeToBytes(clone)
}

// GetFee gets transaction fee limit
func (tx *Tx) GetFee() (uint64, error) {
	if tx.GasLimit == 0 || tx.GasPrice == 0 {
		return 0, nil
	}

	fee := tx.GasLimit * tx.GasPrice
	if fee/tx.GasLimit == tx.GasPrice {
		return fee, nil
	}

	return 0, errors.New("fee overflow")
}

// Sign signs a transaction provided a private key
func (tx *Tx) Sign(privKey ed25519.PrivateKey) error {
	pubKey := make([]byte, ed25519.PublicKeySize)
	copy(pubKey, privKey[32:])
	tx.From.PubKey = pubKey

	sigHash, err := tx.GetSigHash()
	if err != nil {
		return err
	}
	tx.From.Signature = ed25519.Sign(privKey, sigHash)
	return nil
}

// SigVerified checks if transaction signature is correct
func (tx *Tx) SigVerified() bool {
	signature := tx.From.Signature
	sigHash, err := tx.GetSigHash()
	if err != nil {
		return false
	}
	return ed25519.Verify(tx.From.PubKey, sigHash, signature)
}

// Serialize a Tx to bytes
func (tx *Tx) Serialize() []byte {
	bytes, _ := cdc.EncodeToBytes(tx)
	return bytes
}

// Serialize a TxData to bytes
func (txData *TxData) Serialize() []byte {
	bytes, _ := cdc.EncodeToBytes(txData)
	return bytes
}

// Serialize a TxData to bytes
func (txSigner *TxSigner) Serialize() []byte {
	bytes, _ := cdc.EncodeToBytes(txSigner)
	return bytes
}

// Deserialize converts bytes to Tx
func (tx *Tx) Deserialize(bz []byte) error {
	if err := cdc.DecodeBytes(bz, &tx); err != nil {
		return err
	}
	return nil
}

// Deserialize converts bytes to TxData
func (txData *TxData) Deserialize(bz []byte) error {
	return cdc.DecodeBytes(bz, &txData)
}

// Deserialize converts bytes to txSigner
func (txSigner *TxSigner) Deserialize(bz []byte) error {
	return cdc.DecodeBytes(bz, &txSigner)
}
