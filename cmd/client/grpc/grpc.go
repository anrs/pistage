package grpc

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/projecteru2/pistage/errors"
	"github.com/projecteru2/pistage/grpc/client"
)

// Run .
func Run(fn func(*cli.Context, *client.Client) error) cli.ActionFunc {
	return func(c *cli.Context) error {
		grpcClient, err := client.New(context.Background(), c.String("server"))
		if err != nil {
			return errors.Trace(err)
		}

		if err = fn(c, grpcClient); err != nil {
			fmt.Println(errors.Stack(err))
		}

		return err
	}
}
