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
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_tune"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("VatTune Transformer", func() {
	It("transforms VatTune log events", func() {
		blockNumber := int64(8761670)
		config := vat_tune.VatTuneConfig
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		ipc := "https://kovan.infura.io/J5Vd2fRtGsw0zZ0Ov3BL"
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockchain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockchain.Node())
		test_config.CleanTestDB(db)

		err = persistHeader(rpcClient, db, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		initializer := vat_tune.VatTuneTransformerInitializer{Config: config}
		transformer := initializer.NewVatTuneTransformer(db, blockchain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vat_tune.VatTuneModel
		err = db.Select(&dbResult, `SELECT ilk, urn, v, w, dink, dart from maker.vat_tune`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Ilk).To(Equal("ETH"))
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
