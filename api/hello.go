package api

import (
	"net/http"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type HelloArgs struct {
	Who string
}

type HelloReply struct {
	LatestBlockHash string `json:"latest_block_hash"`
}

// HelloService is first service
type HelloService struct {
	client *rpcclient.Client
}

func (api *API) NewHelloService() *HelloService {
	if api.Client == nil {
		panic("api.NewHelloService call without api.c")
	}
	return &HelloService{api.Client}
}

// Say is handler of HelloService
func (service *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
	client := *service.client
	status, err := client.Status()
	if err != nil {
		return err
	}
	reply = &HelloReply{
		LatestBlockHash: status.SyncInfo.LatestBlockHash.String(),
	}
	return nil
}
