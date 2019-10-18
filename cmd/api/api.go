package main

import (
	"os"

	"github.com/QuoineFinancial/vertex/api"
)

func main() {
	apiServer := api.NewAPI(":5555", api.Config{
		HomeDir: os.Getenv("HOME_DIR"),
		NodeURL: os.Getenv("NODE_ADDRESS"),
	})
	apiServer.Serve()
}
