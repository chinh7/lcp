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
	"github.com/QuoineFinancial/liquid-chain/db"
)

// API contains all info to serve an api server
type API struct {
	url    string
	config Config
	server *rpc.Server
	router *mux.Router

	tmAPI    resource.TendermintAPI
	database db.Database
}

// Config to modify the API
type Config struct {
	HomeDir string
	NodeURL string
	DB      db.Database
}

// NewAPI return an new instance of API
func NewAPI(url string, config Config) *API {
	api := &API{
		url:      url,
		config:   config,
		database: config.DB,
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
	api.router = router
}

func (api *API) registerServices() {
	if api.server == nil {
		panic("api.registerServices call without api.server")
	}
	err := api.server.RegisterService(chain.NewService(api.tmAPI, api.database), "chain")
	if err != nil {
		panic(err)
	}
	err = api.server.RegisterService(storage.NewService(api.tmAPI, api.database), "storage")
	if err != nil {
		panic(err)
	}
}

// Serve starts the server to serve request
func (api *API) Serve() {
	log.Println("Server is ready at", api.url)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "DELETE", "PUT", "GET", "HEAD", "OPTIONS"},
	})
	handler := c.Handler(api.router)
	err := http.ListenAndServe(api.url, handler)
	if err != nil {
		panic(err)
	}
	err = http.ListenAndServe(api.url, api.router)
	if err != nil {
		panic(err)
	}
}
