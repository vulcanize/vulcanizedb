package flop_kick_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFlopKick(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FlopKick Suite")
}
