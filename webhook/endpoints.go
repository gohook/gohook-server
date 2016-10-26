package webhook

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type Endpoints struct {
	TriggerEndpoint endpoint.Endpoint
}

// Trigger Endpoint
func (e Endpoints) Trigger(ctx context.Context, trigger TriggerRequest) (*TriggerResponse, error) {
	response, err := e.TriggerEndpoint(ctx, trigger)
	return response.(*TriggerResponse), err
}

func MakeTriggerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (response interface{}, err error) {
		request := req.(TriggerRequest)
		return s.Trigger(ctx, request)
	}
}
