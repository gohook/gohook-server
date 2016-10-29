package inmem

import (
	"errors"
	"sync"

	"github.com/gohook/gohook-server/gohookd"
	"github.com/gohook/gohook-server/user"
)

type InMemHooks struct {
	mtx       sync.RWMutex
	hooks     map[gohookd.HookID]*gohookd.Hook
	accountId user.AccountId
	scoped    bool
}

func NewInMemHooks() gohookd.HookStore {
	return &InMemHooks{
		hooks:  make(map[gohookd.HookID]*gohookd.Hook),
		scoped: false,
	}
}

func (i InMemHooks) Scope(accountId user.AccountId) gohookd.HookStore {
	i.accountId = accountId
	i.scoped = true
	return &i
}

func (i *InMemHooks) Add(m *gohookd.Hook) error {
	i.mtx.Lock()
	defer i.mtx.Unlock()
	if i.scoped {
		m.AccountId = i.accountId
	}
	i.hooks[m.Id] = m
	return nil
}

func (i *InMemHooks) Remove(id gohookd.HookID) (*gohookd.Hook, error) {
	i.mtx.Lock()
	defer i.mtx.Unlock()
	hook, ok := i.hooks[id]
	if ok {
		delete(i.hooks, id)
		return hook, nil
	}
	return nil, errors.New("Not Found")
}

func (i *InMemHooks) Find(id gohookd.HookID) (*gohookd.Hook, error) {
	i.mtx.RLock()
	defer i.mtx.RUnlock()
	if val, ok := i.hooks[id]; ok {
		if i.scoped && val.AccountId != i.accountId {
			return nil, errors.New("Not Found")
		}
		return val, nil
	}
	return nil, errors.New("Not Found")
}

func (i *InMemHooks) FindAll() (gohookd.HookList, error) {
	i.mtx.RLock()
	defer i.mtx.RUnlock()
	h := make(gohookd.HookList, 0, len(i.hooks))
	for _, val := range i.hooks {
		if i.scoped {
			if val.AccountId == i.accountId {
				h = append(h, val)
			}
		} else {
			h = append(h, val)
		}
	}
	return h, nil
}
