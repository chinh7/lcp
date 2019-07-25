package crypto

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/tendermint/go-amino"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"
)

// TxData data for contract deploy/invoke
type TxData struct {
	Method string
	Params []interface{}
}

// TxSigner information about transaction signer
type TxSigner struct {
	PubKey    []byte
	Nonce     uint64
	Signature []byte
}

// Tx transaction
type Tx struct {
	From TxSigner
	Data []byte
	To   Address
}

var cdc = amino.NewCodec()

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
	bz, _ := cdc.MarshalBinaryLengthPrefixed(clone)
	return bz
}

func (tx *Tx) sigVerified() bool {
	signature := tx.From.Signature
	log.Printf("Signature %X\n", signature)
	return ed25519.Verify(tx.From.PubKey, tx.GetSigHash(), signature)
}

// Serialize a Tx to bytes
func (tx *Tx) Serialize() []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(tx)
}

// Serialize a TxData to bytes
func (txData *TxData) Serialize() []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(txData)
}

// Serialize a TxData to bytes
func (txSigner *TxSigner) Serialize() []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(txSigner)
}

// Deserialize converts bytes to Tx
func (tx *Tx) Deserialize(bz []byte) {
	cdc.MustUnmarshalBinaryLengthPrefixed(bz, tx)
}

// Deserialize converts bytes to TxData
func (txData *TxData) Deserialize(bz []byte) {
	cdc.MustUnmarshalBinaryLengthPrefixed(bz, txData)
}

// RegisterCodec registers types that need encoding to the animo codec
func RegisterCodec() {
	log.Println("Registering Codec")
	cdc.RegisterConcrete(&Tx{}, "Tx", nil)
	cdc.RegisterConcrete(&TxSigner{}, "TxSigner", nil)
	cdc.RegisterConcrete(&TxData{}, "TxData", nil)
	cdc.RegisterInterface((*interface{})(nil), nil)
	cdc.RegisterConcrete(int64(0), "int", nil)
	cdc.RegisterConcrete(string(""), "string", nil)
}
