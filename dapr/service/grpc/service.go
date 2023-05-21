package grpc

import (
	"errors"
	"fmt"
	"net"
	"sync/atomic"

	pb "github.com/dapr/go-sdk/dapr/proto/runtime/v1"
	"google.golang.org/grpc"

	"github.com/releaseband/go-micro-tools/dapr/service/common"
)

// NewService creates new Service.
func NewService(address string, opts ...grpc.ServerOption) (s common.Service, err error) {
	if address == "" {
		return nil, errors.New("empty address")
	}
	lis, err := net.Listen("tcp", address)
	if err != nil {
		err = fmt.Errorf("failed to TCP listen on %s: %w", address, err)
		return
	}
	s = newService(lis, opts...)
	return
}

// NewServiceWithListener creates new Service with specific listener.
func NewServiceWithListener(lis net.Listener, opts ...grpc.ServerOption) common.Service {
	return newService(lis, opts...)
}

func newService(lis net.Listener, opts ...grpc.ServerOption) *Server {
	s := &Server{
		listener: lis,
	}

	gs := grpc.NewServer(opts...)
	pb.RegisterAppCallbackHealthCheckServer(gs, s)
	s.grpcServer = gs

	return s
}

// Server is the gRPC service implementation for Dapr.
type Server struct {
	pb.UnimplementedAppCallbackHealthCheckServer
	listener           net.Listener
	healthCheckHandler common.HealthCheckHandler
	grpcServer         *grpc.Server
	started            uint32
}

// Start registers the server and starts it.
func (s *Server) Start() error {
	if !atomic.CompareAndSwapUint32(&s.started, 0, 1) {
		return errors.New("a gRPC server can only be started once")
	}
	return s.grpcServer.Serve(s.listener)
}

// Stop stops the previously-started service.
func (s *Server) Stop() error {
	if atomic.LoadUint32(&s.started) == 0 {
		return nil
	}
	s.grpcServer.Stop()
	s.grpcServer = nil
	return nil
}

// GracefulStop stops the previously-started service gracefully.
func (s *Server) GracefulStop() error {
	if atomic.LoadUint32(&s.started) == 0 {
		return nil
	}
	s.grpcServer.GracefulStop()
	s.grpcServer = nil
	return nil
}

// GrpcServer returns the grpc.Server object managed by the server.
func (s *Server) GrpcServer() *grpc.Server {
	return s.grpcServer
}
