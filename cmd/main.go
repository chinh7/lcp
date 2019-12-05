package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/QuoineFinancial/vertex/cmd/node"
	"github.com/tendermint/tendermint/config"
)

func main() {
	rootDir := os.ExpandEnv(filepath.Join("$HOME", ".vertex"))
	config := config.DefaultConfig()
	fmt.Println(rootDir)
	config.SetRoot(rootDir)

	// TODO: Get gasContractAddress from genesis file
	gasContractAddress := os.Getenv("GAS_CONTRACT_ADDRESS")
	vertexNode := node.New(*config, gasContractAddress)
	vertexNode.Execute()
}
