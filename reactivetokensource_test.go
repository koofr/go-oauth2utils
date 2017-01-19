package oauth2utils_test

import (
	"time"

	. "github.com/koofr/go-oauth2utils"
	"golang.org/x/oauth2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	It("write to close chanel when ReadCloser is closed", func() {
		token := &oauth2.Token{
			AccessToken:  "a1",
			RefreshToken: "r1",
			Expiry:       time.Now().Add(-1 * time.Minute),
		}

		ts := &TestTokenSource{}

		onUpdateCalled := false

		onUpdate := func(t *oauth2.Token) {
			onUpdateCalled = true
		}

		s := NewReactiveTokenSource(token, ts, onUpdate)

		newToken, err := s.Token()
		Expect(err).NotTo(HaveOccurred())

		Expect(newToken.AccessToken).To(Equal("a2"))
		Expect(newToken.RefreshToken).To(Equal("r2"))
		Expect(onUpdateCalled).To(BeTrue())

		newToken, err = s.Token()
		Expect(err).NotTo(HaveOccurred())

		onUpdateCalled = false

		Expect(newToken.AccessToken).To(Equal("a2"))
		Expect(newToken.RefreshToken).To(Equal("r2"))
		Expect(onUpdateCalled).To(BeFalse())
	})
})
