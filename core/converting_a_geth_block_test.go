package core_test

import (
	"math/big"

	"github.com/8thlight/vulcanizedb/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Conversion of GethBlock to core.Block", func() {

	It("converts basic Block metada", func() {
		blockNumber := int64(1)
		gasUsed := int64(100000)
		gasLimit := int64(100000)
		time := int64(140000000)

		header := types.Header{
			GasUsed: big.NewInt(gasUsed),
			Number:  big.NewInt(blockNumber),
			Time:    big.NewInt(time),
			GasLimit: big.NewInt(gasLimit),
		}
		block := types.NewBlock(&header, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{})
		gethBlock := core.GethBlockToCoreBlock(block)

		Expect(gethBlock.Number).To(Equal(blockNumber))
		Expect(gethBlock.GasUsed).To(Equal(gasUsed))
		Expect(gethBlock.GasLimit).To(Equal(gasLimit))
		Expect(gethBlock.Time).To(Equal(time))
	})

	Describe("the converted transations", func() {
		It("is empty", func() {
			header := types.Header{}
			block := types.NewBlock(&header, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{})

			coreBlock := core.GethBlockToCoreBlock(block)

			Expect(len(coreBlock.Transactions)).To(Equal(0))
		})

		It("converts a single transations", func() {
			nonce := uint64(10000)
			header := types.Header{}
			to := common.Address{1}
			amount := big.NewInt(10)
			gasLimit := big.NewInt(5000)
			gasPrice := big.NewInt(3)
			payload := []byte("1234")

			gethTransaction := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, payload)
			gethBlock := types.NewBlock(&header, []*types.Transaction{gethTransaction}, []*types.Header{}, []*types.Receipt{})
			coreBlock := core.GethBlockToCoreBlock(gethBlock)

			Expect(len(coreBlock.Transactions)).To(Equal(1))
			coreTransaction := coreBlock.Transactions[0]
			Expect(coreTransaction.Data).To(Equal(gethTransaction.Data()))
			Expect(coreTransaction.To).To(Equal(gethTransaction.To().Hex()))
			Expect(coreTransaction.GasLimit).To(Equal(gethTransaction.Gas().Int64()))
			Expect(coreTransaction.GasPrice).To(Equal(gethTransaction.GasPrice().Int64()))
			Expect(coreTransaction.Value).To(Equal(gethTransaction.Value().Int64()))
			Expect(coreTransaction.Nonce).To(Equal(gethTransaction.Nonce()))
		})

		It("has an empty to field when transaction creates a new contract", func() {
			gethTransaction := types.NewContractCreation(uint64(10000), big.NewInt(10), big.NewInt(5000), big.NewInt(3), []byte("1234"))
			gethBlock := types.NewBlock(&types.Header{}, []*types.Transaction{gethTransaction}, []*types.Header{}, []*types.Receipt{})

			coreBlock := core.GethBlockToCoreBlock(gethBlock)

			coreTransaction := coreBlock.Transactions[0]
			Expect(coreTransaction.To).To(Equal(""))
		})
	})

})
