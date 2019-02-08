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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_tune"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("VatTune LogNoteTransformer", func() {
	It("transforms VatTune log events", func() {
		blockNumber := int64(8761670)
		config := shared.TransformerConfig{
			TransformerName:     constants.VatTuneLabel,
			ContractAddresses:   []string{test_data.KovanVatContractAddress},
			ContractAbi:         test_data.KovanVatABI,
			Topic:               test_data.KovanVatTuneSignature,
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

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
			Converter:  &vat_tune.VatTuneConverter{},
			Repository: &vat_tune.VatTuneRepository{},
		}.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vat_tune.VatTuneModel
		err = db.Select(&dbResult, `SELECT ilk, urn, v, w, dink, dart from maker.vat_tune`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Ilk).To(Equal("4554480000000000000000000000000000000000000000000000000000000000"))
		Expect(dbResult[0].Urn).To(Equal("0x4F26FfBe5F04ED43630fdC30A87638d53D0b0876"))
		Expect(dbResult[0].V).To(Equal("0x4F26FfBe5F04ED43630fdC30A87638d53D0b0876"))
		Expect(dbResult[0].W).To(Equal("0x4F26FfBe5F04ED43630fdC30A87638d53D0b0876"))
		Expect(dbResult[0].Dink).To(Equal("0"))
		expectedDart := new(big.Int)
		expectedDart.SetString("115792089237316195423570985008687907853269984665640564039455584007913129639936", 10)
		Expect(dbResult[0].Dart).To(Equal(expectedDart.String()))
		Expect(dbResult[0].TransactionIndex).To(Equal(uint(0)))
	})

	It("transforms VatTune log events", func() {
		blockNumber := int64(8761670)
		config := shared.TransformerConfig{
			TransformerName:     constants.VatTuneLabel,
			ContractAddresses:   []string{test_data.KovanVatContractAddress},
			ContractAbi:         test_data.KovanVatABI,
			Topic:               test_data.KovanVatTuneSignature,
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

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
			Converter:  &vat_tune.VatTuneConverter{},
			Repository: &vat_tune.VatTuneRepository{},
		}.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var vatTuneChecked []int
		err = db.Select(&vatTuneChecked, `SELECT vat_tune_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(vatTuneChecked[0]).To(Equal(2))

		var dbResult []vat_tune.VatTuneModel
		err = db.Select(&dbResult, `SELECT ilk, urn, v, w, dink, dart from maker.vat_tune`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Ilk).To(Equal("4554480000000000000000000000000000000000000000000000000000000000"))
		Expect(dbResult[0].Urn).To(Equal("0x4F26FfBe5F04ED43630fdC30A87638d53D0b0876"))
		Expect(dbResult[0].V).To(Equal("0x4F26FfBe5F04ED43630fdC30A87638d53D0b0876"))
		Expect(dbResult[0].W).To(Equal("0x4F26FfBe5F04ED43630fdC30A87638d53D0b0876"))
		Expect(dbResult[0].Dink).To(Equal("0"))
		expectedDart := new(big.Int)
		expectedDart.SetString("115792089237316195423570985008687907853269984665640564039455584007913129639936", 10)
		Expect(dbResult[0].Dart).To(Equal(expectedDart.String()))
		Expect(dbResult[0].TransactionIndex).To(Equal(uint(0)))
	})
})
