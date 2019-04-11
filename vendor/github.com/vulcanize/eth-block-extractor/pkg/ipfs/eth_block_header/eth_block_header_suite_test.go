package eth_block_header_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEthBlockHeader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EthBlockHeader Suite")
}
