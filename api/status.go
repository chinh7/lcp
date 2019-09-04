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
	LatestBlockHash string `json:"latest_block_hash"`
}

// StatusService is first service
type StatusService struct {
	client *rpcclient.Client
}

// NewStatusService returns new instance of StatusService
func (api *API) NewStatusService() *StatusService {
	if api.Client == nil {
		panic("api.NewStatusService call without api.c")
	}
	return &StatusService{api.Client}
}

// Say is handler of StatusService
func (service *StatusService) Say(r *http.Request, args *StatusArgs, reply *StatusReply) error {
	client := *service.client
	status, err := client.Status()
	if err != nil {
		return err
	}
	reply.LatestBlockHash = status.SyncInfo.LatestBlockHash.String()
	return nil
}
