package gohookd

import (
	"fmt"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

type Service interface {
	Tunnel(ctx context.Context) error
	List(ctx context.Context) (HookList, error)
	Create(ctx context.Context, request HookRequest) (*Hook, error)
	Delete(ctx context.Context, id HookID) (*Hook, error)
}

func NewBasicService(store HookStore) Service {
	return &basicService{
		hooks: store,
	}
}

type basicService struct {
	hooks HookStore
}

/*
  Tunnel
  ---------------------------

  All other service requests can go to any instance of the GRPC server, but
tunnel is special. Tunnel is a persistent connection that stays open between
the client and GRPC server. Because of that, when an HTTP hook gets hit, the
message can only propogate down to the client from the already open
connection (if there is one).


  In order to make this work, we will use an external messaging queue that all
GRPC servers are connected to. When a HTTP hook gets hit, it adds a message to
the queue with the clientID and payload data. Any GRPC server with a connected
client matching that clientID will format and send the message.

*/
func (s basicService) Tunnel(_ context.Context) error {
	// setup tunnel here.
	// must figure out the best way to send messages

	return nil
}

func (s basicService) List(_ context.Context) (HookList, error) {
	hookList, err := s.hooks.FindAll()
	if err != nil {
		return nil, err
	}
	return hookList, nil
}

func (s *basicService) Create(_ context.Context, request HookRequest) (*Hook, error) {
	id := uuid.NewV4()
	newHook := &Hook{
		Id:     HookID(id.String()),
		Url:    fmt.Sprintf("http://localhost:8080/hook/%s", id.String()),
		Method: request.Method,
	}
	err := s.hooks.Add(newHook)
	if err != nil {
		return nil, err
	}
	return newHook, nil
}

func (s *basicService) Delete(_ context.Context, id HookID) (*Hook, error) {
	hook, err := s.hooks.Remove(id)
	if err != nil {
		return nil, err
	}
	return hook, nil
}
