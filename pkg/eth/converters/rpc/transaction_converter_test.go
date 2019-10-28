// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package rpc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("RPC transaction converter", func() {
	var converter rpc.RpcTransactionConverter

	BeforeEach(func() {
		converter = rpc.RpcTransactionConverter{}
	})

	It("converts hex fields to integers", func() {
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
		invalidTransaction := getFakeRpcTransaction("invalid")

		_, err := converter.ConvertRpcTransactionsToModels([]core.RpcTransaction{invalidTransaction})

		Expect(err).To(HaveOccurred())
	})

	It("copies RPC transaction hash, from, and to values to model", func() {
		rpcTransaction := getFakeRpcTransaction("0x1")

		transactionModels, err := converter.ConvertRpcTransactionsToModels([]core.RpcTransaction{rpcTransaction})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(transactionModels)).To(Equal(1))
		Expect(transactionModels[0].Hash).To(Equal(rpcTransaction.Hash))
		Expect(transactionModels[0].From).To(Equal(rpcTransaction.From))
		Expect(transactionModels[0].To).To(Equal(rpcTransaction.Recipient))
	})

	It("derives transaction RLP", func() {
		// actual transaction: https://kovan.etherscan.io/tx/0x3b29ef265425d304069c57e5145cd1c7558568b06d231775f50a693bee1aad4f
		rpcTransaction := core.RpcTransaction{
			Nonce:            "0x7aa9",
			GasPrice:         "0x3b9aca00",
			GasLimit:         "0x7a120",
			Recipient:        "0xf88bbdc1e2718f8857f30a180076ec38d53cf296",
			Amount:           "0x0",
			Payload:          "0x18178358",
			V:                "0x78",
			R:                "0x79f6a78ababfdb37b87a4d52795a49b08b5b5171443d1f2fb8f373431e77439c",
			S:                "0x3f1a210dd3b59d161735a314b88568fa91552dfe207c00a2fdbcd52ccb081409",
			Hash:             "0x3b29ef265425d304069c57e5145cd1c7558568b06d231775f50a693bee1aad4f",
			From:             "0x694032e172d9b0ee6aff5d36749bad4947a36e4e",
			TransactionIndex: "0xa",
		}

		transactionModels, err := converter.ConvertRpcTransactionsToModels([]core.RpcTransaction{rpcTransaction})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(transactionModels)).To(Equal(1))
		model := transactionModels[0]
		expectedRLP := []byte{248, 106, 130, 122, 169, 132, 59, 154, 202, 0, 131, 7, 161, 32, 148, 248, 139, 189, 193,
			226, 113, 143, 136, 87, 243, 10, 24, 0, 118, 236, 56, 213, 60, 242, 150, 128, 132, 24, 23, 131, 88, 120, 160,
			121, 246, 167, 138, 186, 191, 219, 55, 184, 122, 77, 82, 121, 90, 73, 176, 139, 91, 81, 113, 68, 61, 31, 47,
			184, 243, 115, 67, 30, 119, 67, 156, 160, 63, 26, 33, 13, 211, 181, 157, 22, 23, 53, 163, 20, 184, 133, 104,
			250, 145, 85, 45, 254, 32, 124, 0, 162, 253, 188, 213, 44, 203, 8, 20, 9}
		Expect(model.Raw).To(Equal(expectedRLP))
	})

	It("does not include transaction receipt", func() {
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
		Payload:          "0x12",
		TransactionIndex: hex,
	}
}
