package client

import (
	"time"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"

	"github.com/gohook/gohook-server/gohookd"
	"github.com/gohook/gohook-server/pb"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

func New(conn *grpc.ClientConn, logger log.Logger) gohookd.Service {

	var listEndpoint endpoint.Endpoint
	{
		listEndpoint = grpctransport.NewClient(
			conn,
			"Gohook",
			"List",
			gohookd.EncodeGRPCListRequest,
			gohookd.DecodeGRPCListResponse,
			pb.ListResponse{},
		).Endpoint()
		listEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "List",
			Timeout: 30 * time.Second,
		}))(listEndpoint)
	}

	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = grpctransport.NewClient(
			conn,
			"Gohook",
			"Create",
			gohookd.EncodeGRPCCreateRequest,
			gohookd.DecodeGRPCCreateResponse,
			pb.CreateResponse{},
		).Endpoint()
		createEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Create",
			Timeout: 30 * time.Second,
		}))(createEndpoint)
	}

	var deleteEndpoint endpoint.Endpoint
	{
		deleteEndpoint = grpctransport.NewClient(
			conn,
			"Gohook",
			"Delete",
			gohookd.EncodeGRPCDeleteRequest,
			gohookd.DecodeGRPCDeleteResponse,
			pb.DeleteResponse{},
		).Endpoint()
		deleteEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Delete",
			Timeout: 30 * time.Second,
		}))(deleteEndpoint)
	}

	return gohookd.Endpoints{
		ListEndpoint:   listEndpoint,
		CreateEndpoint: createEndpoint,
		DeleteEndpoint: deleteEndpoint,
	}
}
