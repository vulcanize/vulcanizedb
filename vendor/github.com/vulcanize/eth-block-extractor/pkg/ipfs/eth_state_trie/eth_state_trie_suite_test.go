package eth_state_trie_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEthStateTrie(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EthStateTrie Suite")
}
