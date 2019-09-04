package main

import (
	"fmt"
	"net/http"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type HelloArgs struct {
	Who string
}

type HelloReply struct {
	Message string
}

// HelloService is first service
type HelloService struct {
	c *rpcclient.Client
}

func (api *API) NewHelloService() *HelloService {
	if api.Client == nil {
		fmt.Println("error")
		panic("api.NewHelloService call without api.c")
	}
	return &HelloService{api.Client}
}

// Say is handler of HelloService
func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
	reply.Message = "Hello, " + args.Who + "!"
	client := *h.c
	fmt.Println(client.Status())
	return nil
}
