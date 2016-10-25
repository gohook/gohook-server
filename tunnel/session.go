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
	Add(sessionId Session) error
	Remove(sessionId SessionID) (*Session, error)
	FindByUserId(userId string) (*Session, error)
	FindBySessionId(sessionId SessionID) (*Session, error)
}
