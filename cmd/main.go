package main

import (
	"os"

	"github.com/QuoineFinancial/liquid-chain/cmd/node"
)

func main() {
	rootDir := os.Getenv("DB_DIR")

	// TODO: Get gasContractAddress from genesis file
	gasContractAddress := os.Getenv("GAS_CONTRACT_ADDRESS")
	liquidNode := node.New(rootDir, gasContractAddress)
	liquidNode.Execute()
}
