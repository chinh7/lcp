package main

import (
	"os"

	"github.com/QuoineFinancial/vertex/api"
)

func main() {
	apiServer := api.NewAPI(":5555", api.Config{
		ChainID:     os.Getenv("CHAIN_ID"),
		HomeDir:     os.Getenv("HOME_DIR"),
		NodeAddress: os.Getenv("NODE_ADDRESS"),
	})
	apiServer.Serve()
}
