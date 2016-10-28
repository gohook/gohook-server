package tunnel

import (
	"errors"
	"sync"
	"time"

	"github.com/gohook/gohook-server/pb"
	"github.com/gohook/gohook-server/user"
)

type SessionId string

type Session struct {
	Id        SessionId
	AccountId user.AccountId
	Start     time.Time
	Stream    pb.Gohook_TunnelServer
}

type SessionList []*Session

type SessionStore struct {
	mtx sync.RWMutex
	// Sessions map with userID as the key and an array of
	// connected sessions as the value.
	sessions map[user.AccountId]SessionList
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[user.AccountId]SessionList),
	}
}

func (s *SessionStore) Add(session *Session) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.sessions[session.AccountId] = append(s.sessions[session.AccountId], session)
	return nil
}

func (s *SessionStore) Remove(accountId user.AccountId, id SessionId) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	sessions, ok := s.sessions[accountId]
	if ok {
		for i, session := range sessions {
			if session.Id == id {
				s.sessions[accountId] = append(s.sessions[accountId][:i], s.sessions[accountId][i+1:]...)
				return nil
			}
		}
	}
	return errors.New("Not Found")
}

func (s *SessionStore) FindBySessionId(id SessionId) (*Session, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	for _, sessions := range s.sessions {
		for _, session := range sessions {
			if session.Id == id {
				return session, nil
			}
		}
	}
	return nil, errors.New("Not Found")
}

func (s *SessionStore) FindByAccountId(accountId user.AccountId) (SessionList, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	sessions, ok := s.sessions[accountId]
	if ok {
		return sessions, nil
	}
	return nil, errors.New("Not Found")
}
