package gohookd

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

func NewGohookdGRPCServer(ctx context.Context, store HookStore, queue HookQueue, logger log.Logger) *GohookdServer {
	// Business domain.
	var service Service
	{
		service = NewBasicService(store, queue)
		service = ServiceLoggingMiddleware(logger)(service)
	}

	// Endpoint domain.
	var tunnelEndpoint endpoint.Endpoint
	{
		tunnelLogger := log.NewContext(logger).With("method", "Tunnel")
		tunnelEndpoint = MakeTunnelEndpoint(service)
		tunnelEndpoint = EndpointLoggingMiddleware(tunnelLogger)(tunnelEndpoint)
	}

	var listEndpoint endpoint.Endpoint
	{
		listLogger := log.NewContext(logger).With("method", "List")
		listEndpoint = MakeListEndpoint(service)
		listEndpoint = EndpointLoggingMiddleware(listLogger)(listEndpoint)
	}

	var createEndpoint endpoint.Endpoint
	{
		createLogger := log.NewContext(logger).With("method", "Create")
		createEndpoint = MakeCreateEndpoint(service)
		createEndpoint = EndpointLoggingMiddleware(createLogger)(createEndpoint)
	}

	var deleteEndpoint endpoint.Endpoint
	{
		deleteLogger := log.NewContext(logger).With("method", "Delete")
		deleteEndpoint = MakeDeleteEndpoint(service)
		deleteEndpoint = EndpointLoggingMiddleware(deleteLogger)(deleteEndpoint)
	}

	endpoints := Endpoints{
		TunnelEndpoint: tunnelEndpoint,
		ListEndpoint:   listEndpoint,
		CreateEndpoint: createEndpoint,
		DeleteEndpoint: deleteEndpoint,
	}

	// Mechanical domain.
	logger = log.NewContext(logger).With("transport", "gRPC")
	return MakeGRPCServer(ctx, endpoints, logger)
}

type HookStore interface {
	Add(hook *Hook) error
	Remove(hookId HookID) (*Hook, error)
	Find(hookId HookID) (*Hook, error)
	FindAll() (HookList, error)
}
