package geth_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGeth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Geth Suite")
}
