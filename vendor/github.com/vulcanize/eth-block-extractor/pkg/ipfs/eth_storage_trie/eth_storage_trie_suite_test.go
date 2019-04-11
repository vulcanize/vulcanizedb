package eth_storage_trie_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEthStorageTrie(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EthStorageTrie Suite")
}
