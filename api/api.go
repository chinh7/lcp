package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"

	"github.com/QuoineFinancial/vertex/api/chain"
	"github.com/QuoineFinancial/vertex/api/resource"
	"github.com/QuoineFinancial/vertex/api/storage"
)

// API contains all info to serve an api server
type API struct {
	address string
	config  Config
	server  *rpc.Server
	router  *mux.Router

	// tAPI is client to interact with Tendermint RPC
	tAPI resource.TendermintAPI
}

// Config to modify the API
type Config struct {
	HomeDir     string
	NodeAddress string
}

// NewAPI return an new instance of API
func NewAPI(address string, config Config) *API {
	api := &API{address: address, config: config}
	api.setupExternalAPIs()
	api.setupServer()
	api.registerServices()
	api.setupRouter()
	return api
}

func (api *API) setupExternalAPIs() {
	tAPI := resource.NewTendermintAPI(
		api.config.HomeDir,
		api.config.NodeAddress,
	)
	api.tAPI = tAPI
}

func (api *API) setupServer() {
	server := rpc.NewServer()
	server.RegisterCodec(json2.NewCodec(), "application/json")
	api.server = server
}

func (api *API) setupRouter() {
	if api.server == nil {
		panic("api.setupRouter call without api.server")
	}
	router := mux.NewRouter()
	router.Handle("/", api.server).Methods("POST")
	api.router = router
}

func (api *API) registerServices() {
	if api.server == nil {
		panic("api.registerServices call without api.server")
	}
	api.server.RegisterService(chain.NewService(api.tAPI), "chain")
	api.server.RegisterService(storage.NewService(api.tAPI), "storage")
}

// Serve starts the server to serve request
func (api *API) Serve() {
	fmt.Println("Server is ready at", api.address)
	err := http.ListenAndServe(api.address, api.router)
	if err != nil {
		panic(err)
	}
}
