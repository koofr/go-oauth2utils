package oauth2utils

import (
	"context"
	"sync"

	"golang.org/x/oauth2"
)

type TokenSourceFactory func(ctx context.Context, t *oauth2.Token) oauth2.TokenSource

// ReactiveTokenSource is a TokenSource that holds a single token in memory
// and validates its expiry before each call to retrieve it with
// Token. If it's expired, it will be auto-refreshed using the
// new TokenSource and onUpdate callback will be called with new token.
type ReactiveTokenSource struct {
	tokenSourceFactory TokenSourceFactory

	mu sync.Mutex // guards t
	t  *oauth2.Token

	onUpdate func(ctx context.Context, t *oauth2.Token)
}

// NewReactiveTokenSource returns a ReactiveTokenSource which repeatedly
// returns the same token as long as it's valid, starting with t.
// When its cached token is invalid, a new token is obtained from src
// and onUpdate callback is called.
func NewReactiveTokenSource(
	tokenSourceFactory TokenSourceFactory,
	t *oauth2.Token,
	onUpdate func(ctx context.Context, t *oauth2.Token),
) *ReactiveTokenSource {
	return &ReactiveTokenSource{
		tokenSourceFactory: tokenSourceFactory,
		t:                  t,
		onUpdate:           onUpdate,
	}
}

// Token returns the current token if it's still valid, else will
// refresh the current token (using r.Context for HTTP client
// information) and return the new one.
func (s *ReactiveTokenSource) Token(ctx context.Context) (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		return s.t, nil
	}
	tokenSource := s.tokenSourceFactory(ctx, s.t)
	t, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	s.onUpdate(ctx, t)
	s.t = t
	return t, nil
}
