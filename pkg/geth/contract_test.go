package geth_test

//import (
//	cfg "github.com/8thlight/vulcanizedb/pkg/config"
//	"github.com/8thlight/vulcanizedb/pkg/geth"
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//)
//
//var _ = Describe("The Geth blockchain", func() {
//
//	Describe("Getting a contract attribute", func() {
//		It("returns the correct attribute for a real contract", func() {
//			config, _ := cfg.NewConfig("public")
//			blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
//			contractHash := "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07"
//
//			name, err := blockchain.GetContractStateAttribute(contractHash, "name")
//
//			Expect(*name).To(Equal("OMGToken"))
//			Expect(err).To(BeNil())
//		})
//
//		It("returns an error when there is no ABI for the given contract", func() {
//			config, _ := cfg.NewConfig("public")
//			blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
//			contractHash := "MISSINGHASH"
//
//			name, err := blockchain.GetContractStateAttribute(contractHash, "name")
//
//			Expect(name).To(BeNil())
//			Expect(err).To(Equal(geth.ErrMissingAbiFile))
//		})
//
//		It("returns an error when asking for an attribute that does not exist", func() {
//			config, _ := cfg.NewConfig("public")
//			blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
//			contractHash := "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07"
//
//			name, err := blockchain.GetContractStateAttribute(contractHash, "missing_attribute")
//
//			Expect(err).To(Equal(geth.ErrInvalidStateAttribute))
//			Expect(name).To(BeNil())
//		})
//	})
//
//})
