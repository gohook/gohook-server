package gohookd

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	TunnelEndpoint endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
	CreateEndpoint endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
}

// Tunnel Endpoint
type tunnelRequest struct{}
type tunnelResponse struct{}

func (e Endpoints) Tunnel(ctx context.Context) error {
	_, err := e.TunnelEndpoint(ctx, tunnelRequest{})
	return err
}

func MakeTunnelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (response interface{}, err error) {
		err = s.Tunnel(ctx)
		return tunnelResponse{}, err
	}
}

// List Endpoint
type listRequest struct{}

func (e Endpoints) List(ctx context.Context) (HookList, error) {
	response, err := e.ListEndpoint(ctx, listRequest{})
	return response.(HookList), err
}

func MakeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (response interface{}, err error) {
		list, err := s.List(ctx)
		if err != nil {
			return nil, err
		}
		return list, nil
	}
}

// Create Endpoint
func (e Endpoints) Create(ctx context.Context, request HookRequest) (*Hook, error) {
	response, err := e.CreateEndpoint(ctx, request)
	if err != nil {
		fmt.Println("Create Error Endpoint: ", err)
		return nil, err
	}
	return response.(*Hook), nil
}

func MakeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(HookRequest)
		hook, err := s.Create(ctx, req)
		if err != nil {
			return nil, err
		}
		return hook, nil
	}
}

// Delete Endpoint
func (e Endpoints) Delete(ctx context.Context, id HookID) (*Hook, error) {
	hook, err := e.DeleteEndpoint(ctx, id)
	if err != nil {
		return nil, err
	}
	return hook.(*Hook), nil
}

func MakeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id := request.(HookID)
		hook, err := s.Delete(ctx, id)
		if err != nil {
			return nil, err
		}
		return hook, nil
	}
}
