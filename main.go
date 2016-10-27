package main

import (
	"fmt"
	"net"
	"net/http"
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
	"github.com/gohook/gohook-server/webhook"
)

const (
	port     = "PORT"
	gRPCPort = "GRPC_PORT"
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

	gRPCPort := os.Getenv(gRPCPort)
	// default for port
	if gRPCPort == "" {
		gRPCPort = "9001"
	}

	// Setup Stores
	hookStore := inmem.NewInMemHooks()
	accountStore := inmem.NewInMemAccounts()

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
	var gohookdService gohookd.Service
	{
		gohookdService = gohookd.NewBasicService(hookStore)
		gohookdService = gohookd.ServiceLoggingMiddleware(logger)(gohookdService)
	}

	var webhookService webhook.Service
	{
		webhookService = webhook.NewBasicService(hookStore, queue)
		webhookService = webhook.ServiceLoggingMiddleware(logger)(webhookService)
	}

	// Endpoint domain.
	var listEndpoint endpoint.Endpoint
	{
		listLogger := log.NewContext(logger).With("method", "List")
		listEndpoint = gohookd.MakeListEndpoint(gohookdService)
		listEndpoint = gohookd.EndpointLoggingMiddleware(listLogger)(listEndpoint)
	}

	var createEndpoint endpoint.Endpoint
	{
		createLogger := log.NewContext(logger).With("method", "Create")
		createEndpoint = gohookd.MakeCreateEndpoint(gohookdService)
		createEndpoint = gohookd.EndpointLoggingMiddleware(createLogger)(createEndpoint)
	}

	var deleteEndpoint endpoint.Endpoint
	{
		deleteLogger := log.NewContext(logger).With("method", "Delete")
		deleteEndpoint = gohookd.MakeDeleteEndpoint(gohookdService)
		deleteEndpoint = gohookd.EndpointLoggingMiddleware(deleteLogger)(deleteEndpoint)
	}

	var triggerEndpoint endpoint.Endpoint
	{
		triggerLogger := log.NewContext(logger).With("method", "Trigger")
		triggerEndpoint = webhook.MakeTriggerEndpoint(webhookService)
		triggerEndpoint = webhook.EndpointLoggingMiddleware(triggerLogger)(triggerEndpoint)
	}

	// Interrupt handler
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP transport
	go func() {
		var webhooks http.Handler
		{
			endpoints := webhook.Endpoints{
				TriggerEndpoint: triggerEndpoint,
			}
			logger := log.NewContext(logger).With("transport", "HTTP")
			webhooks = webhook.MakeWebhookHTTPServer(ctx, endpoints, logger)
		}

		logger.Log("msg", "HTTP Server Started", "port", port)
		errc <- http.ListenAndServe(":"+port, webhooks)
	}()

	// gRPC transport
	go func() {
		lis, err := net.Listen("tcp", ":"+gRPCPort)
		if err != nil {
			errc <- err
			return
		}
		defer lis.Close()

		s := grpc.NewServer()

		// Mechanical domain.
		var gohook pb.GohookServer
		{
			endpoints := gohookd.Endpoints{
				ListEndpoint:   listEndpoint,
				CreateEndpoint: createEndpoint,
				DeleteEndpoint: deleteEndpoint,
			}
			logger := log.NewContext(logger).With("transport", "gRPC")
			g := gohookd.MakeGohookdServer(ctx, endpoints, logger)
			t, err := tunnel.MakeTunnelServer(accountStore, queue, logger)
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

		logger.Log("msg", "GRPC Server Started", "port", gRPCPort)
		errc <- s.Serve(lis)
	}()

	logger.Log("exit", <-errc)
}
