package crypto

import (
	"crypto/ed25519"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/common"
	"golang.org/x/crypto/blake2b"
)

// TxEvent is emitted while executing transactions
type TxEvent struct {
	Contract Address
	Data     []byte
}

// TxReceipt reflects corresponding Transaction execution result
type TxReceipt struct {
	Result  uint64
	GasUsed uint32
	Success bool
	Error   string
	Events  []*TxEvent
}

// TxSender is sender of transaction
type TxSender struct {
	PublicKey ed25519.PublicKey
	Nonce     uint64
}

// TxPayload contains data to interact with smart contract
type TxPayload struct {
	Contract []byte
	Method   string
	Params   []byte
}

// Transaction is transaction of liquid-chain
type Transaction struct {
	Version   uint16
	Sender    *TxSender
	Receiver  *Address
	Payload   *TxPayload
	GasPrice  uint32
	GasLimit  uint32
	Signature []byte
	Receipt   *TxReceipt
}

// Serialize returns bytes representation of transaction
func (tx *Transaction) Serialize() ([]byte, error) {
	return rlp.EncodeToBytes(tx)
}

// Deserialize returns Transaction from bytes representation
func (tx *Transaction) Deserialize(raw []byte) error {
	return rlp.DecodeBytes(raw, &tx)
}

// Hash returns hash for storing transaction
func (tx *Transaction) Hash() common.Hash {
	hash, _ := tx.Serialize()
	return blake2b.Sum256(hash)
}
