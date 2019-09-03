package api

import "net/http"

type helloArgs struct {
	Who string
}

type helloReply struct {
	Message string
}

// HelloService is first service
type HelloService struct{}

// Say is handler of HelloService
func (h *HelloService) Say(r *http.Request, args *helloArgs, reply *helloReply) error {
	reply.Message = "Hello, " + args.Who + "!"
	return nil
}
