package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohae/deepcopy"
	"github.com/ory/fosite"
	"github.com/ory/fosite/token/jwt"
)

type Session struct {
	ID       string
	ClientID string

	Username  string
	Subject   string
	ExpiresAt map[fosite.TokenType]time.Time

	// Default
	Extra map[string]interface{}

	UserID string
	User   User
}

func NewSession(
	ctx context.Context,
	clientID string,
	userID string,
	username string,
	subject string,
	extra map[string]interface{},
) *Session {

	session := &Session{
		ID:       uuid.New().String(),
		UserID:   userID,
		ClientID: clientID,
		Username: username,
		Subject:  subject,
		Extra:    extra,
	}

	return session
}

// SetExpiresAt sets the expiration time of a token.
//
//	session.SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(time.Hour))
func (s *Session) SetExpiresAt(key fosite.TokenType, exp time.Time) {
	expiresAt := make(map[fosite.TokenType]time.Time)

	expiresAt[key] = exp

	s.ExpiresAt = expiresAt
}

// GetExpiresAt returns the expiration time of a token if set, or time.IsZero() if not.
//
//	session.GetExpiresAt(fosite.AccessToken)
func (s *Session) GetExpiresAt(key fosite.TokenType) time.Time {
	if s.ExpiresAt == nil {
		return time.Time{}
	}

	if _, ok := s.ExpiresAt[key]; !ok {
		return time.Time{}
	}

	return s.ExpiresAt[key]
}

// GetUsername returns the username, if set. This is optional and only used during token introspection.
func (s *Session) GetUsername() string {
	if s == nil {
		return ""
	}

	return s.Username
}

func (s *Session) GetExtraClaims() map[string]interface{} {
	if s == nil {
		return nil
	}

	return s.Extra
}

// GetSubject returns the subject, if set. This is optional and only used during token introspection.
func (s *Session) GetSubject() string {
	if s == nil {
		return ""
	}

	return s.Subject
}

// Clone clones the session.
func (s *Session) Clone() fosite.Session {
	if s == nil {
		return nil
	}

	return deepcopy.Copy(s).(fosite.Session)
}

// IDTokenClaims returns a pointer to claims which will be modified in-place by handlers.
// Session should store this pointer and return always the same pointer.
func (s *Session) IDTokenClaims() *jwt.IDTokenClaims {
	return &jwt.IDTokenClaims{}
}

// IDTokenHeaders returns a pointer to header values which will be modified in-place by handlers.
// Session should store this pointer and return always the same pointer.
func (s *Session) IDTokenHeaders() *jwt.Headers {
	return &jwt.Headers{}
}
