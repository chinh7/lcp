package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

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
	api := &API{Client: &c, Address: address}

	// Register server
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")

	// Register our services here
	server.RegisterService(api.NewStatusService(), "")

	// Set up router
	router := mux.NewRouter()
	router.Handle("/rpc", server).Methods("POST")

	api.Server = server
	api.Router = router

	return api
}

// Serve starts the server to serve request
func (api *API) Serve() {
	fmt.Println("Server is ready at", api.Address)
	err := http.ListenAndServe(api.Address, api.Router)
	if err != nil {
		panic(err)
	}
}
