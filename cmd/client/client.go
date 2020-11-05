package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/projecteru2/pistage/cmd/client/ping"
	"github.com/projecteru2/pistage/errors"
	"github.com/projecteru2/pistage/ver"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(ver.Version())
	}

	app := cli.App{
		Commands: []*cli.Command{
			ping.Command(),
		},
		Flags:   globalFlags(),
		Version: ver.Version(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(errors.Stack(err))
	}
}

func globalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "server",
			Usage: "server address, default is 127.0.0.1:8697",
			Value: "127.0.0.1:8697",
		},
	}
}
