package main

import (
	"fmt"
	"net/http"

	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

	command "github.com/tendermint/tendermint/cmd/tendermint/commands"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/lite/proxy"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// API contains all info to serve an api server
type API struct {
	Address string
	Server  *rpc.Server
	Router  *mux.Router
	Client  *rpcclient.Client
}

// NewAPI return an new instance of API
func NewAPI(address string, c rpcclient.Client) *API {
	api := &API{Client: &c}

	// Register server
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")

	// Register our services here
	server.RegisterService((*api).NewHelloService(), "")

	// Set up router
	router := mux.NewRouter()
	router.Handle("/", server)

	api.Server = server
	api.Router = router

	return api
}

// Serve starts the server to serve request
func (api *API) Serve() {
	err := http.ListenAndServe(api.Address, api.Router)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	var (
		logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	)

	nodeAddr := "tcp://localhost:26657"
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
	api := NewAPI("localhost:8008", sc)
	api.Serve()
}
