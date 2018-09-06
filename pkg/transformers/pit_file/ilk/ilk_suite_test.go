package ilk_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIlk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ilk Suite")
}
