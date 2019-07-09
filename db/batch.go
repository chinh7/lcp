package db

import (
	"github.com/tecbot/gorocksdb"
)

// Batch is the batch of write operations of RocksDB
type Batch struct {
	db    *RocksDB
	batch *gorocksdb.WriteBatch
}

// NewBatch returns the rocksdb's write batch
func (db *RocksDB) NewBatch() *Batch {
	batch := gorocksdb.NewWriteBatch()
	return &Batch{
		db:    db,
		batch: batch,
	}
}

// Put inserts a key-value pair to the batch
func (b *Batch) Put(key []byte, value []byte) {
	b.batch.Put(key, value)
}

// Delete removes a key from the batch
func (b *Batch) Delete(key []byte) {
	b.batch.Delete(key)
}

// Write applies the changes of all operations in batch to the database
func (b *Batch) Write() {
	b.db.instance.Write(gorocksdb.NewDefaultWriteOptions(), b.batch)
}

