package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/rs/cors"

	"github.com/QuoineFinancial/liquid-chain/api/chain"
	"github.com/QuoineFinancial/liquid-chain/api/resource"
	"github.com/QuoineFinancial/liquid-chain/api/storage"
	"github.com/QuoineFinancial/liquid-chain/consensus"
)

// API contains all info to serve an api server
type API struct {
	url    string
	config Config
	srv    *http.Server
	server *rpc.Server
	Router *mux.Router

	tmAPI resource.TendermintAPI
	app   *consensus.App
}

// Config to modify the API
type Config struct {
	HomeDir string
	NodeURL string
	App     *consensus.App
}

// NewAPI return an new instance of API
func NewAPI(url string, config Config) *API {
	api := &API{
		url:    url,
		config: config,
		app:    config.App,
	}
	api.setupExternalAPIs()
	api.setupServer()
	api.registerServices()
	api.setupRouter()
	return api
}

func (api *API) setupExternalAPIs() {
	tAPI := resource.NewTendermintAPI(
		api.config.HomeDir,
		api.config.NodeURL,
		0,
	)
	api.tmAPI = tAPI
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
	api.Router = router
}

func (api *API) registerServices() {
	if api.server == nil {
		panic("api.registerServices call without api.server")
	}
	err := api.server.RegisterService(chain.NewService(api.tmAPI, api.app), "chain")
	if err != nil {
		panic(err)
	}
	err = api.server.RegisterService(storage.NewService(api.tmAPI, api.app), "storage")
	if err != nil {
		panic(err)
	}
}

// Serve starts the server to serve request
func (api *API) Serve() error {
	log.Println("Server is ready at", api.url)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "DELETE", "PUT", "GET", "HEAD", "OPTIONS"},
	})
	handler := c.Handler(api.Router)
	// err := http.ListenAndServe(api.url, handler)
	// if err != nil {
	// 	panic(err)
	// }
	// err = http.ListenAndServe(api.url, api.Router)
	// if err != nil {
	// 	panic(err)

	api.srv = &http.Server{Addr: api.url, Handler: handler}
	err := api.srv.ListenAndServe()
	return err
}

// Close will immediately stop the server without waiting for any active connection to complete
// For gracefully shutdown please implement another function and use Server.Shutdown()
func (api *API) Close() {
	log.Println("Closing server")
	if api.srv != nil {
		err := api.srv.Close()
		if err != nil {
			panic(err)
		}
	}
}
