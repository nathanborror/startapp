package auth

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"strings"
)

const secret = "my claims are my password"

var (
	ErrTooManySpaceCharacters = errors.New("Too many space characters in Authorization header")
	ErrMissingSpaceCharacter  = errors.New("Missing space character in Authorization header")
	ErrMissingBearerPrefix    = errors.New("Missing 'Bearer' prefix in Authorization header.")
	ErrMissingAuthorization   = errors.New("Missing Authorization token")
)

type Session struct {
	ID    string
	Token string
	err   error
}

func NewSession(id string) *Session {
	return &Session{
		ID:    id,
		Token: fakeToken(40),
	}
}

func FromRequest(r *http.Request) *Session {
	var s Session
	header := r.Header.Get("Authorization")
	if header == "" {
		return &s
	}
	components := strings.Split(header, " ")
	if len(components) > 2 {
		s.err = ErrTooManySpaceCharacters
	} else if len(components) < 2 {
		s.err = ErrMissingSpaceCharacter
	} else if components[0] != "Bearer" {
		s.err = ErrMissingBearerPrefix
	} else if len(components[1]) == 0 {
		s.err = ErrMissingAuthorization
	}
	if s.err != nil {
		return &s
	}
	s.Token = components[1]
	return &s
}

func FromContext(ctx context.Context) *Session {
	return ctx.Value(sessionKey).(*Session)
}

func (s *Session) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, sessionKey, s)
}

func (s *Session) IsZero() bool {
	return s.ID == "" || s.Token == ""
}

func (s *Session) Err() error {
	return s.err
}

type key int

const sessionKey key = 0

const tokenValidChars = "1234567890ABCDEF"

func fakeToken(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = tokenValidChars[rand.Intn(len(tokenValidChars))]
	}
	return string(b)
}
