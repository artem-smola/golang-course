package grpcserver

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
)

type Config struct {
	Address string        `yaml:"address" env:"LISTEN_ADDRESS" env-default:"localhost:8081"`
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
}

type Server struct {
	server   *grpc.Server
	listener net.Listener
}

func NewServer(address string, opts ...grpc.ServerOption) (*Server, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Server{
		server:   grpc.NewServer(opts...),
		listener: lis,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.server.GracefulStop()
	}()
	return s.server.Serve(s.listener)
}

func (s *Server) GRPC() *grpc.Server {
	return s.server
}
