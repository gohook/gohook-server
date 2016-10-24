package webhook

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type Endpoints struct {
	TriggerEndpoint endpoint.Endpoint
}

// Trigger Endpoint
type triggerRequest struct{}

func (e Endpoints) Trigger(ctx context.Context) (*WebhookStatus, error) {
	response, err := e.TriggerEndpoint(ctx, triggerRequest{})
	return response.(*WebhookStatus), err
}

func MakeTriggerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (response interface{}, err error) {
		status, err := s.Trigger(ctx)
		if err != nil {
			return nil, err
		}
		return status, nil
	}
}
