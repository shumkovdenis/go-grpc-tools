package grpc

import (
	"context"
	"fmt"

	pb "github.com/dapr/go-sdk/dapr/proto/runtime/v1"
	"github.com/shumkovdenis/go-micro-tools/dapr/service/common"

	"google.golang.org/protobuf/types/known/emptypb"
)

// AddHealthCheckHandler appends provided app health check handler.
func (s *Server) AddHealthCheckHandler(_ string, fn common.HealthCheckHandler) error {
	if fn == nil {
		return fmt.Errorf("health check handler required")
	}

	s.healthCheckHandler = fn

	return nil
}

// HealthCheck check app health status.
func (s *Server) HealthCheck(ctx context.Context, _ *emptypb.Empty) (*pb.HealthCheckResponse, error) {
	if s.healthCheckHandler != nil {
		if err := s.healthCheckHandler(ctx); err != nil {
			return &pb.HealthCheckResponse{}, err
		}

		return &pb.HealthCheckResponse{}, nil
	}

	return nil, fmt.Errorf("health check handler not implemented")
}
