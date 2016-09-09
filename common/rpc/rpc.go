package rpc

import (
	"net"
	"net/rpc"
)

var (
	defaultServer = rpc.NewServer()
)

type VoidType struct{}

type ArgType struct {
	Args map[string]interface{}
}

func NewArgType() *ArgType {
	return &ArgType{
		Args: make(map[string]interface{}),
	}
}

type ReplyType struct {
	Reply map[string]interface{}
}

func NewReplyType() *ReplyType {
	return &ReplyType{
		Reply: make(map[string]interface{}),
	}
}

func Register(obj interface{}) {
	if err := defaultServer.Register(obj); err != nil {
		panic(err.Error())
	}
}

func Serve(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err.Error())
	}

	defaultServer.Accept(listener)
}

func Dail(addr string) (*rpc.Client, error) {
	return rpc.Dial("tcp", addr)
}
