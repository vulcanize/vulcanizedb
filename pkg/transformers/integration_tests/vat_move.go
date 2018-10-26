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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_move"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("VatMove LogNoteTransformer", func() {
	It("transforms VatMove log events", func() {
		blockNumber := int64(9004628)
		config := vat_move.VatMoveConfig
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

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Fetcher:    &shared.Fetcher{},
			Converter:  &vat_move.VatMoveConverter{},
			Repository: &vat_move.VatMoveRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db, blockchain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vat_move.VatMoveModel
		err = db.Select(&dbResults, `SELECT src, dst, rad from maker.vat_move`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.Src).To(Equal(common.HexToAddress("0x8868bad8e74fca4505676d1b5b21ecc23328d132").String()))
		Expect(dbResult.Dst).To(Equal(common.HexToAddress("0x0000d8b4147eda80fec7122ae16da2479cbd7ffb").String()))
		Expect(dbResult.Rad).To(Equal("1000000000000000000000000000000000000000000000"))
	})
})
