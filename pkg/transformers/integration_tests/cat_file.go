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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"sort"

	"github.com/ethereum/go-ethereum/ethclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/chop_lump"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/flip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/pit_vow"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Cat File transformer", func() {
	var (
		db         *postgres.DB
		blockChain core.BlockChain
		rpcClient  client.RpcClient
		err        error
		ethClient  *ethclient.Client
		fetcher    *shared.Fetcher
	)

	BeforeEach(func() {
		rpcClient, ethClient, err = getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		fetcher = shared.NewFetcher(blockChain)
	})

	// Cat contract Kovan address: 0x2f34f22a00ee4b7a8f8bbc4eaee1658774c624e0
	It("persists a chop lump event", func() {
		// transaction: 0x98574bfba4d05c3875be10d2376e678d005dbebe9a4520363407508fd21f4014
		chopLumpBlockNumber := int64(8762253)
		header, err := persistHeader(db, chopLumpBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		config := shared.TransformerConfig{
			TransformerName:     constants.CatFileChopLumpLabel,
			ContractAddresses:   []string{test_data.KovanCatContractAddress},
			ContractAbi:         test_data.KovanCatABI,
			Topic:               test_data.KovanCatFileChopLumpSignature,
			StartingBlockNumber: chopLumpBlockNumber,
			EndingBlockNumber:   chopLumpBlockNumber,
		}

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &chop_lump.CatFileChopLumpConverter{},
			Repository: &chop_lump.CatFileChopLumpRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db)

		logs, err := fetcher.FetchLogs(
			[]common.Address{common.HexToAddress(config.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []chop_lump.CatFileChopLumpModel
		err = db.Select(&dbResult, `SELECT ilk, what, data, log_idx FROM maker.cat_file_chop_lump`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(2))
		sort.Sort(byLogIndexChopLump(dbResult))

		Expect(dbResult[0].Ilk).To(Equal("5245500000000000000000000000000000000000000000000000000000000000"))
		Expect(dbResult[0].What).To(Equal("lump"))
		Expect(dbResult[0].Data).To(Equal("10000.000000000000000000"))
		Expect(dbResult[0].LogIndex).To(Equal(uint(3)))

		Expect(dbResult[1].Ilk).To(Equal("5245500000000000000000000000000000000000000000000000000000000000"))
		Expect(dbResult[1].What).To(Equal("chop"))
		Expect(dbResult[1].Data).To(Equal("1.000000000000000000000000000"))
		Expect(dbResult[1].LogIndex).To(Equal(uint(4)))
	})

	It("rechecks header for chop lump event", func() {
		// transaction: 0x98574bfba4d05c3875be10d2376e678d005dbebe9a4520363407508fd21f4014
		chopLumpBlockNumber := int64(8762253)
		header, err := persistHeader(db, chopLumpBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		config := shared.TransformerConfig{
			TransformerName:     constants.CatFileChopLumpLabel,
			ContractAddresses:   []string{test_data.KovanCatContractAddress},
			ContractAbi:         test_data.KovanCatABI,
			Topic:               test_data.KovanCatFileChopLumpSignature,
			StartingBlockNumber: chopLumpBlockNumber,
			EndingBlockNumber:   chopLumpBlockNumber,
		}

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &chop_lump.CatFileChopLumpConverter{},
			Repository: &chop_lump.CatFileChopLumpRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db)

		logs, err := fetcher.FetchLogs(
			[]common.Address{common.HexToAddress(config.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, chopLumpBlockNumber)
		Expect(err).NotTo(HaveOccurred())

		var catChopLumpChecked []int
		err = db.Select(&catChopLumpChecked, `SELECT cat_file_chop_lump_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(catChopLumpChecked[0]).To(Equal(2))
	})

	It("persists a flip event", func() {
		// transaction: 0x44bc18fdb1a5a263db114e7879653304db3e19ceb4e4496f21bc0a76c5faccbe
		flipBlockNumber := int64(8751794)
		header, err := persistHeader(db, flipBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		config := shared.TransformerConfig{
			TransformerName:     constants.CatFileFlipLabel,
			ContractAddresses:   []string{test_data.KovanCatContractAddress},
			ContractAbi:         test_data.KovanCatABI,
			Topic:               test_data.KovanCatFileFlipSignature,
			StartingBlockNumber: flipBlockNumber,
			EndingBlockNumber:   flipBlockNumber,
		}

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &flip.CatFileFlipConverter{},
			Repository: &flip.CatFileFlipRepository{},
		}

		transformer := initializer.NewLogNoteTransformer(db)

		logs, err := fetcher.FetchLogs(
			[]common.Address{common.HexToAddress(config.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []flip.CatFileFlipModel
		err = db.Select(&dbResult, `SELECT ilk, what, flip FROM maker.cat_file_flip`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Ilk).To(Equal("4554480000000000000000000000000000000000000000000000000000000000"))
		Expect(dbResult[0].What).To(Equal("flip"))
		Expect(dbResult[0].Flip).To(Equal("0x32D496Ad866D110060866B7125981C73642cc509"))
	})

	It("rechecks a flip event", func() {
		// transaction: 0x44bc18fdb1a5a263db114e7879653304db3e19ceb4e4496f21bc0a76c5faccbe
		flipBlockNumber := int64(8751794)
		header, err := persistHeader(db, flipBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		config := shared.TransformerConfig{
			TransformerName:     constants.CatFileFlipLabel,
			ContractAddresses:   []string{test_data.KovanCatContractAddress},
			ContractAbi:         test_data.KovanCatABI,
			Topic:               test_data.KovanCatFileFlipSignature,
			StartingBlockNumber: flipBlockNumber,
			EndingBlockNumber:   flipBlockNumber,
		}

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &flip.CatFileFlipConverter{},
			Repository: &flip.CatFileFlipRepository{},
		}

		transformer := initializer.NewLogNoteTransformer(db)

		logs, err := fetcher.FetchLogs(
			[]common.Address{common.HexToAddress(config.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, flipBlockNumber)
		Expect(err).NotTo(HaveOccurred())

		var catFlipChecked []int
		err = db.Select(&catFlipChecked, `SELECT cat_file_flip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(catFlipChecked[0]).To(Equal(2))
	})

	It("persists a pit vow event", func() {
		// transaction: 0x44bc18fdb1a5a263db114e7879653304db3e19ceb4e4496f21bc0a76c5faccbe
		pitVowBlockNumber := int64(8751794)
		header, err := persistHeader(db, pitVowBlockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		config := shared.TransformerConfig{
			TransformerName:     constants.CatFilePitVowLabel,
			ContractAddresses:   []string{test_data.KovanCatContractAddress},
			ContractAbi:         test_data.KovanCatABI,
			Topic:               test_data.KovanCatFilePitVowSignature,
			StartingBlockNumber: pitVowBlockNumber,
			EndingBlockNumber:   pitVowBlockNumber,
		}

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &pit_vow.CatFilePitVowConverter{},
			Repository: &pit_vow.CatFilePitVowRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db)

		logs, err := fetcher.FetchLogs(
			[]common.Address{common.HexToAddress(config.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []pit_vow.CatFilePitVowModel
		err = db.Select(&dbResult, `SELECT what, data, log_idx FROM maker.cat_file_pit_vow`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(2))
		sort.Sort(byLogIndexPitVow(dbResult))
		Expect(dbResult[0].What).To(Equal("vow"))
		Expect(dbResult[0].Data).To(Equal("0x3728e9777B2a0a611ee0F89e00E01044ce4736d1"))
		Expect(dbResult[0].LogIndex).To(Equal(uint(1)))

		Expect(dbResult[1].What).To(Equal("pit"))
		Expect(dbResult[1].Data).To(Equal("0xE7CF3198787C9A4daAc73371A38f29aAeECED87e"))
		Expect(dbResult[1].LogIndex).To(Equal(uint(2)))
	})
})

type byLogIndexChopLump []chop_lump.CatFileChopLumpModel

func (c byLogIndexChopLump) Len() int           { return len(c) }
func (c byLogIndexChopLump) Less(i, j int) bool { return c[i].LogIndex < c[j].LogIndex }
func (c byLogIndexChopLump) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type byLogIndexPitVow []pit_vow.CatFilePitVowModel

func (c byLogIndexPitVow) Len() int           { return len(c) }
func (c byLogIndexPitVow) Less(i, j int) bool { return c[i].LogIndex < c[j].LogIndex }
func (c byLogIndexPitVow) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
