package session

// Implements an InMemory Session Manager
type InMemorySessionManager struct {
	sessions map[int64]string
}

func NewInMemorySessionManager() *InMemorySessionManager {
	return &InMemorySessionManager{
		sessions: make(map[int64]string),
	}
}

func (m *InMemorySessionManager) RefreshSession(userId int64) error {
	return nil
}

func (m *InMemorySessionManager) LookupSession(userId int64) (string, error) {
	sessionId, ok := m.sessions[userId]
	if !ok {
		return "", ENotFound
	}

	return sessionId, nil
}

func (m *InMemorySessionManager) RegisterSession(userId int64, sessionId string) error {
	m.sessions[userId] = sessionId

	return nil
}
