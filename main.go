package main

import (
	"os"

	"github.com/QuoineFinancial/vertex/api"
	command "github.com/tendermint/tendermint/cmd/tendermint/commands"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/lite/proxy"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func main() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	nodeAddr := "tcp://localhost:26657"
	chainID := "test-chain-JyjPVo"
	home := ".tendermint-lite"
	cacheSize := 10

	nodeAddr, err := command.EnsureAddrHasSchemeOrDefaultToTCP(nodeAddr)
	if err != nil {
		cmn.Exit(err.Error())
	}

	node := rpcclient.NewHTTP(nodeAddr, "/websocket")

	cert, err := proxy.NewVerifier(chainID, home, node, logger, cacheSize)
	if err != nil {
		cmn.Exit(err.Error())
	}
	cert.SetLogger(logger)
	sc := proxy.SecureClient(node, cert)
	apiServer := api.NewAPI(":8008", sc)
	apiServer.Serve()
}
