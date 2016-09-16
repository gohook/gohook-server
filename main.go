package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/gohook/pb"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

func main() {
	var (
		grpcAddr = flag.String("addr", ":8080", "gRPC (HTTP) listen address")
	)

	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}
	logger.Log("msg", "Gohookd Started")
	defer logger.Log("msg", "Gohookd Quit")

	// Business domain.
	var service Service
	{
		service = NewBasicService()
		service = ServiceLoggingMiddleware(logger)(service)
	}

	// Endpoint domain.
	var sumEndpoint endpoint.Endpoint
	{
		sumLogger := log.NewContext(logger).With("method", "Sum")
		sumEndpoint = MakeSumEndpoint(service)
		sumEndpoint = EndpointLoggingMiddleware(sumLogger)(sumEndpoint)
	}

	endpoints := Endpoints{
		SumEndpoint: sumEndpoint,
	}

	// Mechanical domain.
	errc := make(chan error)
	ctx := context.Background()

	// Interrupt handler.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// gRPC transport.
	go func() {
		logger := log.NewContext(logger).With("transport", "gRPC")

		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}

		srv := MakeGRPCServer(ctx, endpoints, logger)
		s := grpc.NewServer()
		pb.RegisterAddServer(s, srv)

		logger.Log("addr", *grpcAddr)
		errc <- s.Serve(ln)
	}()

	logger.Log("exit", <-errc)
}
