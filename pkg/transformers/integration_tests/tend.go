// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Tend LogNoteTransformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		config      shared.TransformerConfig
		fetcher     *shared.Fetcher
		initializer factories.LogNoteTransformer
		addresses   []common.Address
		topics      []common.Hash
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		config = shared.TransformerConfig{
			TransformerName:     constants.TendLabel,
			ContractAddresses:   []string{test_data.KovanFlapperContractAddress, test_data.KovanFlipperContractAddress},
			ContractAbi:         test_data.KovanFlipperABI,
			Topic:               test_data.KovanTendFunctionSignature,
			StartingBlockNumber: 0,
			EndingBlockNumber:   -1,
		}

		fetcher = shared.NewFetcher(blockChain)
		addresses = shared.HexStringsToAddresses(config.ContractAddresses)
		topics = []common.Hash{common.HexToHash(config.Topic)}

		initializer = factories.LogNoteTransformer{
			Config:     config,
			Converter:  &tend.TendConverter{},
			Repository: &tend.TendRepository{},
		}
	})

	It("fetches and transforms a Flip Tend event from Kovan chain", func() {
		blockNumber := int64(8935601)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := fetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tend.TendModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, guy, lot FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("4000"))
		Expect(dbResult[0].BidId).To(Equal("3"))
		Expect(dbResult[0].Guy).To(Equal("0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))

		var dbTic int64
		err = db.Get(&dbTic, `SELECT tic FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		actualTic := 1538490276 + constants.TTL
		Expect(dbTic).To(Equal(actualTic))
	})

	It("rechecks tend event", func() {
		blockNumber := int64(8935601)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := fetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var tendChecked []int
		err = db.Select(&tendChecked, `SELECT tend_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(tendChecked[0]).To(Equal(2))

		var dbResult []tend.TendModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, guy, lot from maker.tend where header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("4000"))
		Expect(dbResult[0].BidId).To(Equal("3"))
		Expect(dbResult[0].Guy).To(Equal("0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))

		var dbTic int64
		err = db.Get(&dbTic, `SELECT tic FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		actualTic := 1538490276 + constants.TTL
		Expect(dbTic).To(Equal(actualTic))
	})

	It("fetches and transforms a subsequent Flip Tend event from Kovan chain for the same auction", func() {
		blockNumber := int64(8935731)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := fetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tend.TendModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, guy, lot from maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("4300"))
		Expect(dbResult[0].BidId).To(Equal("3"))
		Expect(dbResult[0].Guy).To(Equal("0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))

		var dbTic int64
		err = db.Get(&dbTic, `SELECT tic FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		actualTic := 1538491224 + constants.TTL
		Expect(dbTic).To(Equal(actualTic))
	})

	It("fetches and transforms a Flap Tend event from the Kovan chain", func() {
		blockNumber := int64(9003177)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := fetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tend.TendModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, guy, lot from maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("1000000000000000"))
		Expect(dbResult[0].BidId).To(Equal("1"))
		Expect(dbResult[0].Guy).To(Equal("0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))

		var dbTic int64
		err = db.Get(&dbTic, `SELECT tic FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		actualTic := 1538992860 + constants.TTL
		Expect(dbTic).To(Equal(actualTic))
	})
})
