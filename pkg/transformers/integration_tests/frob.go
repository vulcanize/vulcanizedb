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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"

	shared_t "github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Frob Transformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		fetcher     *shared.Fetcher
		config      shared_t.TransformerConfig
		initializer factories.Transformer
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		fetcher = shared.NewFetcher(blockChain)
		config = shared_t.TransformerConfig{
			TransformerName:     constants.FrobLabel,
			ContractAddresses:   []string{test_data.KovanPitContractAddress},
			ContractAbi:         test_data.KovanPitABI,
			Topic:               test_data.KovanFrobSignature,
			StartingBlockNumber: 0,
			EndingBlockNumber:   -1,
		}

		initializer = factories.Transformer{
			Config:     config,
			Converter:  &frob.FrobConverter{},
			Repository: &frob.FrobRepository{},
		}
	})

	It("fetches and transforms a Frob event from Kovan chain", func() {
		blockNumber := int64(8935258)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := fetcher.FetchLogs(
			shared_t.HexStringsToAddresses(config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []frob.FrobModel
		err = db.Select(&dbResult, `SELECT art, dart, dink, iart, ilk, ink, urn from maker.frob`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Art).To(Equal("10000000000000000"))
		Expect(dbResult[0].Dart).To(Equal("0"))
		Expect(dbResult[0].Dink).To(Equal("10000000000000"))
		Expect(dbResult[0].IArt).To(Equal("1495509999999999999992"))
		ilkID, err := shared.GetOrCreateIlk("4554480000000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].Ink).To(Equal("10050100000000000"))
		Expect(dbResult[0].Urn).To(Equal("000000000000000000000000c8e093e5f3f9b5aa6a6b33ea45960b93c161430c"))
	})

	It("rechecks frob event", func() {
		blockNumber := int64(8935258)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := fetcher.FetchLogs(
			shared.HexStringsToAddresses(config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var frobChecked []int
		err = db.Select(&frobChecked, `SELECT frob_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(frobChecked[0]).To(Equal(2))

		var dbResult []frob.FrobModel
		err = db.Select(&dbResult, `SELECT art, dart, dink, iart, ilk, ink, urn from maker.frob`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Art).To(Equal("10000000000000000"))
		Expect(dbResult[0].Dart).To(Equal("0"))
		Expect(dbResult[0].Dink).To(Equal("10000000000000"))
		Expect(dbResult[0].IArt).To(Equal("1495509999999999999992"))
		ilkID, err := shared.GetOrCreateIlk("4554480000000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Ilk).To(Equal(strconv.Itoa(ilkID)))
		Expect(dbResult[0].Ink).To(Equal("10050100000000000"))
		Expect(dbResult[0].Urn).To(Equal("000000000000000000000000c8e093e5f3f9b5aa6a6b33ea45960b93c161430c"))
	})

	It("unpacks an event log", func() {
		address := common.HexToAddress(test_data.KovanPitContractAddress)
		abi, err := geth.ParseAbi(test_data.KovanPitABI)
		Expect(err).NotTo(HaveOccurred())

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		entity := &frob.FrobEntity{}

		var eventLog = test_data.EthFrobLog

		err = contract.UnpackLog(entity, "Frob", eventLog)
		Expect(err).NotTo(HaveOccurred())

		expectedEntity := test_data.FrobEntity
		Expect(entity.Art).To(Equal(expectedEntity.Art))
		Expect(entity.IArt).To(Equal(expectedEntity.IArt))
		Expect(entity.Ilk).To(Equal(expectedEntity.Ilk))
		Expect(entity.Ink).To(Equal(expectedEntity.Ink))
		Expect(entity.Urn).To(Equal(expectedEntity.Urn))
	})
})
