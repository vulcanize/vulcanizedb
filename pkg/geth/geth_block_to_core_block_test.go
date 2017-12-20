package geth_test

import (
	"math/big"

	"context"

	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeGethClient struct{}

func (client *FakeGethClient) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	return common.HexToAddress("0x123"), nil
}

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
		client := &FakeGethClient{}
		gethBlock := geth.GethBlockToCoreBlock(block, client)

		Expect(gethBlock.Difficulty).To(Equal(difficulty.Int64()))
		Expect(gethBlock.GasLimit).To(Equal(gasLimit))
		Expect(gethBlock.GasUsed).To(Equal(gasUsed))
		Expect(gethBlock.Hash).To(Equal(block.Hash().Hex()))
		Expect(gethBlock.Nonce).To(Equal(hexutil.Encode(header.Nonce[:])))
		Expect(gethBlock.Number).To(Equal(number))
		Expect(gethBlock.ParentHash).To(Equal(block.ParentHash().Hex()))
		Expect(gethBlock.Size).To(Equal(block.Size().Int64()))
		Expect(gethBlock.Time).To(Equal(time))
		Expect(gethBlock.UncleHash).To(Equal(block.UncleHash().Hex()))
		Expect(gethBlock.IsFinal).To(BeFalse())
	})

	Describe("the converted transations", func() {
		It("is empty", func() {
			header := types.Header{}
			block := types.NewBlock(&header, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{})
			client := &FakeGethClient{}
			coreBlock := geth.GethBlockToCoreBlock(block, client)

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
			client := &FakeGethClient{}
			gethBlock := types.NewBlock(&header, []*types.Transaction{gethTransaction}, []*types.Header{}, []*types.Receipt{})
			coreBlock := geth.GethBlockToCoreBlock(gethBlock, client)

			Expect(len(coreBlock.Transactions)).To(Equal(1))
			coreTransaction := coreBlock.Transactions[0]
			Expect(coreTransaction.Data).To(Equal(gethTransaction.Data()))
			Expect(coreTransaction.To).To(Equal(gethTransaction.To().Hex()))
			Expect(coreTransaction.From).To(Equal("0x0000000000000000000000000000000000000123"))
			Expect(coreTransaction.GasLimit).To(Equal(gethTransaction.Gas().Int64()))
			Expect(coreTransaction.GasPrice).To(Equal(gethTransaction.GasPrice().Int64()))
			Expect(coreTransaction.Value).To(Equal(gethTransaction.Value().Int64()))
			Expect(coreTransaction.Nonce).To(Equal(gethTransaction.Nonce()))
		})

		It("has an empty to field when transaction creates a new contract", func() {
			gethTransaction := types.NewContractCreation(uint64(10000), big.NewInt(10), big.NewInt(5000), big.NewInt(3), []byte("1234"))
			gethBlock := types.NewBlock(&types.Header{}, []*types.Transaction{gethTransaction}, []*types.Header{}, []*types.Receipt{})
			client := &FakeGethClient{}

			coreBlock := geth.GethBlockToCoreBlock(gethBlock, client)

			coreTransaction := coreBlock.Transactions[0]
			Expect(coreTransaction.To).To(Equal(""))
		})
	})

})
