package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Conversion of GethBlock to core.Block", func() {

	It("Converts a GethBlock to core.Block (metadata, without transactions)", func() {
		blockNumber := big.NewInt(1)
		gasUsed := big.NewInt(100000)
		gasLimit := big.NewInt(100000)
		time := big.NewInt(140000000)
		transaction := types.Transaction{}

		header := types.Header{Number: blockNumber, GasUsed: gasUsed, Time: time, GasLimit: gasLimit}
		block := types.NewBlock(&header, []*types.Transaction{&transaction}, []*types.Header{}, []*types.Receipt{})
		gethBlock := GethBlockToCoreBlock(block)

		Expect(gethBlock.Number).To(Equal(blockNumber))
		Expect(gethBlock.GasUsed).To(Equal(gasUsed))
		Expect(gethBlock.GasLimit).To(Equal(gasLimit))
		Expect(gethBlock.Time).To(Equal(time))
		Expect(gethBlock.NumberOfTransactions).To(Equal(1))
	})

})
