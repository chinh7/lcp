package storage

import (
	"encoding/binary"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
)

// IndexStorage is storage of indexes
type IndexStorage struct {
	database db.Database
}

// NewIndexStorage returns new instance of IndexStorage
func NewIndexStorage(db db.Database) *IndexStorage {
	return &IndexStorage{db}
}

// StoreBlockIndexes extracts all indexes and store it
func (index *IndexStorage) StoreBlockIndexes(block *crypto.Block) error {
	index.database.Put(
		index.encodeHeightToBlockHashKey(block.Header.Height),
		block.Header.Hash().Bytes(),
	)

	var blockHeightByte []byte
	binary.LittleEndian.PutUint64(blockHeightByte, block.Header.Height)
	for _, tx := range block.Transactions {
		index.database.Put(
			index.encodeTxHashToHeightKey(tx.Hash()),
			blockHeightByte,
		)
	}

	return nil
}

// HeightToBlockHash retrieves block hash by its height
func (index *IndexStorage) HeightToBlockHash(height uint64) common.Hash {
	hash := index.database.Get(index.encodeHeightToBlockHashKey(height))
	return common.BytesToHash(hash)
}

// TxHashToBlockHeight retrieves height of block which contains tx
func (index *IndexStorage) TxHashToBlockHeight(txHash common.Hash) uint64 {
	blockHeightByte := index.database.Get(index.encodeTxHashToHeightKey(txHash))
	return binary.LittleEndian.Uint64(blockHeightByte)
}
