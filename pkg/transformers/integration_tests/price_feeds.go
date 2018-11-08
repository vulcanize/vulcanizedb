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
		db         *postgres.DB
		blockChain core.BlockChain
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)
	})

	It("persists a ETH/USD price feed event", func() {
		blockNumber := int64(8763054)
		err := persistHeader(db, blockNumber)
		Expect(err).NotTo(HaveOccurred())
		config := price_feeds.PriceFeedConfig
		config.ContractAddresses = []string{constants.PipContractAddress}
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		transformerInitializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &price_feeds.PriceFeedConverter{},
			Repository: &price_feeds.PriceFeedRepository{},
			Fetcher:    &shared.Fetcher{},
		}
		transformer := transformerInitializer.NewLogNoteTransformer(db, blockChain)

		err = transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("207.314891143"))
		Expect(model.MedianizerAddress).To(Equal(config.ContractAddresses[0]))
	})

	It("persists a MKR/USD price feed event", func() {
		blockNumber := int64(8763059)
		err := persistHeader(db, blockNumber)
		Expect(err).NotTo(HaveOccurred())
		config := price_feeds.PriceFeedConfig
		config.ContractAddresses = []string{constants.PepContractAddress}
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		transformerInitializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &price_feeds.PriceFeedConverter{},
			Repository: &price_feeds.PriceFeedRepository{},
			Fetcher:    &shared.Fetcher{},
		}
		transformer := transformerInitializer.NewLogNoteTransformer(db, blockChain)

		err = transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("391.803979212"))
		Expect(model.MedianizerAddress).To(Equal(config.ContractAddresses[0]))
	})

	It("persists a REP/USD price feed event", func() {
		blockNumber := int64(8763062)
		err := persistHeader(db, blockNumber)
		Expect(err).NotTo(HaveOccurred())
		config := price_feeds.PriceFeedConfig
		config.ContractAddresses = []string{constants.RepContractAddress}
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		transformerInitializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &price_feeds.PriceFeedConverter{},
			Repository: &price_feeds.PriceFeedRepository{},
			Fetcher:    &shared.Fetcher{},
		}
		transformer := transformerInitializer.NewLogNoteTransformer(db, blockChain)

		err = transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("12.8169284827"))
		Expect(model.MedianizerAddress).To(Equal(config.ContractAddresses[0]))
	})
})
