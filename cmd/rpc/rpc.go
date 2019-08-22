package main

import (
	"net/http"
	"os"

	amino "github.com/tendermint/go-amino"
	command "github.com/tendermint/tendermint/cmd/tendermint/commands"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/lite/proxy"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcserver "github.com/tendermint/tendermint/rpc/lib/server"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
)

// RPCRoutes just routes everything to the given client, as if it were
// a tendermint fullnode.
//
// if we want security, the client must implement it as a secure client
func RPCRoutes(c rpcclient.Client) map[string]*rpcserver.RPCFunc {
	routes := proxy.RPCRoutes(c)
	// we can add new endpoints here
	routes["status"] = rpcserver.NewRPCFunc(makeStatusFunc(c), "")
	return routes
}

func makeStatusFunc(c rpcclient.Client) func(ctx *rpctypes.Context) (*ctypes.ResultStatus, error) {
	return func(ctx *rpctypes.Context) (*ctypes.ResultStatus, error) {
		return c.Status()
	}
}

func main() {
	var (
		mux    = http.NewServeMux()
		cdc    = amino.NewCodec()
		logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	)
	ctypes.RegisterAmino(cdc)

	// vars on connection
	nodeAddr := "tcp://localhost:26657"
	listenAddr := "tcp://localhost:8008"
	chainID := "test-chain-GSukkq"
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

	routes := RPCRoutes(sc)
	// Stop upon receiving SIGTERM or CTRL-C.
	cmn.TrapSignal(logger, func() {})

	rpcserver.RegisterRPCFuncs(mux, routes, cdc, logger)
	config := rpcserver.DefaultConfig()
	listener, err := rpcserver.Listen(listenAddr, config)
	if err != nil {
		cmn.Exit(err.Error())
	}
	rpcserver.StartHTTPServer(listener, mux, logger, config)
}
