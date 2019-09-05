package api

import (
	"net/http"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// StatusArgs is params of StatusService
type StatusArgs struct {
}

// StatusReply is response of StatusService
type StatusReply struct {
	LatestBlockHash   string `json:"latest_block_hash"`
	LatestBlockHeight int64  `json:"latest_block_height"`
	ChainID           string `json:"chain_id"`
}

// StatusService is first service
type StatusService struct {
	client *rpcclient.Client
}

// NewStatusService returns new instance of StatusService
func (api *API) NewStatusService() *StatusService {
	if api.Client == nil {
		panic("api.NewStatusService call without api.Client")
	}
	return &StatusService{api.Client}
}

// Get is handler of StatusService
func (service *StatusService) Get(r *http.Request, args *StatusArgs, reply *StatusReply) error {
	client := *service.client
	status, err := client.Status()
	if err != nil {
		return err
	}
	reply.LatestBlockHash = status.SyncInfo.LatestBlockHash.String()
	reply.ChainID = status.NodeInfo.Network
	reply.LatestBlockHeight = status.SyncInfo.LatestBlockHeight
	return nil
}
