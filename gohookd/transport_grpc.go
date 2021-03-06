package gohookd

import (
	"errors"

	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/gohook/gohook-server/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

type GohookdServer struct {
	list   grpctransport.Handler
	create grpctransport.Handler
	delete grpctransport.Handler
}

func extractAuthToken(ctx context.Context, md *metadata.MD) context.Context {
	if token, ok := (*md)["token"]; ok && len(token) > 0 {
		return context.WithValue(ctx, "token", token[0])
	}
	return ctx
}

func MakeGohookdServer(ctx context.Context, endpoints Endpoints, logger log.Logger) *GohookdServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
		grpctransport.ServerBefore(extractAuthToken),
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
	id := request.(HookID)
	return &pb.DeleteRequest{string(id)}, nil
}

func DecodeGRPCDeleteRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeleteRequest)
	return HookID(req.Id), nil
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
