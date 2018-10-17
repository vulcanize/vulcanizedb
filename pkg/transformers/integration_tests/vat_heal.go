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

	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_heal"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("VatHeal Transformer", func() {
	It("transforms VatHeal log events", func() {
		blockNumber := int64(8935578)
		config := vat_heal.VatHealConfig
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockchain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockchain.Node())
		test_config.CleanTestDB(db)

		err = persistHeader(rpcClient, db, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		initializer := vat_heal.VatHealTransformerInitializer{Config: config}
		transformer := initializer.NewVatHealTransformer(db, blockchain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResults []vat_heal.VatHealModel
		err = db.Select(&dbResults, `SELECT urn, v, rad from maker.vat_heal`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.Urn).To(Equal(common.HexToAddress("0x0000000000000000000000003728e9777b2a0a611ee0f89e00e01044ce4736d1").String()))
		Expect(dbResult.V).To(Equal(common.HexToAddress("0x0000000000000000000000003728e9777b2a0a611ee0f89e00e01044ce4736d1").String()))
		Expect(dbResult.Rad).To(Equal("1000000000000000000000000000"))
	})
})
