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
	hook, err := s.hooks.Find(trigger.HookId)
	if err != nil {
		return nil, err
	}

	// Broadcast message with the userid and hook data
	err = s.queue.Broadcast(&tunnel.QueueMessage{
		AccountId: "myid",
		// AccountId: hook.AccountId,
		Hook: tunnel.HookCall{
			Id:     string(hook.Id),
			Method: trigger.Method,
			Body:   trigger.Body,
		},
	})
	if err != nil {
		return nil, err
	}
	return &TriggerResponse{200}, nil
}
