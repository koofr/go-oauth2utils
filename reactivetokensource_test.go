package oauth2utils_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"

	. "github.com/koofr/go-oauth2utils"
)

type TestTokenSource struct {
}

func (s *TestTokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken:  "a2",
		RefreshToken: "r2",
		Expiry:       time.Now().Add(1 * time.Hour),
	}

	return token, nil
}

var _ = Describe("ReactiveTokenSource", func() {
	It("should call onUpdate when token is refreshed", func() {
		token := &oauth2.Token{
			AccessToken:  "a1",
			RefreshToken: "r1",
			Expiry:       time.Now().Add(-1 * time.Minute),
		}

		ts := &TestTokenSource{}

		onUpdateCalled := false

		tokenCtx := context.WithValue(context.Background(), "foo", "bar")

		onUpdate := func(ctx context.Context, t *oauth2.Token) {
			Expect(ctx).To(Equal(tokenCtx))
			onUpdateCalled = true
		}

		tokenSourceFactory := func(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
			return ts
		}
		s := NewReactiveTokenSource(tokenSourceFactory, token, onUpdate)

		newToken, err := s.Token(tokenCtx)
		Expect(err).NotTo(HaveOccurred())

		Expect(newToken.AccessToken).To(Equal("a2"))
		Expect(newToken.RefreshToken).To(Equal("r2"))
		Expect(onUpdateCalled).To(BeTrue())

		newToken, err = s.Token(tokenCtx)
		Expect(err).NotTo(HaveOccurred())

		onUpdateCalled = false

		Expect(newToken.AccessToken).To(Equal("a2"))
		Expect(newToken.RefreshToken).To(Equal("r2"))
		Expect(onUpdateCalled).To(BeFalse())
	})
})
