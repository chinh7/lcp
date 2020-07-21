package storage

import (
	"time"

	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// TransactionType to clarify type of transaction
type TransactionType uint

const (
	// TransactionTypeCreate indicates that tx is to create new smart contract from compiled data
	TransactionTypeCreate TransactionType = 0x0

	// TransactionTypeInvoke indicates that tx is to trigger existed smart contract on chain
	TransactionTypeInvoke TransactionType = 0x1
)

// Transaction contains
type Transaction struct {
	Version float64
	Type    TransactionType

	Nonce     uint64
	PublicKey []byte
	Signature []byte

	To       crypto.Address
	Data     string
	Memo     string
	GasPrice uint32
	GasLimit uint32
}

// Event is emitted while executing transactions
type Event struct {
	Name     string
	Data     string
	Contract *crypto.Address
}

// TransactionReceipt reflects corresponding Transaction execution result
type TransactionReceipt struct {
	TransactionHash string
	Result          uint64
	GasUsed         uint64
	Events          []*Event
}

// Block contains basic info and root hash of storage, transactions and reciepts
type Block struct {
	CreatedAt            time.Time
	StorageHash          string
	TransactionsRootHash string
	ReceiptsRootHash     string
}
