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

package eth_test

import (
	"context"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

var (
	expectedBlock = map[string]interface{}{
		"number":           (*hexutil.Big)(mocks.MockBlock.Number()),
		"hash":             mocks.MockBlock.Hash(),
		"parentHash":       mocks.MockBlock.ParentHash(),
		"nonce":            mocks.MockBlock.Header().Nonce,
		"mixHash":          mocks.MockBlock.MixDigest(),
		"sha3Uncles":       mocks.MockBlock.UncleHash(),
		"logsBloom":        mocks.MockBlock.Bloom(),
		"stateRoot":        mocks.MockBlock.Root(),
		"miner":            mocks.MockBlock.Coinbase(),
		"difficulty":       (*hexutil.Big)(mocks.MockBlock.Difficulty()),
		"extraData":        hexutil.Bytes(mocks.MockBlock.Header().Extra),
		"gasLimit":         hexutil.Uint64(mocks.MockBlock.GasLimit()),
		"gasUsed":          hexutil.Uint64(mocks.MockBlock.GasUsed()),
		"timestamp":        hexutil.Uint64(mocks.MockBlock.Time()),
		"transactionsRoot": mocks.MockBlock.TxHash(),
		"receiptsRoot":     mocks.MockBlock.ReceiptHash(),
		"totalDifficulty":  (*hexutil.Big)(mocks.MockBlock.Difficulty()),
		"size":             hexutil.Uint64(mocks.MockBlock.Size()),
	}
	expectedHeader = map[string]interface{}{
		"number":           (*hexutil.Big)(mocks.MockBlock.Header().Number),
		"hash":             mocks.MockBlock.Header().Hash(),
		"parentHash":       mocks.MockBlock.Header().ParentHash,
		"nonce":            mocks.MockBlock.Header().Nonce,
		"mixHash":          mocks.MockBlock.Header().MixDigest,
		"sha3Uncles":       mocks.MockBlock.Header().UncleHash,
		"logsBloom":        mocks.MockBlock.Header().Bloom,
		"stateRoot":        mocks.MockBlock.Header().Root,
		"miner":            mocks.MockBlock.Header().Coinbase,
		"difficulty":       (*hexutil.Big)(mocks.MockBlock.Header().Difficulty),
		"extraData":        hexutil.Bytes(mocks.MockBlock.Header().Extra),
		"size":             hexutil.Uint64(mocks.MockBlock.Header().Size()),
		"gasLimit":         hexutil.Uint64(mocks.MockBlock.Header().GasLimit),
		"gasUsed":          hexutil.Uint64(mocks.MockBlock.Header().GasUsed),
		"timestamp":        hexutil.Uint64(mocks.MockBlock.Header().Time),
		"transactionsRoot": mocks.MockBlock.Header().TxHash,
		"receiptsRoot":     mocks.MockBlock.Header().ReceiptHash,
		"totalDifficulty":  (*hexutil.Big)(mocks.MockBlock.Header().Difficulty),
	}
	expectedTransaction = eth.NewRPCTransaction(mocks.MockTransactions[0], mocks.MockBlock.Hash(), mocks.MockBlock.NumberU64(), 0)
)

var _ = Describe("API", func() {
	var (
		db                *postgres.DB
		retriever         *eth.CIDRetriever
		fetcher           *eth.IPLDPGFetcher
		indexAndPublisher *eth.IPLDPublisherAndIndexer
		backend           *eth.Backend
		api               *eth.PublicEthAPI
	)
	BeforeEach(func() {
		var err error
		db, err = shared.SetupDB()
		Expect(err).ToNot(HaveOccurred())
		retriever = eth.NewCIDRetriever(db)
		fetcher = eth.NewIPLDPGFetcher(db)
		indexAndPublisher = eth.NewIPLDPublisherAndIndexer(db)
		backend = &eth.Backend{
			Retriever: retriever,
			Fetcher:   fetcher,
			DB:        db,
		}
		api = eth.NewPublicEthAPI(backend)
		_, err = indexAndPublisher.Publish(mocks.MockConvertedPayload)
		Expect(err).ToNot(HaveOccurred())
		uncles := mocks.MockBlock.Uncles()
		uncleHashes := make([]common.Hash, len(uncles))
		for i, uncle := range uncles {
			uncleHashes[i] = uncle.Hash()
		}
		expectedBlock["uncles"] = uncleHashes
	})
	AfterEach(func() {
		eth.TearDownDB(db)
	})
	Describe("BlockNumber", func() {
		It("Retrieves the head block number", func() {
			bn := api.BlockNumber()
			ubn := (uint64)(bn)
			subn := strconv.FormatUint(ubn, 10)
			Expect(subn).To(Equal(mocks.MockCIDPayload.HeaderCID.BlockNumber))
		})
	})

	Describe("GetTransactionByHash", func() {
		It("Retrieves the head block number", func() {
			hash := mocks.MockTransactions[0].Hash()
			tx, err := api.GetTransactionByHash(context.Background(), hash)
			Expect(err).ToNot(HaveOccurred())
			Expect(tx).To(Equal(expectedTransaction))
		})
	})

	Describe("GetBlockByNumber", func() {
		It("Retrieves a block by number", func() {
			// without full txs
			number, err := strconv.ParseInt(mocks.MockCIDPayload.HeaderCID.BlockNumber, 10, 64)
			Expect(err).ToNot(HaveOccurred())
			block, err := api.GetBlockByNumber(context.Background(), rpc.BlockNumber(number), false)
			Expect(err).ToNot(HaveOccurred())
			transactionHashes := make([]interface{}, len(mocks.MockBlock.Transactions()))
			for i, trx := range mocks.MockBlock.Transactions() {
				transactionHashes[i] = trx.Hash()
			}
			expectedBlock["transactions"] = transactionHashes
			for key, val := range expectedBlock {
				Expect(val).To(Equal(block[key]))
			}
			// with full txs
			block, err = api.GetBlockByNumber(context.Background(), rpc.BlockNumber(number), true)
			Expect(err).ToNot(HaveOccurred())
			transactions := make([]interface{}, len(mocks.MockBlock.Transactions()))
			for i, trx := range mocks.MockBlock.Transactions() {
				transactions[i] = eth.NewRPCTransactionFromBlockHash(mocks.MockBlock, trx.Hash())
			}
			expectedBlock["transactions"] = transactions
			for key, val := range expectedBlock {
				Expect(val).To(Equal(block[key]))
			}
		})
	})

	Describe("GetHeaderByNumber", func() {
		It("Retrieves a header by number", func() {
			number, err := strconv.ParseInt(mocks.MockCIDPayload.HeaderCID.BlockNumber, 10, 64)
			Expect(err).ToNot(HaveOccurred())
			header, err := api.GetHeaderByNumber(context.Background(), rpc.BlockNumber(number))
			Expect(header).To(Equal(expectedHeader))
		})
	})

	Describe("GetBlockByHash", func() {
		It("Retrieves a block by hash", func() {
			// without full txs
			block, err := api.GetBlockByHash(context.Background(), mocks.MockBlock.Hash(), false)
			Expect(err).ToNot(HaveOccurred())
			transactionHashes := make([]interface{}, len(mocks.MockBlock.Transactions()))
			for i, trx := range mocks.MockBlock.Transactions() {
				transactionHashes[i] = trx.Hash()
			}
			expectedBlock["transactions"] = transactionHashes
			for key, val := range expectedBlock {
				Expect(val).To(Equal(block[key]))
			}
			// with full txs
			block, err = api.GetBlockByHash(context.Background(), mocks.MockBlock.Hash(), true)
			Expect(err).ToNot(HaveOccurred())
			transactions := make([]interface{}, len(mocks.MockBlock.Transactions()))
			for i, trx := range mocks.MockBlock.Transactions() {
				transactions[i] = eth.NewRPCTransactionFromBlockHash(mocks.MockBlock, trx.Hash())
			}
			expectedBlock["transactions"] = transactions
			for key, val := range expectedBlock {
				Expect(val).To(Equal(block[key]))
			}
		})
	})

	Describe("GetLogs", func() {
		It("Retrieves receipt logs that match the provided topcis within the provided range", func() {
			crit := ethereum.FilterQuery{
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
					},
				},
				FromBlock: mocks.MockBlock.Number(),
				ToBlock:   mocks.MockBlock.Number(),
			}
			logs, err := api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1}))

			crit = ethereum.FilterQuery{
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
						common.HexToHash("0x05"),
					},
				},
				FromBlock: mocks.MockBlock.Number(),
				ToBlock:   mocks.MockBlock.Number(),
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(2))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1, mocks.MockLog2}))

			crit = ethereum.FilterQuery{
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
						common.HexToHash("0x06"),
					},
				},
				FromBlock: mocks.MockBlock.Number(),
				ToBlock:   mocks.MockBlock.Number(),
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1}))

			crit = ethereum.FilterQuery{
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
					},
					{
						common.HexToHash("0x07"),
					},
				},
				FromBlock: mocks.MockBlock.Number(),
				ToBlock:   mocks.MockBlock.Number(),
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(0))

			crit = ethereum.FilterQuery{
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
					},
					{
						common.HexToHash("0x06"),
					},
				},
				FromBlock: mocks.MockBlock.Number(),
				ToBlock:   mocks.MockBlock.Number(),
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1}))

			crit = ethereum.FilterQuery{
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x05"),
					},
					{
						common.HexToHash("0x07"),
					},
				},
				FromBlock: mocks.MockBlock.Number(),
				ToBlock:   mocks.MockBlock.Number(),
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog2}))

			crit = ethereum.FilterQuery{
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x05"),
					},
					{
						common.HexToHash("0x06"),
						common.HexToHash("0x07"),
					},
				},
				FromBlock: mocks.MockBlock.Number(),
				ToBlock:   mocks.MockBlock.Number(),
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog2}))

			crit = ethereum.FilterQuery{
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
						common.HexToHash("0x05"),
					},
					{
						common.HexToHash("0x06"),
						common.HexToHash("0x07"),
					},
				},
				FromBlock: mocks.MockBlock.Number(),
				ToBlock:   mocks.MockBlock.Number(),
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(2))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1, mocks.MockLog2}))

			crit = ethereum.FilterQuery{
				Topics: [][]common.Hash{
					{},
					{
						common.HexToHash("0x07"),
					},
				},
				FromBlock: mocks.MockBlock.Number(),
				ToBlock:   mocks.MockBlock.Number(),
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog2}))

			crit = ethereum.FilterQuery{
				Topics: [][]common.Hash{
					{},
					{
						common.HexToHash("0x06"),
					},
				},
				FromBlock: mocks.MockBlock.Number(),
				ToBlock:   mocks.MockBlock.Number(),
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1}))

			crit = ethereum.FilterQuery{
				Topics:    [][]common.Hash{},
				FromBlock: mocks.MockBlock.Number(),
				ToBlock:   mocks.MockBlock.Number(),
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(2))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1, mocks.MockLog2}))
		})

		It("Uses the provided blockhash if one is provided", func() {
			hash := mocks.MockBlock.Hash()
			crit := ethereum.FilterQuery{
				BlockHash: &hash,
				Topics: [][]common.Hash{
					{},
					{
						common.HexToHash("0x06"),
					},
				},
			}
			logs, err := api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1}))

			crit = ethereum.FilterQuery{
				BlockHash: &hash,
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
					},
					{
						common.HexToHash("0x06"),
					},
				},
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1}))

			crit = ethereum.FilterQuery{
				BlockHash: &hash,
				Topics: [][]common.Hash{
					{},
					{
						common.HexToHash("0x07"),
					},
				},
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog2}))

			crit = ethereum.FilterQuery{
				BlockHash: &hash,
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x05"),
					},
					{
						common.HexToHash("0x07"),
					},
				},
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog2}))

			crit = ethereum.FilterQuery{
				BlockHash: &hash,
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
					},
					{
						common.HexToHash("0x07"),
					},
				},
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(0))

			crit = ethereum.FilterQuery{
				BlockHash: &hash,
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
						common.HexToHash("0x05"),
					},
					{
						common.HexToHash("0x07"),
					},
				},
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog2}))

			crit = ethereum.FilterQuery{
				BlockHash: &hash,
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
						common.HexToHash("0x05"),
					},
				},
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(2))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1, mocks.MockLog2}))

			crit = ethereum.FilterQuery{
				BlockHash: &hash,
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
						common.HexToHash("0x05"),
					},
					{
						common.HexToHash("0x06"),
						common.HexToHash("0x07"),
					},
				},
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(2))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1, mocks.MockLog2}))

			crit = ethereum.FilterQuery{
				BlockHash: &hash,
				Topics:    [][]common.Hash{},
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(2))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1, mocks.MockLog2}))
		})

		It("Filters on contract address if any are provided", func() {
			hash := mocks.MockBlock.Hash()
			crit := ethereum.FilterQuery{
				BlockHash: &hash,
				Addresses: []common.Address{
					mocks.Address,
				},
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
						common.HexToHash("0x05"),
					},
					{
						common.HexToHash("0x06"),
						common.HexToHash("0x07"),
					},
				},
			}
			logs, err := api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(1))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1}))

			hash = mocks.MockBlock.Hash()
			crit = ethereum.FilterQuery{
				BlockHash: &hash,
				Addresses: []common.Address{
					mocks.Address,
					mocks.AnotherAddress,
				},
				Topics: [][]common.Hash{
					{
						common.HexToHash("0x04"),
						common.HexToHash("0x05"),
					},
					{
						common.HexToHash("0x06"),
						common.HexToHash("0x07"),
					},
				},
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(2))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1, mocks.MockLog2}))

			hash = mocks.MockBlock.Hash()
			crit = ethereum.FilterQuery{
				BlockHash: &hash,
				Addresses: []common.Address{
					mocks.Address,
					mocks.AnotherAddress,
				},
			}
			logs, err = api.GetLogs(context.Background(), crit)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(logs)).To(Equal(2))
			Expect(logs).To(Equal([]*types.Log{mocks.MockLog1, mocks.MockLog2}))
		})
	})
})
