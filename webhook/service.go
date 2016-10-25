package webhook

import (
	"github.com/gohook/gohook-server/gohookd"
	"golang.org/x/net/context"
)

type Service interface {
	Trigger(ctx context.Context, hookId string) (*WebhookStatus, error)
}

func NewBasicService(store gohookd.HookStore, queue gohookd.HookQueue) Service {
	return &basicService{
		hooks: store,
		queue: queue,
	}
}

type basicService struct {
	hooks gohookd.HookStore
	queue gohookd.HookQueue
}

func (s basicService) Trigger(_ context.Context, hookId string) (*WebhookStatus, error) {
	// Check to see if this hookid exists
	// If it does, check to see if there is a session associated with this user
	// If there is, broadcast a struct containing the sessionid and the hook payload
	err := s.queue.Broadcast(hookId)
	if err != nil {
		return nil, err
	}
	return &WebhookStatus{200}, nil
}
