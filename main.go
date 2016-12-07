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
	"gopkg.in/mgo.v2"

	"github.com/gohook/gohook-server/auth"
	"github.com/gohook/gohook-server/gohookd"
	"github.com/gohook/gohook-server/mongo"
	"github.com/gohook/gohook-server/pb"
	"github.com/gohook/gohook-server/redis"
	"github.com/gohook/gohook-server/tunnel"
	"github.com/gohook/gohook-server/webhook"
)

const (
	port             = "PORT"
	gRPCPort         = "GRPC_PORT"
	httpServerOrigin = "HTTP_ORIGIN"
	mongoAddr        = "MONGO_URL"
	redisAddr        = "REDIS_ADDR"
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

	httpServerOrigin := os.Getenv(httpServerOrigin)
	// default for port
	if httpServerOrigin == "" {
		httpServerOrigin = "localhost:8080"
	}

	mongoAddr := os.Getenv(mongoAddr)
	// default for mongo
	if mongoAddr == "" {
		mongoAddr = "127.0.0.1"
	}

	redisAddr := os.Getenv(redisAddr)
	// default for mongo
	if redisAddr == "" {
		redisAddr = ":6379"
	}

	// Setup Stores
	// Setup Mongo DB connection
	session, err := mgo.Dial(mongoAddr)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	hookStore, err := mongo.NewMongoHookStore("gohook", session)
	if err != nil {
		panic(err)
	}

	accountStore := mongo.NewMongoAccountStore("gohook", session)

	// Setup AuthService
	authService := auth.NewAuthService(accountStore)

	// Setup Queue
	queue, err := redis.NewRedisQueue(redisAddr)
	if err != nil {
		panic(err)
	}

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
		gohookdService = gohookd.NewBasicService(hookStore, authService, gohookd.WithOrigin(httpServerOrigin))
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
		listEndpoint = gohookd.EndpointAuthMiddleware(listLogger, authService)(listEndpoint)
		listEndpoint = gohookd.EndpointLoggingMiddleware(listLogger)(listEndpoint)
	}

	var createEndpoint endpoint.Endpoint
	{
		createLogger := log.NewContext(logger).With("method", "Create")
		createEndpoint = gohookd.MakeCreateEndpoint(gohookdService)
		createEndpoint = gohookd.EndpointAuthMiddleware(createLogger, authService)(createEndpoint)
		createEndpoint = gohookd.EndpointLoggingMiddleware(createLogger)(createEndpoint)
	}

	var deleteEndpoint endpoint.Endpoint
	{
		deleteLogger := log.NewContext(logger).With("method", "Delete")
		deleteEndpoint = gohookd.MakeDeleteEndpoint(gohookdService)
		deleteEndpoint = gohookd.EndpointAuthMiddleware(deleteLogger, authService)(deleteEndpoint)
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
			webhooks = webhook.MakeWebhookHTTPServer(ctx, endpoints, logger, httpServerOrigin)
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
			t, err := tunnel.MakeTunnelServer(authService, queue, logger)
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
