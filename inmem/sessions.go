package inmem

import (
	"errors"
	"github.com/gohook/gohook-server/tunnel"
	"sync"
)

type InMemSessions struct {
	mtx      sync.RWMutex
	sessions map[tunnel.SessionID]*tunnel.Session
}

func NewInMemSessions() tunnel.SessionStore {
	return &InMemSessions{
		sessions: make(map[tunnel.SessionID]*tunnel.Session),
	}
}

func (i *InMemSessions) Add(s *tunnel.Session) error {
	i.mtx.Lock()
	defer i.mtx.Unlock()
	i.sessions[s.Id] = s
	return nil
}

func (i *InMemSessions) Remove(id tunnel.SessionID) error {
	i.mtx.Lock()
	defer i.mtx.Unlock()
	_, ok := i.sessions[id]
	if ok {
		delete(i.sessions, id)
		return nil
	}
	return errors.New("Not Found")
}

func (i *InMemSessions) FindBySessionId(id tunnel.SessionID) (*tunnel.Session, error) {
	i.mtx.RLock()
	defer i.mtx.RUnlock()
	if val, ok := i.sessions[id]; ok {
		return val, nil
	}
	return nil, errors.New("Not Found")
}

func (i *InMemSessions) FindByUserId(id string) (*tunnel.Session, error) {
	i.mtx.RLock()
	defer i.mtx.RUnlock()
	for _, s := range i.sessions {
		if s.UserId == id {
			return s, nil
		}
	}
	return nil, errors.New("Not Found")
}
