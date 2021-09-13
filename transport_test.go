package oauth2utils_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"

	. "github.com/koofr/go-oauth2utils"
)

var _ = Describe("Transport", func() {
	It("add the access token to the request", func() {
		var expectedAccessToken string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			go GinkgoRecover()

			Expect(r.Header.Get("Authorization")).To(Equal("Bearer " + expectedAccessToken))

			w.Write([]byte("resp"))
		}))
		defer server.Close()

		token := &oauth2.Token{
			AccessToken:  "a1",
			RefreshToken: "r1",
			Expiry:       time.Now().Add(15 * time.Minute),
		}

		onUpdateCalled := false

		reqCtx := context.WithValue(context.Background(), "foo", "bar")

		onUpdate := func(ctx context.Context, t *oauth2.Token) {
			Expect(ctx).To(Equal(reqCtx))
			onUpdateCalled = true
		}

		tokenSourceFactory := func(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
			Expect(ctx).To(Equal(reqCtx))
			return &TestTokenSource{}
		}
		s := NewReactiveTokenSource(tokenSourceFactory, token, onUpdate)

		client := &http.Client{
			Transport: &Transport{
				Base:   server.Client().Transport,
				Source: s,
			},
		}

		expectedAccessToken = token.AccessToken
		req, err := http.NewRequest("GET", server.URL, nil)
		Expect(err).NotTo(HaveOccurred())
		req = req.WithContext(reqCtx)
		res, err := client.Do(req)
		Expect(err).NotTo(HaveOccurred())
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		Expect(body).To(Equal([]byte("resp")))
		Expect(onUpdateCalled).To(BeFalse())

		token.Expiry = time.Now().Add(-15 * time.Minute)

		expectedAccessToken = "a2"
		req, err = http.NewRequest("GET", server.URL, nil)
		Expect(err).NotTo(HaveOccurred())
		req = req.WithContext(reqCtx)
		res, err = client.Do(req)
		Expect(err).NotTo(HaveOccurred())
		body, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
		Expect(body).To(Equal([]byte("resp")))
		Expect(onUpdateCalled).To(BeTrue())
	})
})
