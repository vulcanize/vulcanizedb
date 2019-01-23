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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/debt_ceiling"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("PitFileDebtCeiling LogNoteTransformer", func() {
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

	It("fetches and transforms a PitFileDebtCeiling event from Kovan chain", func() {
		blockNumber := int64(8535578)
		config := shared.TransformerConfig{
			TransformerName:     constants.PitFileDebtCeilingLabel,
			ContractAddresses:   []string{test_data.KovanPitContractAddress},
			ContractAbi:         test_data.KovanPitABI,
			Topic:               test_data.KovanPitFileDebtCeilingSignature,
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

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &debt_ceiling.PitFileDebtCeilingConverter{},
			Repository: &debt_ceiling.PitFileDebtCeilingRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []debt_ceiling.PitFileDebtCeilingModel
		err = db.Select(&dbResult, `SELECT what, data from maker.pit_file_debt_ceiling`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].What).To(Equal("Line"))
		Expect(dbResult[0].Data).To(Equal("10000000.000000000000000000"))
	})
})
