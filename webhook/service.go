package webhook

import (
	"github.com/gohook/gohook-server/gohookd"
	"golang.org/x/net/context"
)

type Service interface {
	Trigger(ctx context.Context) (*WebhookStatus, error)
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

func (s basicService) Trigger(_ context.Context) (*WebhookStatus, error) {
	err := s.queue.Broadcast("From Hook")
	if err != nil {
		return nil, err
	}
	return &WebhookStatus{200}, nil
}
