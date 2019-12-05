package main

import (
	"os"
	"path/filepath"

	"github.com/QuoineFinancial/vertex/cmd/node"
)

func main() {
	rootDir := os.ExpandEnv(filepath.Join("$HOME", ".vertex"))

	// TODO: Get gasContractAddress from genesis file
	gasContractAddress := os.Getenv("GAS_CONTRACT_ADDRESS")
	vertexNode := node.New(rootDir, gasContractAddress)
	vertexNode.Execute()
}
