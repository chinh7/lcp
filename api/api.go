package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

// API contains all info to serve an api server
type API struct {
	address string
	server  *rpc.Server
	router  *mux.Router
}

// NewAPI return an new instance of API
func NewAPI(address string) *API {
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")

	// Register our services here
	server.RegisterService(new(HelloService), "")

	// Set up router
	router := mux.NewRouter()
	router.Handle("/", server)

	return &API{address, server, router}
}

// Serve starts the server to serve request
func (api *API) Serve() {
	log.Printf("API server is ready at %s", api.address)
	http.ListenAndServe(api.address, api.router)
}
