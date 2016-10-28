package gohookd

import (
	"fmt"

	"github.com/gohook/gohook-server/user"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

type Service interface {
	List(ctx context.Context) (HookList, error)
	Create(ctx context.Context, request HookRequest) (*Hook, error)
	Delete(ctx context.Context, id HookID) (*Hook, error)
}

func NewBasicService(store HookStore, authService user.AuthService) Service {
	return &basicService{
		hooks: store,
		auth:  authService,
	}
}

type basicService struct {
	hooks HookStore
	auth  user.AuthService
}

func (s basicService) List(ctx context.Context) (HookList, error) {
	account := ctx.Value("account").(*user.Account)
	hookList, err := s.hooks.Scope(account.Id).FindAll()
	if err != nil {
		return nil, err
	}
	return hookList, nil
}

func (s *basicService) Create(ctx context.Context, request HookRequest) (*Hook, error) {
	account := ctx.Value("account").(*user.Account)
	id := uuid.NewV4()
	newHook := &Hook{
		Id:     HookID(id.String()),
		Url:    fmt.Sprintf("http://localhost:8080/hook/%s", id.String()),
		Method: request.Method,
	}
	err := s.hooks.Scope(account.Id).Add(newHook)
	if err != nil {
		return nil, err
	}
	return newHook, nil
}

func (s *basicService) Delete(ctx context.Context, id HookID) (*Hook, error) {
	account := ctx.Value("account").(*user.Account)
	hook, err := s.hooks.Scope(account.Id).Remove(id)
	if err != nil {
		return nil, err
	}
	return hook, nil
}
