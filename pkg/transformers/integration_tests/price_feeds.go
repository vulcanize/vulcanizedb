// VulcanizeDB // Copyright © 2018 Vulcanize  // This program is free software: you can redistribute it and/or modify // it under the terms of the GNU Affero General Public License as published by // the Free Software Foundation, either version 3 of the License, or // (at your option) any later version.  // This program is distributed in the hope that it will be useful, // but WITHOUT ANY WARRANTY; without even the implied warranty of // MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the // GNU Affero General Public License for more details.  // You should have received a copy of the GNU Affero General Public License // along with this program.  If not, see <http://www.gnu.org/licenses/>.

package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"

	shared_t "github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Price feeds transformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		config      shared_t.TransformerConfig
		fetcher     *shared.Fetcher
		initializer factories.LogNoteTransformer
		topics      []common.Hash
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		config = shared_t.TransformerConfig{
			TransformerName: constants.PriceFeedLabel,
			ContractAddresses: []string{
				test_data.KovanPepContractAddress,
				test_data.KovanPipContractAddress,
				test_data.KovanRepContractAddress,
			},
			ContractAbi:         test_data.KovanMedianizerABI,
			Topic:               test_data.KovanLogValueSignature,
			StartingBlockNumber: 0,
			EndingBlockNumber:   -1,
		}

		topics = []common.Hash{common.HexToHash(config.Topic)}

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
		addresses := []string{test_data.KovanPipContractAddress}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := fetcher.FetchLogs(
			shared_t.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("207.314891143000011198"))
		Expect(model.MedianizerAddress).To(Equal(addresses[0]))
	})

	It("rechecks price feed event", func() {
		blockNumber := int64(8763054)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{test_data.KovanPipContractAddress}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := fetcher.FetchLogs(
			shared.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var priceFeedChecked []int
		err = db.Select(&priceFeedChecked, `SELECT price_feeds_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(priceFeedChecked[0]).To(Equal(2))

		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("207.314891143000011198"))
		Expect(model.MedianizerAddress).To(Equal(addresses[0]))
	})

	It("persists a MKR/USD price feed event", func() {
		blockNumber := int64(8763059)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{test_data.KovanPepContractAddress}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := fetcher.FetchLogs(
			shared_t.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("391.803979212000001553"))
		Expect(model.MedianizerAddress).To(Equal(addresses[0]))
	})

	It("persists a REP/USD price feed event", func() {
		blockNumber := int64(8763062)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())
		addresses := []string{test_data.KovanRepContractAddress}
		initializer.Config.ContractAddresses = addresses
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		logs, err := fetcher.FetchLogs(
			shared_t.HexStringsToAddresses(addresses),
			topics,
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, initializer.Config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("12.816928482699999847"))
		Expect(model.MedianizerAddress).To(Equal(addresses[0]))
	})
})
