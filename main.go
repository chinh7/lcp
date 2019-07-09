package main

import (
	"fmt"

	"github.com/QuoineFinancial/vertex-storage/db"
)

func main() {
	db := db.NewRocksDB("data")
	db.Put([]byte("Hello"), []byte("World"))
	db.Put([]byte("Dang"), []byte("Nguyen"))
	db.Put([]byte("Kha"), []byte("Do"))
	fmt.Println(string(db.Get([]byte("Dang"))))
	fmt.Println(string(db.Get([]byte("Kha"))))
	fmt.Println(string(db.Get([]byte("Hello"))))
}
