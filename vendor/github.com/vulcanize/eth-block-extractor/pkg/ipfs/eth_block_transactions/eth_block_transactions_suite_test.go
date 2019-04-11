package eth_block_transactions_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEthBlockTransactions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EthBlockTransactions Suite")
}
