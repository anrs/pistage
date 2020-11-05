package server

import (
	"context"

	pb "github.com/projecteru2/pistage/grpc/gen"
	"github.com/projecteru2/pistage/log"
	"github.com/projecteru2/pistage/ver"
)

type app struct {
}

// Ping .
func (a *app) Ping(_ context.Context, _ *pb.Empty) (*pb.Pong, error) {
	log.Infof("[server] recv Ping request")
	return &pb.Pong{Version: ver.Version()}, nil
}
