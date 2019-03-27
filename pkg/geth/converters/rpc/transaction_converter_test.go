package rpc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
)

var _ = Describe("RPC transaction converter", func() {
	It("converts hex fields to integers", func() {
		converter := rpc.RpcTransactionConverter{}
		rpcTransaction := getFakeRpcTransaction("0x1")

		transactionModels, err := converter.ConvertRpcTransactionsToModels([]core.RpcTransaction{rpcTransaction})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(transactionModels)).To(Equal(1))
		Expect(transactionModels[0].GasLimit).To(Equal(uint64(1)))
		Expect(transactionModels[0].GasPrice).To(Equal(int64(1)))
		Expect(transactionModels[0].Nonce).To(Equal(uint64(1)))
		Expect(transactionModels[0].TxIndex).To(Equal(int64(1)))
		Expect(transactionModels[0].Value).To(Equal("1"))
	})

	It("returns error if invalid hex cannot be converted", func() {
		converter := rpc.RpcTransactionConverter{}
		invalidTransaction := getFakeRpcTransaction("invalid")

		_, err := converter.ConvertRpcTransactionsToModels([]core.RpcTransaction{invalidTransaction})

		Expect(err).To(HaveOccurred())
	})

	It("copies RPC transaction hash, from, and to values to model", func() {
		converter := rpc.RpcTransactionConverter{}
		rpcTransaction := getFakeRpcTransaction("0x1")

		transactionModels, err := converter.ConvertRpcTransactionsToModels([]core.RpcTransaction{rpcTransaction})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(transactionModels)).To(Equal(1))
		Expect(transactionModels[0].Hash).To(Equal(rpcTransaction.Hash))
		Expect(transactionModels[0].From).To(Equal(rpcTransaction.From))
		Expect(transactionModels[0].To).To(Equal(rpcTransaction.Recipient))
	})

	XIt("derives transaction RLP", func() {
		// actual transaction: https://kovan.etherscan.io/tx/0x73aefdf70fc5650e0dd82affbb59d107f12dfabc50a78625b434ea68b7a69ee6
		// actual RLP hex: 0x2926af093b6b72e3f10089bde6da0f99b0d4e13354f6f37c8334efc9d7e99a47

	})

	It("does not include transaction receipt", func() {
		converter := rpc.RpcTransactionConverter{}
		rpcTransaction := getFakeRpcTransaction("0x1")

		transactionModels, err := converter.ConvertRpcTransactionsToModels([]core.RpcTransaction{rpcTransaction})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(transactionModels)).To(Equal(1))
		Expect(transactionModels[0].Receipt).To(Equal(core.Receipt{}))
	})
})

func getFakeRpcTransaction(hex string) core.RpcTransaction {
	return core.RpcTransaction{
		Hash:             "0x2",
		Amount:           hex,
		GasLimit:         hex,
		GasPrice:         hex,
		Nonce:            hex,
		From:             fakes.FakeAddress.Hex(),
		Recipient:        fakes.FakeAddress.Hex(),
		V:                "0x2",
		R:                "0x2",
		S:                "0x2",
		Payload:          nil,
		TransactionIndex: hex,
	}
}
