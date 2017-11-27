package auth

import (
	"context"
	"net/http/httptest"
	"testing"
)

var ctx context.Context

func TestNewSession(t *testing.T) {
	s := NewSession("1")

	if len(s.Token) != 40 {
		t.Errorf("token != 40")
	}
}

func TestRequestHeader(t *testing.T) {
	mock := NewSession("1")
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer "+mock.Token)

	session := FromRequest(req)

	if err := session.Err(); err != nil {
		t.Error(err)
	}
	if session.Token != mock.Token {
		t.Errorf("%s != %s", session.Token, mock.Token)
	}

	session.ID = "1"
	ctx = session.Context(context.Background())
}
func TestEmptyAuthorizationHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	session := FromRequest(req)
	if err := session.Err(); err != nil {
		t.Error(err)
	}
	if !session.IsZero() {
		t.Errorf("not empty: %+v", session)
	}
}

func TestAuthorizationHeaderErrors(t *testing.T) {
	tables := []struct {
		header string
		err    error
	}{
		{"Token TEST_TOKEN", ErrMissingBearerPrefix},
		{"Bearer", ErrMissingSpaceCharacter},
		{"Bearer ", ErrMissingAuthorization},
		{"Bearer TEST TOKEN", ErrTooManySpaceCharacters},
		{"BearerTEST_TOKEN", ErrMissingSpaceCharacter},
	}
	for _, table := range tables {
		req := httptest.NewRequest("GET", "http://example.com", nil)
		req.Header.Set("Authorization", table.header)
		session := FromRequest(req)
		if err := session.Err(); err != table.err {
			t.Errorf("%s = %v", table.header, err)
		}
	}
}

func TestFromContext(t *testing.T) {
	session := FromContext(ctx)

	if err := session.Err(); err != nil {
		t.Error(err)
	}
	if session.IsZero() {
		t.Errorf("%+v == empty", session)
	}
}
