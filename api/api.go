package api

import (
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

// API contains all info to serve an api server
type API struct {
	server *rpc.Server
}

// NewAPI return an new instance of API
func NewAPI() *API {
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")

	// Register our services here
	server.RegisterService(new(HelloService), "")

	return &API{server}
}
