package storage

import (
	"errors"
	"time"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
)

// BlockStorage is storage for block
type BlockStorage struct {
	db.Database
	currentBlock *crypto.Block
}

// NewBlockStorage returns new instance of IndexStorage
func NewBlockStorage(db db.Database) *BlockStorage {
	return &BlockStorage{db, nil}
}

// ComposeBlock compose currentBlock based on parent and proposed time
func (bs *BlockStorage) ComposeBlock(parent *crypto.Block, time time.Time) {
	bs.currentBlock = crypto.NewEmptyBlock(parent.Header.Hash(), parent.Header.Height+1, time)
}

// FinalizeBlock completes currentBlock with committed state
func (bs *BlockStorage) FinalizeBlock(stateRoot, txRoot common.Hash) error {
	if bs.currentBlock == nil {
		return errors.New("BlockStorage.currentBlock is nil")
	}

	bs.currentBlock.Header.SetStateRoot(stateRoot)
	bs.currentBlock.Header.SetTransactionRoot(txRoot)
	return nil
}

// Commit puts currentBlock to storage
func (bs *BlockStorage) Commit() common.Hash {
	if bs.currentBlock == nil {
		panic("BlockStorage.currentBlock is nil")
	}

	hash := bs.currentBlock.Header.Hash()
	rawBlock, err := bs.currentBlock.Encode()
	if err != nil {
		panic(err)
	}
	bs.Put(hash.Bytes(), rawBlock)
	return hash
}

// AddTransaction add tx to currentBlock
func (bs *BlockStorage) AddTransaction(tx *crypto.Transaction) error {
	if bs.currentBlock == nil {
		return errors.New("BlockStorage.currentBlock is nil")
	}

	bs.currentBlock.Transactions = append(bs.currentBlock.Transactions, tx)
	return nil
}

// GetBlock retrieves block by its hash
func (bs *BlockStorage) GetBlock(hash common.Hash) (*crypto.Block, error) {
	if hash == common.EmptyHash {
		return &crypto.GenesisBlock, nil
	}
	rawBlock := bs.Get(hash.Bytes())
	return crypto.DecodeBlock(rawBlock)
}

// MustGetBlock retrieves block by its hash, panic if failed
func (bs *BlockStorage) MustGetBlock(hash common.Hash) *crypto.Block {
	block, err := bs.GetBlock(hash)
	if err != nil {
		panic(err)
	}
	return block
}
