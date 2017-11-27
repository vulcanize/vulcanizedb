package geth_test

import (
	"path/filepath"

	cfg "github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/ethereum/go-ethereum/accounts/abi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reading ABI files", func() {

	It("loads a valid ABI file", func() {
		contractHash := "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07"
		path := filepath.Join(cfg.ProjectRoot(), "contracts", "public", contractHash+".json")

		contractAbi, err := geth.ParseAbiFile(path)

		Expect(contractAbi).NotTo(BeNil())
		Expect(err).To(BeNil())
	})

	It("returns an error when the file does not exist", func() {
		path := filepath.Join(cfg.ProjectRoot(), "contracts", "public", "missing_file.json")

		contractAbi, err := geth.ParseAbiFile(path)

		Expect(contractAbi).To(Equal(abi.ABI{}))
		Expect(err).To(Equal(geth.ErrMissingAbiFile))
	})

	It("returns an error when the file has invalid contents", func() {
		path := filepath.Join(cfg.ProjectRoot(), "pkg", "geth", "testing", "invalid_abi.json")

		contractAbi, err := geth.ParseAbiFile(path)

		Expect(contractAbi).To(Equal(abi.ABI{}))
		Expect(err).To(Equal(geth.ErrInvalidAbiFile))
	})

})
