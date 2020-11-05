package ping

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/projecteru2/pistage/cmd/client/grpc"
	"github.com/projecteru2/pistage/errors"
	"github.com/projecteru2/pistage/grpc/client"
	pb "github.com/projecteru2/pistage/grpc/gen"
)

// Command .
func Command() *cli.Command {
	return &cli.Command{
		Name:   "ping",
		Action: grpc.Run(ping),
	}
}

func ping(c *cli.Context, grpcClient *client.Client) error {
	pong, err := grpcClient.Ping(context.Background(), &pb.Empty{})
	if err != nil {
		return errors.Trace(err)
	}

	fmt.Printf("server version: %s\n", pong.Version)

	return nil
}
