package server

import (
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/projecteru2/pistage/config"
	pb "github.com/projecteru2/pistage/grpc/gen"
	"github.com/projecteru2/pistage/log"
	"github.com/projecteru2/pistage/netx"
)

// Server .
type Server struct {
	Listener net.Listener
	IP       string

	grpc *grpc.Server
	app  *app

	exit struct {
		sync.Once
		C chan struct{}
	}
}

// Listen .
func Listen(addr string) (srv *Server, err error) {
	network := "tcp"

	srv = &Server{}
	if srv.Listener, err = net.Listen(network, addr); err != nil {
		return
	}

	srv.grpc = grpc.NewServer()
	srv.app = &app{}
	srv.exit.C = make(chan struct{}, 1)
	srv.IP, err = netx.GetLocalIP(network, srv.Listener.Addr().String())

	return
}

// Serve .
func (s *Server) Serve() error {
	defer func() {
		s.Close()
		log.Warnf("[server] GRPC server main loop %p is terminated", s)
	}()

	pb.RegisterPistagedServer(s.grpc, s.app)

	return s.grpc.Serve(s.Listener)
}

// Close .
func (s *Server) Close() {
	s.exit.Do(func() {
		close(s.exit.C)

		done := make(chan struct{})
		go func() {
			defer close(done)
			s.grpc.GracefulStop()
		}()

		timer := time.NewTimer(config.Conf.GracefulTimeout.Duration())
		select {
		case <-done:
			log.Infof("[server] terminating GRPC server gracefully")
		case <-timer.C:
			log.Warnf("[server] terminating GRPC server forcefully")
			s.grpc.Stop()
		}
	})
}

// Reload .
func (s *Server) Reload() error {
	return nil
}

// Exit .
func (s *Server) Exit() <-chan struct{} {
	return s.exit.C
}
