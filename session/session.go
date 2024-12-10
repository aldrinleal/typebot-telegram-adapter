package session

import "errors"

var ENotFound = errors.New("sessionId not found.")

type SessionManager interface {
	// LookupSession returns a session by its int64 id
	LookupSession(userId int64) (string, error)
	RefreshSession(userId int64) error
	RegisterSession(userId int64, sessionId string) error
	UnregisterSession(userId int64) error
}

func (m *InMemorySessionManager) UnregisterSession(userId int64) error {
	delete(m.sessions, userId)

	return nil
}
