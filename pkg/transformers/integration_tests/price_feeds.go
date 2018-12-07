// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Price feeds transformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		config      shared.TransformerConfig
		fetcher     shared.Fetcher
		initializer factories.LogNoteTransformer
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)
		config = price_feeds.PriceFeedConfig
		fetcher = shared.NewFetcher(blockChain)

		initializer = factories.LogNoteTransformer{
			Config:     config,
			Converter:  &price_feeds.PriceFeedConverter{},
			Repository: &price_feeds.PriceFeedRepository{},
		}
	})

	It("persists a ETH/USD price feed event", func() {
		blockNumber := int64(8763054)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		initializer.Config.ContractAddresses = []string{constants.PipContractAddress}
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := fetcher.FetchLogs(
			shared.HexStringsToAddresses(initializer.Config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header)

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("207.314891143000011198"))
		Expect(model.MedianizerAddress).To(Equal(config.ContractAddresses[0]))
	})

	It("persists a MKR/USD price feed event", func() {
		blockNumber := int64(8763059)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		initializer.Config.ContractAddresses = []string{constants.PepContractAddress}
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := fetcher.FetchLogs(
			shared.HexStringsToAddresses(initializer.Config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header)

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("391.803979212000001553"))
		Expect(model.MedianizerAddress).To(Equal(config.ContractAddresses[0]))
	})

	It("persists a REP/USD price feed event", func() {
		blockNumber := int64(8763062)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		initializer.Config.ContractAddresses = []string{constants.RepContractAddress}
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := fetcher.FetchLogs(
			shared.HexStringsToAddresses(config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header)

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("12.816928482699999847"))
		Expect(model.MedianizerAddress).To(Equal(config.ContractAddresses[0]))
	})
})
