package run

import (
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/pistage/config"
	"github.com/projecteru2/pistage/errors"
	"github.com/projecteru2/pistage/executor"
	"github.com/projecteru2/pistage/log"
	"github.com/projecteru2/pistage/metrics"
	"github.com/projecteru2/pistage/store"
)

// Run .
func Run(fn cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		filepaths := c.StringSlice("config")
		if err := setup(filepaths); err != nil {
			return errors.Trace(err)
		}

		if err := log.Setup(config.Conf.LogLevel, config.Conf.LogFile); err != nil {
			return errors.Trace(err)
		}

		if err := fn(c); err != nil {
			log.ErrorStack(err)
			metrics.IncrError()
			return errors.Trace(err)
		}

		return nil
	}
}

func setup(filepaths []string) error {
	if err := config.Conf.ParseFiles(filepaths...); err != nil {
		return errors.Trace(err)
	}

	return store.Setup(config.Conf.MetaType)
}

// NewExecutor .
func NewExecutor(c *cli.Context) (executor.Executor, error) {
	return executor.NewSimple()
}
