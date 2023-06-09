package core

import (
	"context"

	"github.com/gotd/td/session"
)

// Session абстракция над сессией телеграм аккаунта
type Session struct {
	// mux  sync.RWMutex
	SessionID int
	Data      []byte
}

// Session абстракция над сессией телеграм аккаунта
type SessionData struct {
	Version int          `json:"Version"`
	Data    session.Data `json:"Data"`
}

// LoadSession loads session from memory.
func (s *Session) LoadSession(context.Context) ([]byte, error) {
	if s == nil {
		return nil, session.ErrNotFound
	}

	// s.mux.RLock()
	// defer s.mux.RUnlock()

	if len(s.Data) == 0 {
		return nil, session.ErrNotFound
	}

	cpy := append([]byte(nil), s.Data...)

	return cpy, nil
}

// StoreSession stores session to memory.
func (s *Session) StoreSession(ctx context.Context, data []byte) error {
	// s.mux.Lock()
	s.Data = data
	// s.mux.Unlock()
	return nil
}

type ExtractSessionResult struct {
	SaveZipCounter int
	UnZipCounter   int
}
