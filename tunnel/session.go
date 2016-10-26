package tunnel

import (
	"errors"
	"sync"
	"time"

	"github.com/gohook/gohook-server/pb"
)

type SessionID string

type Session struct {
	Id     SessionID
	UserId string
	Start  time.Time
	Stream pb.Gohook_TunnelServer
}

type SessionList []*Session

type SessionStore struct {
	mtx sync.RWMutex
	// Sessions map with userID as the key and an array of
	// connected sessions as the value.
	sessions map[string]SessionList
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]SessionList),
	}
}

func (s *SessionStore) Add(session *Session) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.sessions[session.UserId] = append(s.sessions[session.UserId], session)
	return nil
}

func (s *SessionStore) Remove(userId string, id SessionID) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	sessions, ok := s.sessions[userId]
	if ok {
		for i, session := range sessions {
			if session.Id == id {
				s.sessions[userId] = append(s.sessions[userId][:i], s.sessions[userId][i+1:]...)
				return nil
			}
		}
	}
	return errors.New("Not Found")
}

func (s *SessionStore) FindBySessionId(id SessionID) (*Session, error) {
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

func (s *SessionStore) FindByUserId(userId string) (SessionList, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	sessions, ok := s.sessions[userId]
	if ok {
		return sessions, nil
	}
	return nil, errors.New("Not Found")
}
