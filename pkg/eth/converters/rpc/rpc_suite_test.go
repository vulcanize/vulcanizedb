package rpc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRpc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rpc Suite")
}
