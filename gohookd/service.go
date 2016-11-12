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

type ServiceOpts struct {
	Origin   string
	Protocol string
}

var DefaultServiceOpts = ServiceOpts{
	Protocol: "http",
	Origin:   "localhost",
}

func NewServiceOpts() *ServiceOpts {
	return &DefaultServiceOpts
}

func WithOrigin(origin string) *ServiceOpts {
	opts := &DefaultServiceOpts
	opts.Origin = origin
	return opts
}

func NewBasicService(store HookStore, authService user.AuthService, opts *ServiceOpts) Service {
	return &basicService{
		hooks: store,
		auth:  authService,
		opts:  opts,
	}
}

type basicService struct {
	hooks HookStore
	auth  user.AuthService
	opts  *ServiceOpts
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
		Url:    fmt.Sprintf("%s://%s/%s/%s", s.opts.Protocol, s.opts.Origin, account.Id, id.String()),
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
