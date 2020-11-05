package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/projecteru2/core/client/interceptor"

	"github.com/projecteru2/pistage/errors"
	pb "github.com/projecteru2/pistage/grpc/gen"
	"github.com/projecteru2/pistage/log"
)

// Client .
type Client struct {
	pb.PistagedClient
	addr string
	conn *grpc.ClientConn
}

// New .
func New(ctx context.Context, addr string) (*Client, error) {
	cli := &Client{addr: addr}
	if err := cli.connect(ctx); err != nil {
		return nil, errors.Trace(err)
	}
	return cli, nil
}

func (c *Client) connect(ctx context.Context) (err error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: 6 * 60 * time.Second, Timeout: time.Second}),
		grpc.WithUnaryInterceptor(interceptor.NewUnaryRetry(interceptor.RetryOptions{Max: 1})),
		grpc.WithStreamInterceptor(interceptor.NewStreamRetry(interceptor.RetryOptions{Max: 1})),
	}

	if c.conn, err = grpc.Dial(c.addr, opts...); err != nil {
		log.Warnf(errors.Stack(err))
		return
	}

	c.PistagedClient = pb.NewPistagedClient(c.conn)

	return
}
