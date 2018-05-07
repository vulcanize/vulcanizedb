package crypto_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/crypto"
)

var _ = Describe("Public key parser", func() {
	It("parses public key from private key", func() {
		privKey := "0000000000000000000000000000000000000000000000000000000000000001"
		parser := crypto.EthPublicKeyParser{}

		pubKey, err := parser.ParsePublicKey(privKey)

		Expect(err).NotTo(HaveOccurred())
		Expect(pubKey).To(Equal("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"))
	})
})
