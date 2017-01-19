package oauth2utils

import (
	"sync"

	"golang.org/x/oauth2"
)

// ReactiveTokenSource is a TokenSource that holds a single token in memory
// and validates its expiry before each call to retrieve it with
// Token. If it's expired, it will be auto-refreshed using the
// new TokenSource and onUpdate callback will be called with new token.
type ReactiveTokenSource struct {
	new oauth2.TokenSource // called when t is expired.

	mu sync.Mutex // guards t
	t  *oauth2.Token

	onUpdate func(t *oauth2.Token)
}

// NewReactiveTokenSource returns a ReactiveTokenSource which repeatedly
// returns the same token as long as it's valid, starting with t.
// When its cached token is invalid, a new token is obtained from src
// and onUpdate callback is called.
func NewReactiveTokenSource(t *oauth2.Token, new oauth2.TokenSource, onUpdate func(t *oauth2.Token)) *ReactiveTokenSource {
	return &ReactiveTokenSource{
		new:      new,
		t:        t,
		onUpdate: onUpdate,
	}
}

// Token returns the current token if it's still valid, else will
// refresh the current token (using r.Context for HTTP client
// information) and return the new one.
func (s *ReactiveTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		return s.t, nil
	}
	t, err := s.new.Token()
	if err != nil {
		return nil, err
	}
	s.onUpdate(t)
	s.t = t
	return t, nil
}
