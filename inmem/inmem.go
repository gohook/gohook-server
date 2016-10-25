package inmem

import (
	"errors"
	"sync"

	"github.com/gohook/gohook-server/gohookd"
	"github.com/gohook/gohook-server/tunnel"
)

/*
InMemQueue
----------

InMemQueue impliments the HookQueue interface in the most
basic way that allows the message sending to happen from
within the same process. This makes testing much easier and
allows the single gohookd process to run without external
dependencies.

THIS DOES NOT SCALE. Only use for testing and single client
setups.

Since gohookd processes can't communicate when they receive
a hook message, there is no guarantee that the hook message
will go to the process that the client is connected to.
*/

type InMemQueue struct {
	receivec tunnel.ReceiveC
}

func NewInMemQueue() tunnel.HookQueue {
	return InMemQueue{
		receivec: make(tunnel.ReceiveC),
	}
}

func (i InMemQueue) Broadcast(m *tunnel.QueueMessage) error {
	i.receivec <- m
	return nil
}

func (i InMemQueue) Listen() (tunnel.ReceiveC, error) {
	return i.receivec, nil
}

type InMemHooks struct {
	mtx   sync.RWMutex
	hooks map[gohookd.HookID]*gohookd.Hook
}

func NewInMemHooks() gohookd.HookStore {
	return &InMemHooks{
		hooks: make(map[gohookd.HookID]*gohookd.Hook),
	}
}

func (i *InMemHooks) Add(m *gohookd.Hook) error {
	i.mtx.Lock()
	defer i.mtx.Unlock()
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
		return val, nil
	}
	return nil, errors.New("Not Found")
}

func (i *InMemHooks) FindAll() (gohookd.HookList, error) {
	i.mtx.RLock()
	defer i.mtx.RUnlock()
	h := make(gohookd.HookList, 0, len(i.hooks))
	for _, val := range i.hooks {
		h = append(h, val)
	}
	return h, nil
}
