package geth_test

import (
	"math/big"

	"github.com/8thlight/vulcanizedb/geth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
)

var _ = Describe("Conversion of GethBlock to core.Block", func() {

	It("converts basic Block metada", func() {
		difficulty := big.NewInt(1)
		gasLimit := int64(100000)
		gasUsed := int64(100000)
		nonce := types.BlockNonce{10}
		number := int64(1)
		time := int64(140000000)

		header := types.Header{
			Difficulty: difficulty,
			GasLimit:   big.NewInt(gasLimit),
			GasUsed:    big.NewInt(gasUsed),
			Nonce:      nonce,
			Number:     big.NewInt(number),
			ParentHash: common.Hash{64},
			Time:       big.NewInt(time),
			UncleHash:  common.Hash{128},
		}
		block := types.NewBlock(&header, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{})
		gethBlock := geth.GethBlockToCoreBlock(block)

		Expect(gethBlock.Difficulty).To(Equal(difficulty.Int64()))
		Expect(gethBlock.GasLimit).To(Equal(gasLimit))
		Expect(gethBlock.GasUsed).To(Equal(gasUsed))
		Expect(gethBlock.Hash).To(Equal(block.Hash().Hex()))
		Expect(gethBlock.Nonce).To(Equal((strconv.FormatUint(block.Nonce(), 10))))
		Expect(gethBlock.Number).To(Equal(number))
		Expect(gethBlock.ParentHash).To(Equal(block.ParentHash().Hex()))
		Expect(gethBlock.Size).To(Equal(block.Size().Int64()))
		Expect(gethBlock.Time).To(Equal(time))
		Expect(gethBlock.UncleHash).To(Equal(block.UncleHash().Hex()))
	})

	Describe("the converted transations", func() {
		It("is empty", func() {
			header := types.Header{}
			block := types.NewBlock(&header, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{})

			coreBlock := geth.GethBlockToCoreBlock(block)

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
			coreBlock := geth.GethBlockToCoreBlock(gethBlock)

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

			coreBlock := geth.GethBlockToCoreBlock(gethBlock)

			coreTransaction := coreBlock.Transactions[0]
			Expect(coreTransaction.To).To(Equal(""))
		})
	})

})
