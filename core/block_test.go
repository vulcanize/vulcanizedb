package core

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Converting core.Block to DB record", func() {

	It("Converts core.Block to BlockRecord", func() {
		blockNumber := big.NewInt(1)
		gasLimit := big.NewInt(100000)
		gasUsed := big.NewInt(10)
		blockTime := big.NewInt(1508981640)
		block := Block{Number: blockNumber, GasLimit: gasLimit, GasUsed: gasUsed, Time: blockTime}
		blockRecord := BlockToBlockRecord(block)
		Expect(blockRecord.BlockNumber).To(Equal(int64(1)))
		Expect(blockRecord.GasLimit).To(Equal(int64(100000)))
		Expect(blockRecord.GasUsed).To(Equal(int64(10)))

	})
})
