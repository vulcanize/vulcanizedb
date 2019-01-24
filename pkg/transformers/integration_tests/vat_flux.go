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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_flux"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("VatFlux LogNoteTransformer", func() {
	It("transforms VatFlux log events", func() {
		blockNumber := int64(9004474)
		config := shared.TransformerConfig{
			TransformerName:     constants.VatFluxLabel,
			ContractAddresses:   []string{test_data.KovanVatContractAddress},
			ContractAbi:         test_data.KovanVatABI,
			Topic:               test_data.KovanVatFluxSignature,
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

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &vat_flux.VatFluxConverter{},
			Repository: &vat_flux.VatFluxRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db)

		err = transformer.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vat_flux.VatFluxModel
		err = db.Select(&dbResult, `SELECT ilk, src, dst, rad from maker.vat_flux`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Ilk).To(Equal("REP"))
		Expect(dbResult[0].Src).To(Equal("0xC0851F73CC8DD5c0765E71980eC7E7Fd1EF74434"))
		Expect(dbResult[0].Dst).To(Equal("0x0000d8b4147eDa80Fec7122AE16DA2479Cbd7ffB"))
		Expect(dbResult[0].Rad).To(Equal("1800000000000000000000000000000000000000000000"))
		Expect(dbResult[0].TransactionIndex).To(Equal(uint(0)))
	})
})
