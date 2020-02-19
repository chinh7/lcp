package db

import (
	"github.com/tecbot/gorocksdb"
)

// RocksDB use map to store and retrieve value
type RocksDB struct {
	instance *gorocksdb.DB
	cache    map[*[]byte]*[]byte
}

// NewRocksDB returns a new instance of the RocksDB
func NewRocksDB(path string) *RocksDB {
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	instance, err := gorocksdb.OpenDb(opts, path)
	if err != nil {
		panic(err)
	}
	return &RocksDB{
		instance: instance,
		cache:    make(map[*[]byte]*[]byte),
	}
}

// Get returns the value based on key
func (db *RocksDB) Get(key []byte) []byte {
	ro := gorocksdb.NewDefaultReadOptions()
	ro.SetFillCache(true)
	value, err := db.instance.Get(ro, key)
	if err != nil {
		panic(err)
	}
	return value.Data()
}

// Put inserts an key-value pair to database
func (db *RocksDB) Put(key []byte, value []byte) {
	wo := gorocksdb.NewDefaultWriteOptions()
	wo.SetSync(false)
	if err := db.instance.Put(wo, key, value); err != nil {
		panic(err)
	}
}

func (db *RocksDB) Close() {
	db.instance.Close()
}
