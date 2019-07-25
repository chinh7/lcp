package trie

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/QuoineFinancial/vertex-storage/db"
)

var items = []struct {
	key   string
	value string
}{
	{"do", "verb"},
	{"dog", "puppy"},
	{"doge", "coin"},
	{"horse", "stallion"},
}

var expectedRootHash, _ = hex.DecodeString("5991bb8c6514148a29db676a14ac506cd2cd5775ace63c30a4fe457715e9ac84")

func TestOperation(t *testing.T) {
	// Setup
	path := "./test-trie-operation"
	database := db.NewRocksDB(path)
	tree := New(Hash{}, database)

	// Put
	for _, item := range items {
		tree.Update([]byte(item.key), []byte(item.value))
	}

	// Get
	for _, item := range items {
		actual, err := tree.Get([]byte(item.key))
		if err != nil {
			t.Errorf("Get error %v", err)
		}
		if !bytes.Equal(actual, []byte(item.value)) {
			t.Errorf("Value getting from trie is different from expected. Expected: %v. Actual: %v", item.value, actual)
		}
	}

	// Hash
	actualRootHash := tree.Hash()
	if !bytes.Equal(expectedRootHash, actualRootHash[:]) {
		t.Errorf("Root hash incorrect. Expected: %v. Actual: %v", string(expectedRootHash), string(actualRootHash[:]))
	}

	// Tear down
	os.RemoveAll(path)
}

func TestLoading(t *testing.T) {
	// Setup
	path := "./test-trie-loading"
	database := db.NewRocksDB(path)
	tree := New(Hash{}, database)
	for _, item := range items {
		tree.Update([]byte(item.key), []byte(item.value))
	}
	rootHash := tree.Commit()
	database.GetInstance().Close()

	// Load new trie from old rootHash
	newDatabase := db.NewRocksDB(path)
	newTree := New(rootHash, newDatabase)
	fmt.Println(newTree.Hash())
	for _, item := range items {
		actual, err := newTree.Get([]byte(item.key))
		if err != nil {
			t.Errorf("Get error %v", err)
		}
		if !bytes.Equal(actual, []byte(item.value)) {
			t.Errorf("Value getting from trie is different from expected. Expected: %v. Actual: %v", item.value, actual)
		}
	}

	// Hash
	actualRootHash := newTree.Hash()
	if !bytes.Equal(expectedRootHash, actualRootHash[:]) {
		t.Errorf("Root hash incorrect. Expected: %v. Actual: %v", string(expectedRootHash), string(actualRootHash[:]))
	}

	// Tear down
	os.RemoveAll(path)
}
