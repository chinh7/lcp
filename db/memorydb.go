package db

import (
	"encoding/hex"
	"fmt"
)

// MemoryDB simple memory database
type MemoryDB struct {
	cache map[string][]byte
}

// NewMemoryDB return new in-memory database
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		cache: make(map[string][]byte),
	}
}

// Get returns the value based on key
func (db *MemoryDB) Get(key []byte) []byte {
	fmt.Printf("GET: %x\n", key)
	return db.cache[hex.EncodeToString(key)]
}

// Put inserts an key-value pair to database
func (db *MemoryDB) Put(key []byte, value []byte) {
	db.cache[hex.EncodeToString(key)] = value
}

// Delete removes a key from database
func (db *MemoryDB) Delete(key []byte) {
	fmt.Printf("DELETE: %x\n", key)
	delete(db.cache, hex.EncodeToString(key))
}
