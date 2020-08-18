package storage

import (
	"encoding/binary"

	"github.com/QuoineFinancial/liquid-chain/common"
)

// metaKeyPrefix is type of prefix keys for indexing
type metaKeyPrefix byte

const (
	blockHeightToBlockHashPrefix metaKeyPrefix = 0x0
	txHashToBlockHeightPrefix    metaKeyPrefix = 0x1
	latestBlockHeightPrefix      metaKeyPrefix = 0x2
)

func (index *MetaStorage) encodeTxHashToBlockHeightKey(hash common.Hash) []byte {
	return index.encodeKey(txHashToBlockHeightPrefix, hash[:])
}

func (index *MetaStorage) encodeBlockHeightToBlockHashKey(height uint64) []byte {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, height)
	return index.encodeKey(blockHeightToBlockHashPrefix, key)
}

func (index *MetaStorage) encodeLatestBlockHeightKey() []byte {
	return index.encodeKey(latestBlockHeightPrefix, []byte{})
}

func (index *MetaStorage) encodeKey(prefix metaKeyPrefix, key []byte) []byte {
	return append([]byte{byte(prefix)}, key...)
}
