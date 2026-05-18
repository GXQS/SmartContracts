package rpc

import "context"

type Executor interface {
	Execute(context.Context, []byte) ([]byte, error)
}

type GRPCServer struct {
	exec Executor
}

func NewGRPCServer(exec Executor) *GRPCServer { return &GRPCServer{exec: exec} }
