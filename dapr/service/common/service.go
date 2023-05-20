package common

import (
	"context"

	"google.golang.org/grpc"
)

type (
	HealthCheckHandler func(context.Context) error
)

// Service represents Dapr callback service.
type Service interface {
	GrpcServer() *grpc.Server
	// AddHealthCheckHandler sets a health check handler, name: http (router) and grpc (invalid).
	AddHealthCheckHandler(name string, fn HealthCheckHandler) error
	// Start starts service.
	Start() error
	// Stop stops the previously started service.
	Stop() error
	// Gracefully stops the previous started service
	GracefulStop() error
}
