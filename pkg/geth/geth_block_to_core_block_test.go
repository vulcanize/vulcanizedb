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

type FakeGethClient struct {
	receipts map[string]*types.Receipt
}

func NewFakeClient() *FakeGethClient {
	return &FakeGethClient{
		receipts: make(map[string]*types.Receipt),
	}
}

func (client *FakeGethClient) AddReceipts(receipts []*types.Receipt) {
	for _, receipt := range receipts {
		client.receipts[receipt.TxHash.Hex()] = receipt
	}
}

func (client *FakeGethClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	if gasUsed, ok := client.receipts[txHash.Hex()]; ok {
		return gasUsed, nil
	}
	return &types.Receipt{GasUsed: big.NewInt(0)}, nil
}

func (client *FakeGethClient) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	return common.HexToAddress("0x123"), nil
}

var _ = Describe("Conversion of GethBlock to core.Block", func() {

	It("converts basic Block metadata", func() {
		difficulty := big.NewInt(1)
		gasLimit := int64(100000)
		gasUsed := int64(100000)
		miner := common.HexToAddress("0x0000000000000000000000000000000000000123")
		extraData, _ := hexutil.Decode("0xe4b883e5bda9e7a59ee4bb99e9b1bc")
		nonce := types.BlockNonce{10}
		number := int64(1)
		time := int64(140000000)

		header := types.Header{
			Difficulty: difficulty,
			GasLimit:   big.NewInt(gasLimit),
			GasUsed:    big.NewInt(gasUsed),
			Extra:      extraData,
			Coinbase:   miner,
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
		Expect(gethBlock.Miner).To(Equal(miner.Hex()))
		Expect(gethBlock.GasUsed).To(Equal(gasUsed))
		Expect(gethBlock.Hash).To(Equal(block.Hash().Hex()))
		Expect(gethBlock.Nonce).To(Equal(hexutil.Encode(header.Nonce[:])))
		Expect(gethBlock.Number).To(Equal(number))
		Expect(gethBlock.ParentHash).To(Equal(block.ParentHash().Hex()))
		Expect(gethBlock.ExtraData).To(Equal(hexutil.Encode(block.Extra())))
		Expect(gethBlock.Size).To(Equal(block.Size().Int64()))
		Expect(gethBlock.Time).To(Equal(time))
		Expect(gethBlock.UncleHash).To(Equal(block.UncleHash().Hex()))
		Expect(gethBlock.IsFinal).To(BeFalse())
	})

	Describe("The block and uncle rewards calculations", func() {
		It("calculates block rewards for a block", func() {
			transaction := types.NewTransaction(
				uint64(226823),
				common.HexToAddress("0x108fedb097c1dcfed441480170144d8e19bb217f"),
				big.NewInt(1080900090000000000),
				big.NewInt(90000),
				big.NewInt(50000000000),
				[]byte{},
			)
			transactions := []*types.Transaction{transaction}

			txHash := transaction.Hash()
			receipt := types.Receipt{TxHash: txHash, GasUsed: big.NewInt(21000)}
			receipts := []*types.Receipt{&receipt}

			number := int64(1071819)
			header := types.Header{
				Number: big.NewInt(number),
			}
			uncles := []*types.Header{{Number: big.NewInt(1071817)}, {Number: big.NewInt(1071818)}}
			block := types.NewBlock(&header, transactions, uncles, []*types.Receipt{})

			client := NewFakeClient()
			client.AddReceipts(receipts)

			Expect(geth.CalcBlockReward(block, client)).To(Equal(5.31355))
		})

		It("calculates the uncles reward for a block", func() {
			transaction := types.NewTransaction(
				uint64(226823),
				common.HexToAddress("0x108fedb097c1dcfed441480170144d8e19bb217f"),
				big.NewInt(1080900090000000000),
				big.NewInt(90000),
				big.NewInt(50000000000),
				[]byte{})
			transactions := []*types.Transaction{transaction}

			receipt := types.Receipt{
				TxHash:  transaction.Hash(),
				GasUsed: big.NewInt(21000),
			}
			receipts := []*types.Receipt{&receipt}

			header := types.Header{
				Number: big.NewInt(int64(1071819)),
			}
			uncles := []*types.Header{
				{Number: big.NewInt(1071816)},
				{Number: big.NewInt(1071817)},
			}
			block := types.NewBlock(&header, transactions, uncles, []*types.Receipt{})

			client := NewFakeClient()
			client.AddReceipts(receipts)

			Expect(geth.CalcUnclesReward(block)).To(Equal(6.875))
		})

		It("decreases the static block reward from 5 to 3 for blocks after block 4,269,999", func() {
			transactionOne := types.NewTransaction(
				uint64(8072),
				common.HexToAddress("0xebd17720aeb7ac5186c5dfa7bafeb0bb14c02551 "),
				big.NewInt(0),
				big.NewInt(500000),
				big.NewInt(42000000000),
				[]byte{},
			)

			transactionTwo := types.NewTransaction(uint64(8071),
				common.HexToAddress("0x3cdab63d764c8c5048ed5e8f0a4e95534ba7e1ea"),
				big.NewInt(0),
				big.NewInt(500000),
				big.NewInt(42000000000),
				[]byte{})

			transactions := []*types.Transaction{transactionOne, transactionTwo}

			receiptOne := types.Receipt{
				TxHash:  transactionOne.Hash(),
				GasUsed: big.NewInt(297508),
			}
			receiptTwo := types.Receipt{
				TxHash:  transactionTwo.Hash(),
				GasUsed: big.NewInt(297508),
			}
			receipts := []*types.Receipt{&receiptOne, &receiptTwo}

			number := int64(4370055)
			header := types.Header{
				Number: big.NewInt(number),
			}
			var uncles []*types.Header
			block := types.NewBlock(&header, transactions, uncles, receipts)

			client := NewFakeClient()
			client.AddReceipts(receipts)

			Expect(geth.CalcBlockReward(block, client)).To(Equal(3.024990672))
		})
	})

	Describe("the converted transactions", func() {
		It("is empty", func() {
			header := types.Header{}
			block := types.NewBlock(&header, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{})
			client := &FakeGethClient{}
			coreBlock := geth.GethBlockToCoreBlock(block, client)

			Expect(len(coreBlock.Transactions)).To(Equal(0))
		})

		It("converts a single transaction", func() {
			nonce := uint64(10000)
			header := types.Header{}
			to := common.Address{1}
			amount := big.NewInt(10)
			gasLimit := big.NewInt(5000)
			gasPrice := big.NewInt(3)
			input := "0xf7d8c8830000000000000000000000000000000000000000000000000000000000037788000000000000000000000000000000000000000000000000000000000003bd14"
			payload, _ := hexutil.Decode(input)

			gethTransaction := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, payload)

			client := NewFakeClient()

			gethBlock := types.NewBlock(&header, []*types.Transaction{gethTransaction}, []*types.Header{}, []*types.Receipt{})
			coreBlock := geth.GethBlockToCoreBlock(gethBlock, client)

			Expect(len(coreBlock.Transactions)).To(Equal(1))
			coreTransaction := coreBlock.Transactions[0]
			Expect(coreTransaction.Data).To(Equal(input))
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

			client := NewFakeClient()
			coreBlock := geth.GethBlockToCoreBlock(gethBlock, client)

			coreTransaction := coreBlock.Transactions[0]
			Expect(coreTransaction.To).To(Equal(""))
		})
	})

})
