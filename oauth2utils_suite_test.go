package oauth2utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOAuth2utils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OAuth2utils Suite")
}
