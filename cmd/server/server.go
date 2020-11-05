package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/projecteru2/pistage/cmd/server/act"
	"github.com/projecteru2/pistage/cmd/server/daemon"
	"github.com/projecteru2/pistage/errors"
	"github.com/projecteru2/pistage/ver"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(ver.Version())
	}

	app := cli.App{
		Commands: []*cli.Command{
			act.Command(),
			daemon.Command(),
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
		&cli.StringSliceFlag{
			Name:     "config",
			Usage:    "config files",
			Required: true,
		},
	}
}
