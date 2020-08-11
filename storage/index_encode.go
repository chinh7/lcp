package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/QuoineFinancial/liquid-chain/common"
)

// IndexPrefix is type of prefix keys for indexing
type IndexPrefix string

// IndexPrefix values
const (
	HeightToBlockHashPrefix IndexPrefix = "h2bh:"
	TxHashToHeightPrefix    IndexPrefix = "txh2h:"
)

func (index *IndexStorage) encodeTxHashToHeightKey(hash common.Hash) []byte {
	return index.encodeKey(TxHashToHeightPrefix, hash[:])
}

func (index *IndexStorage) encodeHeightToBlockHashKey(height uint64) []byte {
	var key []byte
	binary.LittleEndian.PutUint64(key, height)
	return index.encodeKey(HeightToBlockHashPrefix, key)
}

func (index *IndexStorage) encodeKey(prefix IndexPrefix, key []byte) []byte {
	return append([]byte(prefix), key...)
}

func (index *IndexStorage) decodeKey(prefix IndexPrefix, encodedKey []byte) ([]byte, error) {
	if !bytes.Equal([]byte(prefix), encodedKey[:len(prefix)]) {
		return nil, fmt.Errorf("Prefix %s not matched with key %v", prefix, encodedKey)
	}

	return encodedKey[len(prefix):], nil
}
