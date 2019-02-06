package pit_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pit Suite")
}
