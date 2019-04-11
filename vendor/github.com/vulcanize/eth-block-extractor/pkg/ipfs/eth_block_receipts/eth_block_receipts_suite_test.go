package eth_block_receipts_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEthBlockReceipts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EthBlockReceipts Suite")
}
