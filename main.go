package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/gohook/gohook-server/gohookd"
	"github.com/gohook/gohook-server/inmem"
	"github.com/gohook/gohook-server/pb"
	"github.com/gohook/gohook-server/tunnel"
)

const (
	port = "PORT"
)

type GohookGRPCServer struct {
	*gohookd.GohookdServer
	*tunnel.GohookTunnelServer
}

func main() {
	port := os.Getenv(port)
	// default for port
	if port == "" {
		port = "8080"
	}

	// Setup Store
	store := inmem.NewInMemHooks()

	// Setup Queue
	queue := inmem.NewInMemQueue()

	// Context
	ctx := context.Background()

	// Error chan
	errc := make(chan error)

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	// Business domain.
	var service gohookd.Service
	{
		service = gohookd.NewBasicService(store, queue)
		service = gohookd.ServiceLoggingMiddleware(logger)(service)
	}

	// Endpoint domain.
	var listEndpoint endpoint.Endpoint
	{
		listLogger := log.NewContext(logger).With("method", "List")
		listEndpoint = gohookd.MakeListEndpoint(service)
		listEndpoint = gohookd.EndpointLoggingMiddleware(listLogger)(listEndpoint)
	}

	var createEndpoint endpoint.Endpoint
	{
		createLogger := log.NewContext(logger).With("method", "Create")
		createEndpoint = gohookd.MakeCreateEndpoint(service)
		createEndpoint = gohookd.EndpointLoggingMiddleware(createLogger)(createEndpoint)
	}

	var deleteEndpoint endpoint.Endpoint
	{
		deleteLogger := log.NewContext(logger).With("method", "Delete")
		deleteEndpoint = gohookd.MakeDeleteEndpoint(service)
		deleteEndpoint = gohookd.EndpointLoggingMiddleware(deleteLogger)(deleteEndpoint)
	}

	endpoints := gohookd.Endpoints{
		ListEndpoint:   listEndpoint,
		CreateEndpoint: createEndpoint,
		DeleteEndpoint: deleteEndpoint,
	}

	// Interrupt handler
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// gRPC transport
	go func() {
		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			errc <- err
			return
		}
		defer lis.Close()

		s := grpc.NewServer()

		// Mechanical domain.
		var gohook pb.GohookServer
		{
			logger := log.NewContext(logger).With("transport", "gRPC")
			g := gohookd.MakeGohookdServer(ctx, endpoints, logger)
			t, err := tunnel.MakeTunnelServer(queue, logger)
			if err != nil {
				errc <- err
				return
			}

			gohook = &GohookGRPCServer{
				GohookTunnelServer: t,
				GohookdServer:      g,
			}
		}

		pb.RegisterGohookServer(s, gohook)

		logger.Log("msg", "GRPC Server Started", "port", port)
		errc <- s.Serve(lis)
	}()

	logger.Log("exit", <-errc)
}
