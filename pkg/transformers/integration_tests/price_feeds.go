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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Price feeds transformer", func() {
	var (
		db         *postgres.DB
		blockChain core.BlockChain
	)

	BeforeEach(func() {
		ipc := "https://kovan.infura.io/J5Vd2fRtGsw0zZ0Ov3BL"
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		for i := 8763054; i < 8763063; i++ {
			err = persistHeader(rpcClient, db, int64(i))
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("persists a ETH/USD price feed event", func() {
		config := price_feeds.IPriceFeedConfig{
			ContractAddresses:   []string{"0x9FfFE440258B79c5d6604001674A4722FfC0f7Bc"},
			StartingBlockNumber: 8763054,
			EndingBlockNumber:   8763054,
		}
		transformerInitializer := price_feeds.PriceFeedTransformerInitializer{Config: config}
		transformer := transformerInitializer.NewPriceFeedTransformer(db, blockChain)

		err := transformer.Execute()

		time.AfterFunc(5*time.Second, func() {
			defer GinkgoRecover()
			Expect(err).NotTo(HaveOccurred())
			var model price_feeds.PriceFeedModel
			err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, config.StartingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(model.UsdValue).To(Equal("207.314891143"))
		})
	})

	It("persists a MKR/USD price feed event", func() {
		config := price_feeds.IPriceFeedConfig{
			ContractAddresses:   []string{"0xB1997239Cfc3d15578A3a09730f7f84A90BB4975"},
			StartingBlockNumber: 8763059,
			EndingBlockNumber:   8763059,
		}
		transformerInitializer := price_feeds.PriceFeedTransformerInitializer{Config: config}
		transformer := transformerInitializer.NewPriceFeedTransformer(db, blockChain)

		err := transformer.Execute()

		time.AfterFunc(5*time.Second, func() {
			defer GinkgoRecover()
			Expect(err).NotTo(HaveOccurred())
			var model price_feeds.PriceFeedModel
			err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, config.StartingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(model.UsdValue).To(Equal("391.803979212"))
		})
	})

	It("persists a REP/USD price feed event", func() {
		config := price_feeds.IPriceFeedConfig{
			ContractAddresses:   []string{"0xf88bBDc1E2718F8857F30A180076ec38d53cf296"},
			StartingBlockNumber: 8763062,
			EndingBlockNumber:   8763062,
		}
		transformerInitializer := price_feeds.PriceFeedTransformerInitializer{Config: config}
		transformer := transformerInitializer.NewPriceFeedTransformer(db, blockChain)

		err := transformer.Execute()

		time.AfterFunc(5*time.Second, func() {
			defer GinkgoRecover()
			Expect(err).NotTo(HaveOccurred())
			var model price_feeds.PriceFeedModel
			err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, config.StartingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(model.UsdValue).To(Equal("12.8169284827"))
		})
	})
})
