package webhook

import (
	"github.com/gohook/gohook-server/gohookd"
	"github.com/gohook/gohook-server/tunnel"
	"golang.org/x/net/context"
)

type Service interface {
	Trigger(ctx context.Context, trigger TriggerRequest) (*TriggerResponse, error)
}

func NewBasicService(store gohookd.HookStore, queue tunnel.HookQueue) Service {
	return &basicService{
		hooks: store,
		queue: queue,
	}
}

type basicService struct {
	hooks gohookd.HookStore
	queue tunnel.HookQueue
}

func (s basicService) Trigger(_ context.Context, trigger TriggerRequest) (*TriggerResponse, error) {
	// + Check to see if this hookid exists
	// If it does, check to see if there is a session associated with this user
	// If there is, broadcast a struct containing the sessionid and the hook payload

	hook, err := s.hooks.Find(trigger.HookId)
	if err != nil {
		return nil, err
	}

	err = s.queue.Broadcast(&tunnel.QueueMessage{
		SessionId: string(hook.Id),
	})
	if err != nil {
		return nil, err
	}
	return &TriggerResponse{200}, nil
}
