package geth_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

var _ = Describe("Conversion of GethReceipt to core.Receipt", func() {

	It(`converts geth receipt to internal receipt format (pre Byzantium has post-transaction stateroot)`, func() {
		receipt := types.Receipt{
			Bloom:             types.BytesToBloom(hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")),
			ContractAddress:   common.Address{},
			CumulativeGasUsed: big.NewInt(21000),
			GasUsed:           big.NewInt(21000),
			Logs:              []*types.Log{},
			PostState:         hexutil.MustDecode("0x88abf7e73128227370aa7baa3dd4e18d0af70e92ef1f9ef426942fbe2dddb733"),
			TxHash:            common.HexToHash("0x97d99bc7729211111a21b12c933c949d4f31684f1d6954ff477d0477538ff017"),
		}

		expected := core.Receipt{
			Bloom:             "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			ContractAddress:   "",
			CumulativeGasUsed: 21000,
			GasUsed:           21000,
			Logs:              []core.Log{},
			StateRoot:         "0x88abf7e73128227370aa7baa3dd4e18d0af70e92ef1f9ef426942fbe2dddb733",
			Status:            -99,
			TxHash:            receipt.TxHash.Hex(),
		}

		coreReceipt := geth.ReceiptToCoreReceipt(&receipt)
		Expect(coreReceipt.Bloom).To(Equal(expected.Bloom))
		Expect(coreReceipt.ContractAddress).To(Equal(expected.ContractAddress))
		Expect(coreReceipt.CumulativeGasUsed).To(Equal(expected.CumulativeGasUsed))
		Expect(coreReceipt.GasUsed).To(Equal(expected.GasUsed))
		Expect(coreReceipt.Logs).To(Equal(expected.Logs))
		Expect(coreReceipt.StateRoot).To(Equal(expected.StateRoot))
		Expect(coreReceipt.Status).To(Equal(expected.Status))
		Expect(coreReceipt.TxHash).To(Equal(expected.TxHash))

	})

	It("converts geth receipt to internal receipt format (post Byzantium has status", func() {
		receipt := types.Receipt{
			Bloom:             types.BytesToBloom(hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")),
			ContractAddress:   common.HexToAddress("x0123"),
			CumulativeGasUsed: big.NewInt(7996119),
			GasUsed:           big.NewInt(21000),
			Logs:              []*types.Log{},
			Status:            uint(1),
			TxHash:            common.HexToHash("0xe340558980f89d5f86045ac11e5cc34e4bcec20f9f1e2a427aa39d87114e8223"),
		}

		expected := core.Receipt{
			Bloom:             "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			ContractAddress:   receipt.ContractAddress.Hex(),
			CumulativeGasUsed: 7996119,
			GasUsed:           21000,
			Logs:              []core.Log{},
			StateRoot:         "",
			Status:            1,
			TxHash:            receipt.TxHash.Hex(),
		}

		coreReceipt := geth.ReceiptToCoreReceipt(&receipt)
		Expect(coreReceipt.Bloom).To(Equal(expected.Bloom))
		Expect(coreReceipt.ContractAddress).To(Equal(""))
		Expect(coreReceipt.CumulativeGasUsed).To(Equal(expected.CumulativeGasUsed))
		Expect(coreReceipt.GasUsed).To(Equal(expected.GasUsed))
		Expect(coreReceipt.Logs).To(Equal(expected.Logs))
		Expect(coreReceipt.StateRoot).To(Equal(expected.StateRoot))
		Expect(coreReceipt.Status).To(Equal(expected.Status))
		Expect(coreReceipt.TxHash).To(Equal(expected.TxHash))

	})

})
