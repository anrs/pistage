package act

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/projecteru2/pistage/action"
	"github.com/projecteru2/pistage/cmd/server/run"
	"github.com/projecteru2/pistage/errors"
	"github.com/projecteru2/pistage/io"
)

// Command .
func Command() *cli.Command {
	return &cli.Command{
		Name: "action",
		Subcommands: []*cli.Command{
			{
				Name:   "register",
				Action: run.Run(register),
			},
			{
				Name:   "execute",
				Action: run.Run(execute),
			},
		},
	}
}

func execute(c *cli.Context) error {
	complex, err := parseComplex(c)
	if err != nil {
		return errors.Trace(err)
	}

	executor, err := run.NewExecutor(c)
	if err != nil {
		return errors.Trace(err)
	}

	_, err = executor.SyncStart(context.Background(), complex)

	return errors.Trace(err)
}

func register(c *cli.Context) error {
	complex, err := parseComplex(c)
	if err != nil {
		return errors.Trace(err)
	}

	if err := complex.Save(context.Background()); err != nil {
		return errors.Trace(err)
	}

	fmt.Printf("%s has been registered\n", complex.Name)

	return nil
}

func parseComplex(c *cli.Context) (*action.Complex, error) {
	cont, err := readSpec(c)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return action.Parse(string(cont))
}

func readSpec(c *cli.Context) (string, error) {
	specFile := c.Args().First()
	if len(specFile) < 1 {
		return "", errors.New("spec filepath is required")
	}

	buf, err := io.ReadFile(specFile)
	if err != nil {
		return "", errors.Trace(err)
	}

	return string(buf), nil
}
