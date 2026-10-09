package core_grpc_server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	grpcServer *grpc.Server
	config     Config
}

func NewGRPCServer(config Config, opts ...grpc.ServerOption) *GRPCServer {
	return &GRPCServer{
		grpcServer: grpc.NewServer(opts...),
		config:     config,
	}
}

func (s *GRPCServer) Register(registerFns ...func(*grpc.Server)) {
	for _, fn := range registerFns {
		fn(s.grpcServer)
	}
}

func (s *GRPCServer) Start(ctx context.Context) error {
	reflection.Register(s.grpcServer)
	lis, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.config.Addr, err)
	}

	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	if err := s.grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

func (s *GRPCServer) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}
