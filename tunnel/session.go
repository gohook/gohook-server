package tunnel

import (
	"time"
)

type SessionID string

type Session struct {
	Id     SessionID
	UserId string
	Start  time.Time
}

type SessionStore interface {
	Add(session *Session) error
	Remove(sessionId SessionID) error
	FindByUserId(userId string) (*Session, error)
	FindBySessionId(sessionId SessionID) (*Session, error)
}
