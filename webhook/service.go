package webhook

import (
	"github.com/gohook/gohook-server/gohookd"
	"github.com/gohook/gohook-server/tunnel"
	"golang.org/x/net/context"
)

type Service interface {
	Trigger(ctx context.Context, trigger TriggerRequest) (*TriggerResponse, error)
}

func NewBasicService(store gohookd.HookStore, sessions tunnel.SessionStore, queue tunnel.HookQueue) Service {
	return &basicService{
		hooks:    store,
		sessions: sessions,
		queue:    queue,
	}
}

type basicService struct {
	hooks    gohookd.HookStore
	sessions tunnel.SessionStore
	queue    tunnel.HookQueue
}

func (s basicService) Trigger(_ context.Context, trigger TriggerRequest) (*TriggerResponse, error) {
	// + Check to see if this hookid exists
	// If it does, check to see if there is a session associated with this user
	// If there is, broadcast a struct containing the sessionid and the hook payload

	hook, err := s.hooks.Find(trigger.HookId)
	if err != nil {
		return nil, err
	}

	// hook.UserId is hardcoded right now
	session, err := s.sessions.FindByUserId("myid")
	if err != nil {
		return nil, err
	}

	// Broadcast message with the sessionid and hook data
	err = s.queue.Broadcast(&tunnel.QueueMessage{
		SessionId: session.Id,
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
