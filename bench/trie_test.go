package bench

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/QuoineFinancial/vertex/db"
	"github.com/QuoineFinancial/vertex/trie"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type Node struct {
	key   []byte
	value []byte
}

const nodeCount = 1000000
const keyLength = 32
const valueLength = 128

var nodes []Node

func randomBytes(n int) []byte {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return bytes
}

func init() {
	for i := 0; i < nodeCount; i++ {
		key := randomBytes(keyLength)
		value := randomBytes(valueLength)
		nodes = append(nodes, Node{key, value})
	}
}

func benchmarkInsert(n int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		db := db.NewMemoryDB()
		root := common.HexToHash("")
		b.ReportAllocs()
		tree := trie.New(root, db)
		for i := 0; i < n; i++ {
			tree.Update(nodes[i].key, nodes[i].value)
		}
		tree.Commit()
	}
}

func benchmarkInsertDisk(n int, b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		id, _ := uuid.NewUUID()
		path := fmt.Sprintf("./data-"+id.String(), n, i)
		database := db.NewRocksDB(path)
		root := common.HexToHash("")
		tree := trie.New(root, database)

		for i := 0; i < n; i++ {
			tree.Update(nodes[i].key, nodes[i].value)
		}
		tree.Commit()
		os.RemoveAll(path)
	}

}

// Memory
func BenchmarkInsert1(b *testing.B)       { benchmarkInsert(1, b) }
func BenchmarkInsert100(b *testing.B)     { benchmarkInsert(100, b) }
func BenchmarkInsert10000(b *testing.B)   { benchmarkInsert(10000, b) }
func BenchmarkInsert1000000(b *testing.B) { benchmarkInsert(1000000, b) }

// Disk
func BenchmarkInsertDisk1(b *testing.B)       { benchmarkInsertDisk(1, b) }
func BenchmarkInsertDisk100(b *testing.B)     { benchmarkInsertDisk(100, b) }
func BenchmarkInsertDisk10000(b *testing.B)   { benchmarkInsertDisk(10000, b) }
func BenchmarkInsertDisk1000000(b *testing.B) { benchmarkInsertDisk(1000000, b) }
