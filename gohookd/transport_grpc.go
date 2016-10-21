package gohookd

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"time"

	"github.com/gohook/gohook-server/pb"

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
func (s *GohookdServer) Tunnel(req *pb.TunnelRequest, stream pb.Gohook_TunnelServer) error {
	// pass the stream and request data down to the service?
	// that way we can handle the logic there... but... then we are mixing transport and
	// business logic...
	// is there some way we can expose an interface to the service?
	for {
		err := stream.Send(&pb.TunnelResponse{
			Event: &pb.TunnelResponse_Hook{
				Hook: &pb.HookCall{
					Id: "Hello, World",
				},
			},
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(1 * time.Second)
	}
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

// List transforms
func EncodeGRPCListRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return &pb.ListRequest{}, nil
}

func DecodeGRPCListRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return listRequest{}, nil
}

func EncodeGRPCListResponse(_ context.Context, response interface{}) (interface{}, error) {
	hooks := response.(HookList)
	pbHooks := []*pb.Hook{}

	for _, h := range hooks {
		method, ok := pb.Method_value[h.Method]
		if !ok {
			return nil, errors.New("Invalid Method Name")
		}
		pbHooks = append(pbHooks, &pb.Hook{
			Id:     string(h.Id),
			Url:    h.Url,
			Method: pb.Method(method),
		})
	}
	return &pb.ListResponse{pbHooks}, nil
}

func DecodeGRPCListResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	resp := grpcReply.(*pb.ListResponse)
	modelHooks := HookList{}
	for _, h := range resp.Hooks {
		methodName, ok := pb.Method_name[int32(h.Method)]
		if !ok {
			return nil, errors.New("Invalid Method ID")
		}
		modelHooks = append(modelHooks, &Hook{
			Id:     HookID(h.Id),
			Url:    h.Url,
			Method: methodName,
		})
	}
	return modelHooks, nil
}

// Create transforms
func EncodeGRPCCreateRequest(_ context.Context, request interface{}) (interface{}, error) {
	hook := request.(HookRequest)
	methodID, ok := pb.Method_value[hook.Method]
	if !ok {
		return nil, errors.New("1 Invalid Method Name")
	}
	createReq := &pb.HookRequest{
		Method: pb.Method(methodID),
	}
	return &pb.CreateRequest{createReq}, nil
}

func DecodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	createReq := grpcReq.(*pb.CreateRequest)
	hookReq := createReq.Hook
	method, ok := pb.Method_name[int32(hookReq.Method)]
	if !ok {
		return nil, errors.New("2 Invalid Method Name")
	}
	hook := HookRequest{
		Method: method,
	}
	return hook, nil
}

func EncodeGRPCCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	createRes := response.(*Hook)
	method, ok := pb.Method_value[createRes.Method]
	if !ok {
		return nil, errors.New("3 Invalid Method Name")
	}
	hook := &pb.Hook{
		Id:     string(createRes.Id),
		Url:    createRes.Url,
		Method: pb.Method(method),
	}
	return &pb.CreateResponse{hook}, nil
}

func DecodeGRPCCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	createRes := response.(*pb.CreateResponse)
	hookRes := createRes.Hook
	method, ok := pb.Method_name[int32(hookRes.Method)]
	if !ok {
		return nil, errors.New("4 Invalid Method Name")
	}
	hook := &Hook{
		Id:     HookID(hookRes.Id),
		Url:    hookRes.Url,
		Method: method,
	}
	return hook, nil
}

// Delete transforms
func EncodeGRPCDeleteRequest(_ context.Context, request interface{}) (interface{}, error) {
	id := request.(string)
	return &pb.DeleteRequest{id}, nil
}

func DecodeGRPCDeleteRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeleteRequest)
	return req.Id, nil
}

func EncodeGRPCDeleteResponse(_ context.Context, response interface{}) (interface{}, error) {
	deleteRes := response.(*Hook)
	method, ok := pb.Method_value[deleteRes.Method]
	if !ok {
		return nil, errors.New("Invalid Method Name")
	}
	hook := &pb.Hook{
		Id:     string(deleteRes.Id),
		Url:    deleteRes.Url,
		Method: pb.Method(method),
	}
	return &pb.DeleteResponse{hook}, nil
}

func DecodeGRPCDeleteResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	deleteRes := grpcReply.(*pb.DeleteResponse)
	hookRes := deleteRes.Hook
	method, ok := pb.Method_name[int32(hookRes.Method)]
	if !ok {
		return nil, errors.New("Invalid Method Name")
	}
	hook := &Hook{
		Id:     HookID(hookRes.Id),
		Url:    hookRes.Url,
		Method: method,
	}
	return hook, nil
}
