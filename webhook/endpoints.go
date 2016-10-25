package webhook

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type Endpoints struct {
	TriggerEndpoint endpoint.Endpoint
}

// Trigger Endpoint
type triggerRequest struct {
	hookId string `json:"hookId"`
}

func (e Endpoints) Trigger(ctx context.Context, hookId string) (*WebhookStatus, error) {
	response, err := e.TriggerEndpoint(ctx, triggerRequest{hookId})
	return response.(*WebhookStatus), err
}

func MakeTriggerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (response interface{}, err error) {
		request := req.(triggerRequest)
		status, err := s.Trigger(ctx, request.hookId)
		if err != nil {
			return nil, err
		}
		return status, nil
	}
}
