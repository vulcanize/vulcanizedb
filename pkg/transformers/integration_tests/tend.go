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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
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

		config = tend.TendConfig

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
		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tend.TendModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, guy, lot FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("4000"))
		Expect(dbResult[0].BidId).To(Equal("3"))
		Expect(dbResult[0].Guy).To(Equal("0x0000d8b4147eDa80Fec7122AE16DA2479Cbd7ffB"))
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
		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tend.TendModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, guy, lot from maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("4300"))
		Expect(dbResult[0].BidId).To(Equal("3"))
		Expect(dbResult[0].Guy).To(Equal("0x0000d8b4147eDa80Fec7122AE16DA2479Cbd7ffB"))
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
		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tend.TendModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, guy, lot from maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("1000000000000000"))
		Expect(dbResult[0].BidId).To(Equal("1"))
		Expect(dbResult[0].Guy).To(Equal("0x0000d8b4147eDa80Fec7122AE16DA2479Cbd7ffB"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))

		var dbTic int64
		err = db.Get(&dbTic, `SELECT tic FROM maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		actualTic := 1538992860 + constants.TTL
		Expect(dbTic).To(Equal(actualTic))
	})
})
