package storage

import (
	"encoding/binary"

	"github.com/QuoineFinancial/liquid-chain/common"
)

// metaKeyPrefix is type of prefix keys for indexing
type metaKeyPrefix byte

const (
	heightToBlockHashPrefix metaKeyPrefix = 0x0
	txHashToHeightPrefix    metaKeyPrefix = 0x1
	latestBlockHeightPrefix metaKeyPrefix = 0x2
)

func (index *MetaStorage) encodeTxHashToHeightKey(hash common.Hash) []byte {
	return index.encodeKey(txHashToHeightPrefix, hash[:])
}

func (index *MetaStorage) encodeHeightToBlockHashKey(height uint64) []byte {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, height)
	return index.encodeKey(heightToBlockHashPrefix, key)
}

func (index *MetaStorage) encodeLatestBlockHeightKey() []byte {
	return index.encodeKey(latestBlockHeightPrefix, []byte{})
}

func (index *MetaStorage) encodeKey(prefix metaKeyPrefix, key []byte) []byte {
	return append([]byte{byte(prefix)}, key...)
}
