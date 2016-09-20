package gohookd

import (
	"errors"
	"golang.org/x/net/context"

	"github.com/gohook/pb"

	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type GohookdServer struct {
	list   grpctransport.Handler
	create grpctransport.Handler
	delete grpctransport.Handler
}

func MakeGRPCServer(ctx context.Context, endpoints Endpoints, logger log.Logger) *GohookdServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &GohookdServer{
		list: grpctransport.NewServer(
			ctx,
			endpoints.ListEndpoint,
			DecodeGRPCListRequest,
			EncodeGRPCListResponse,
			options...,
		),
		create: grpctransport.NewServer(
			ctx,
			endpoints.CreateEndpoint,
			DecodeGRPCCreateRequest,
			EncodeGRPCCreateResponse,
			options...,
		),
		delete: grpctransport.NewServer(
			ctx,
			endpoints.DeleteEndpoint,
			DecodeGRPCDeleteRequest,
			EncodeGRPCDeleteResponse,
			options...,
		),
	}
}

// Tunnel transport handler
func (s *GohookdServer) Tunnel(req *pb.TunnelRequest, _ pb.Gohook_TunnelServer) error {
	return errors.New("Not Implimented. Use other tunnel method.")
}

// List transport handler
func (s *GohookdServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	_, rep, err := s.list.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListResponse), nil
}

// Create transport handler
func (s *GohookdServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	_, rep, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateResponse), nil
}

// Delete transport handler
func (s *GohookdServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	_, rep, err := s.delete.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.DeleteResponse), nil
}
