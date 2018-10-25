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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Tend Transformer", func() {
	It("fetches and transforms a Tend event from Kovan chain", func() {
		blockNumber := int64(8935601)
		config := tend.TendConfig
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockchain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockchain.Node())
		test_config.CleanTestDB(db)

		err = persistHeader(db, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		initializer := factories.Transformer{
			Config:     config,
			Fetcher:    &shared.Fetcher{},
			Converter:  &tend.TendConverter{},
			Repository: &tend.TendRepository{},
		}
		transformer := initializer.NewTransformer(db, blockchain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tend.TendModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, guy, lot, tic from maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("4000"))
		Expect(dbResult[0].BidId).To(Equal("3"))
		Expect(dbResult[0].Guy).To(Equal("0x0000d8b4147eDa80Fec7122AE16DA2479Cbd7ffB"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))
		Expect(dbResult[0].Tic).To(Equal("0"))
	})

	It("fetches and transforms a subsequent Tend event from Kovan chain for the same auction", func() {
		blockNumber := int64(8935731)
		config := tend.TendConfig
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockchain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockchain.Node())
		test_config.CleanTestDB(db)

		err = persistHeader(db, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		initializer := factories.Transformer{
			Config:     config,
			Fetcher:    &shared.Fetcher{},
			Converter:  &tend.TendConverter{},
			Repository: &tend.TendRepository{},
		}
		transformer := initializer.NewTransformer(db, blockchain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResult []tend.TendModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, guy, lot, tic from maker.tend`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("4300"))
		Expect(dbResult[0].BidId).To(Equal("3"))
		Expect(dbResult[0].Guy).To(Equal("0x0000d8b4147eDa80Fec7122AE16DA2479Cbd7ffB"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))
		Expect(dbResult[0].Tic).To(Equal("0"))
	})
})
