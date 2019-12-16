package main

import (
	"os"
	"path/filepath"

	"github.com/QuoineFinancial/liquid-chain/cmd/node"
)

func main() {
	rootDir := os.ExpandEnv(filepath.Join("$HOME", ".liquid"))

	// TODO: Get gasContractAddress from genesis file
	gasContractAddress := os.Getenv("GAS_CONTRACT_ADDRESS")
	liquidNode := node.NewNode(rootDir, gasContractAddress)
	liquidNode.Execute()
}
