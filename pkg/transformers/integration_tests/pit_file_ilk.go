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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("PitFileIlk LogNoteTransformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
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
		config := shared.TransformerConfig{
			TransformerName:     constants.PitFileIlkLabel,
			ContractAddresses:   []string{test_data.KovanPitContractAddress},
			ContractAbi:         test_data.KovanPitABI,
			Topic:               test_data.KovanPitFileIlkSignature,
			StartingBlockNumber: 0,
			EndingBlockNumber:   -1,
		}

		addresses = shared.HexStringsToAddresses(config.ContractAddresses)
		topics = []common.Hash{common.HexToHash(config.Topic)}

		initializer = factories.LogNoteTransformer{
			Config:     config,
			Converter:  &ilk.PitFileIlkConverter{},
			Repository: &ilk.PitFileIlkRepository{},
		}
	})

	It("fetches and transforms a Pit.file ilk 'spot' event from Kovan", func() {
		blockNumber := int64(9103223)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		fetcher := shared.NewFetcher(blockChain)
		logs, err := fetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []ilk.PitFileIlkModel
		err = db.Select(&dbResult, `SELECT ilk, what, data from maker.pit_file_ilk`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Ilk).To(Equal("0x4554480000000000000000000000000000000000000000000000000000000000"))
		Expect(dbResult[0].What).To(Equal("spot"))
		Expect(dbResult[0].Data).To(Equal("139.840000000000003410605131648"))
	})

	It("fetches and transforms a Pit.file ilk 'line' event from Kovan", func() {
		blockNumber := int64(8762253)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		fetcher := shared.NewFetcher(blockChain)
		logs, err := fetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewLogNoteTransformer(db)
		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []ilk.PitFileIlkModel
		err = db.Select(&dbResult, `SELECT ilk, what, data from maker.pit_file_ilk`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(2))
		var pitFileIlkLineModel ilk.PitFileIlkModel
		for _, result := range dbResult {
			if result.What == "line" {
				pitFileIlkLineModel = result
			}
		}
		Expect(pitFileIlkLineModel.Ilk).To(Equal("0x5245500000000000000000000000000000000000000000000000000000000000"))
		Expect(pitFileIlkLineModel.Data).To(Equal("2000000.000000000000000000"))
	})
})
