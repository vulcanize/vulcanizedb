package ipfs_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIpfs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ipfs Suite")
}
