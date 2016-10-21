package gohookd

import (
	"fmt"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

type Service interface {
	List(ctx context.Context) (HookList, error)
	Create(ctx context.Context, request HookRequest) (*Hook, error)
	Delete(ctx context.Context, id HookID) (*Hook, error)
}

func NewBasicService(store HookStore, queue HookQueue) Service {
	return &basicService{
		hooks: store,
		queue: queue,
	}
}

type basicService struct {
	hooks HookStore
	queue HookQueue
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
