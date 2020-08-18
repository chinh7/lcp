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
	Index    uint32
	Data     []byte
}

// TxReceipt reflects corresponding Transaction execution result
type TxReceipt struct {
	Result  uint64
	GasUsed uint32
	Code    ReceiptCode
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
	Receiver  Address
	Payload   *TxPayload
	GasPrice  uint32
	GasLimit  uint32
	Signature []byte
	Receipt   *TxReceipt
}

// Encode returns bytes representation of transaction
func (tx Transaction) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(tx)
}

// DecodeTransaction returns Transaction from bytes representation
func DecodeTransaction(raw []byte) (*Transaction, error) {
	var tx Transaction
	if err := rlp.DecodeBytes(raw, &tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

// Hash returns hash for storing transaction
func (tx Transaction) Hash() common.Hash {
	hash, _ := tx.Encode()
	return blake2b.Sum256(hash)
}
