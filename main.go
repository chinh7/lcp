package main

import (
	"fmt"

	"github.com/QuoineFinancial/vertex/db"
	"github.com/QuoineFinancial/vertex/trie"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	db := db.NewRocksDB("data")

	root := common.HexToHash("0x5991bb8c6514148a29db676a14ac506cd2cd5775ace63c30a4fe457715e9ac84")
	tree := trie.New(root, db)

	// tree := trie.New(trie.Hash{}, db)

	// Update
	tree.Update([]byte("do"), []byte("verb"))
	tree.Update([]byte("dog"), []byte("puppy"))
	tree.Update([]byte("doge"), []byte("coin"))

	// Delete by update to nil
	tree.Update([]byte("hors"), []byte(nil))

	// Get data
	v, _ := tree.Get([]byte("do"))
	v, _ = tree.Get([]byte("do"))
	v, _ = tree.Get([]byte("do"))
	fmt.Println(string(v))

	// Compute hash
	newRootHash := tree.Hash()
	fmt.Println(common.ToHex(newRootHash[:]))
}
