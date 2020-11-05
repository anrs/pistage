package daemon

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // nolint
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/projecteru2/pistage/cmd/server/run"
	"github.com/projecteru2/pistage/config"
	"github.com/projecteru2/pistage/errors"
	"github.com/projecteru2/pistage/log"
	"github.com/projecteru2/pistage/metrics"
	"github.com/projecteru2/pistage/pistaged"
	"github.com/projecteru2/pistage/pistaged/server"
)

var signals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGQUIT,
	syscall.SIGUSR2,
}

// Command .
func Command() *cli.Command {
	return &cli.Command{
		Name:   "server",
		Action: run.Run(daemon),
	}
}

func daemon(c *cli.Context) error {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	go prof()

	dump, err := config.Conf.Dump()
	if err != nil {
		return errors.Trace(err)
	}
	log.Infof(dump)

	srv, err := server.Listen(config.Conf.BindGRPCAddr)
	if err != nil {
		return errors.Trace(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleSignals(srv)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Serve(); err != nil {
			log.ErrorStack(err)
			metrics.IncrError()
		}
	}()

	log.Infof("[server] pistaged is running")
	wg.Wait()

	log.Infof("[server] pistaged is terminated")
	return nil
}

func handleSignals(srv pistaged.Server) {
	defer func() {
		log.Warnf("[server] signals handler %p exit", srv)
		srv.Close()
	}()

	sch := make(chan os.Signal, 1)
	signal.Notify(sch, signals...)

	for {
		select {
		case sign := <-sch:
			switch sign {
			case syscall.SIGHUP, syscall.SIGUSR2:
				log.Warnf("[server] recv signal %d to reload", sign)
				log.ErrorStack(srv.Reload())

			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Warnf("[server] recv signal %d to exit", sign)
				return

			default:
				log.Warnf("[server] recv signal %d to ignore", sign)
			}

		case <-srv.Exit():
			log.Warnf("[server] recv from server %p exit ch", srv)
			return
		}
	}
}

func prof() {
	switch flag := strings.ToLower(os.Getenv("PISTAGED_PPROF")); flag {
	case "", "0", "false", "off":
		return
	}

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.Conf.ProfHTTPPort), nil) // nolint
}
