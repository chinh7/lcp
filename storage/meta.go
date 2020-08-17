package storage

import (
	"encoding/binary"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
)

// MetaStorage is storage of indexes
type MetaStorage struct {
	db.Database
}

// NewMetaStorage returns new instance of IndexStorage
func NewMetaStorage(db db.Database) *MetaStorage {
	return &MetaStorage{db}
}

// StoreBlockIndexes extracts all indexes and store it
func (ms *MetaStorage) StoreBlockIndexes(block *crypto.Block) error {
	ms.Put(
		ms.encodeBlockHeightToBlockHashKey(block.Header.Height),
		block.Header.Hash().Bytes(),
	)

	blockHeightByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(blockHeightByte, block.Header.Height)
	for _, tx := range block.Transactions {
		ms.Put(
			ms.encodeTxHashToBlockHeightKey(tx.Hash()),
			blockHeightByte,
		)
	}

	if block.Header.Height > ms.LatestBlockHeight() {
		ms.Put(
			ms.encodeLatestBlockHeightKey(),
			blockHeightByte,
		)
	}

	return nil
}

// LatestBlockHeight retrieves latest block height
func (ms *MetaStorage) LatestBlockHeight() uint64 {
	blockHeightByte := ms.Get(ms.encodeLatestBlockHeightKey())
	if len(blockHeightByte) == 0 {
		return crypto.GenesisBlock.Header.Height
	}
	return binary.LittleEndian.Uint64(blockHeightByte)
}

// BlockHeightToBlockHash retrieves block hash by its height
func (ms *MetaStorage) BlockHeightToBlockHash(height uint64) common.Hash {
	hash := ms.Get(ms.encodeBlockHeightToBlockHashKey(height))
	return common.BytesToHash(hash)
}

// TxHashToBlockHeight retrieves height of block which contains tx
func (ms *MetaStorage) TxHashToBlockHeight(txHash common.Hash) uint64 {
	blockHeightByte := ms.Get(ms.encodeTxHashToBlockHeightKey(txHash))
	return binary.LittleEndian.Uint64(blockHeightByte)
}
