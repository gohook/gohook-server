package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/gohook/gohook-server/gohookd"
	"github.com/gohook/gohook-server/inmem"
	"github.com/gohook/pb"
)

const (
	port = "PORT"
)

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
		gohook := gohookd.NewGohookdGRPCServer(ctx, store, queue, logger)
		pb.RegisterGohookServer(s, gohook)

		logger.Log("msg", "GRPC Server Started", "port", port)
		errc <- s.Serve(lis)
	}()

	logger.Log("exit", <-errc)
}
