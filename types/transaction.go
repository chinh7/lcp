package types

import (
	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// Transaction contains
type Transaction struct {
	Version float64

	Nonce     uint64
	PublicKey []byte
	Signature []byte

	To       *crypto.Address
	Data     []byte
	Memo     []byte
	GasPrice uint32
	GasLimit uint32
}

// Serialize returns bytes representation of transaction
func (transaction *Transaction) Serialize() ([]byte, error) {
	return rlp.EncodeToBytes(transaction)
}

// Deserialize returns Transaction from bytes representation
func Deserialize(raw []byte) (*Transaction, error) {
	var transaction Transaction
	if err := rlp.DecodeBytes(raw, &transaction); err != nil {
		return nil, err
	}
	return &transaction, nil
}
