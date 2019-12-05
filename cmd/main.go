package main

import (
	"os"
	"path/filepath"

	"github.com/QuoineFinancial/vertex/cmd/node"
	"github.com/tendermint/tendermint/config"
)

func main() {
	rootDir := os.ExpandEnv(filepath.Join("$HOME", ".vertex"))
	config := config.DefaultConfig()
	config.SetRoot(rootDir)
	vertexNode := node.New(*config)
	vertexNode.Execute()
}
