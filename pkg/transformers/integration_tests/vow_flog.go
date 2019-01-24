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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vow_flog"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("VowFlog LogNoteTransformer", func() {
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

	It("transforms VowFlog log events", func() {
		blockNumber := int64(8946819)
		config := shared.TransformerConfig{
			TransformerName:     constants.VowFlogLabel,
			ContractAddresses:   []string{test_data.KovanVowContractAddress},
			ContractAbi:         test_data.KovanVowABI,
			Topic:               test_data.KovanVowFlogSignature,
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		fetcher := shared.NewFetcher(blockChain)
		logs, err := fetcher.FetchLogs(
			shared.HexStringsToAddresses(config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &vow_flog.VowFlogConverter{},
			Repository: &vow_flog.VowFlogRepository{},
		}.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vow_flog.VowFlogModel
		err = db.Select(&dbResult, `SELECT era, log_idx, tx_idx from maker.vow_flog`)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult[0].Era).To(Equal("1538558052"))
		Expect(dbResult[0].LogIndex).To(Equal(uint(2)))
		Expect(dbResult[0].TransactionIndex).To(Equal(uint(2)))
	})
})
