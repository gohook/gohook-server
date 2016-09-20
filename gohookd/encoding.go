package gohookd

import (
	"errors"
	"golang.org/x/net/context"

	"github.com/gohook/pb"
)

// Tunnel transforms
func EncodeGRPCTunnelRequest(_ context.Context, request interface{}) (interface{}, error) {
	_ = request.(tunnelRequest)
	return &pb.TunnelRequest{}, nil
}

func DecodeGRPCTunnelRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	_ = grpcReq.(*pb.TunnelRequest)
	return tunnelRequest{}, nil
}

func DecodeGRPCTunnelResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	_ = grpcReply.(*pb.TunnelResponse)
	return tunnelResponse{}, nil
}

func EncodeGRPCTunnelResponse(_ context.Context, response interface{}) (interface{}, error) {
	_ = response.(tunnelResponse)
	return &pb.TunnelResponse{}, nil
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
