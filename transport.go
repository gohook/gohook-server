package main

import (
	"golang.org/x/net/context"

	pb "github.com/vems/pb/add"

	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

func MakeGRPCServer(ctx context.Context, endpoints Endpoints, logger log.Logger) pb.AddServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		sum: grpctransport.NewServer(
			ctx,
			endpoints.SumEndpoint,
			DecodeGRPCSumRequest,
			EncodeGRPCSumResponse,
			options...,
		),
	}
}

type grpcServer struct {
	sum grpctransport.Handler
}

func (s *grpcServer) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumReply, error) {
	_, rep, err := s.sum.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SumReply), nil
}

func DecodeGRPCSumRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SumRequest)
	return sumRequest{A: int(req.A), B: int(req.B)}, nil
}

func DecodeGRPCSumResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SumReply)
	return sumResponse{V: int(reply.V), Err: str2err(reply.Err)}, nil
}

func EncodeGRPCSumResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(sumResponse)
	return &pb.SumReply{V: int64(resp.V), Err: err2str(resp.Err)}, nil
}

func EncodeGRPCSumRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(sumRequest)
	return &pb.SumRequest{A: int64(req.A), B: int64(req.B)}, nil
}
