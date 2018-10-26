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

	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_grab"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat Grab Transformer", func() {
	It("transforms VatGrab log events", func() {
		blockNumber := int64(8958230)
		config := vat_grab.VatGrabConfig
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

		transformer := factories.Transformer{
			Config:     vat_grab.VatGrabConfig,
			Converter:  &vat_grab.VatGrabConverter{},
			Repository: &vat_grab.VatGrabRepository{},
			Fetcher:    &shared.Fetcher{},
		}.NewTransformer(db, blockchain)

		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vat_grab.VatGrabModel
		err = db.Select(&dbResult, `SELECT ilk, urn, v, w, dink, dart from maker.vat_grab`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Ilk).To(Equal("REP"))
		Expect(dbResult[0].Urn).To(Equal("0x6a3AE20C315E845B2E398e68EfFe39139eC6060C"))
		Expect(dbResult[0].V).To(Equal("0x2F34f22a00eE4b7a8F8BBC4eAee1658774c624e0")) //cat contract address
		Expect(dbResult[0].W).To(Equal("0x3728e9777B2a0a611ee0F89e00E01044ce4736d1"))
		expectedDink := new(big.Int)
		expectedDink.SetString("115792089237316195423570985008687907853269984665640564039455584007913129639936", 10)
		Expect(dbResult[0].Dink).To(Equal(expectedDink.String()))
		expectedDart := new(big.Int)
		expectedDart.SetString("115792089237316195423570985008687907853269984665640564039441803007913129639936", 10)
		Expect(dbResult[0].Dart).To(Equal(expectedDart.String()))
		Expect(dbResult[0].TransactionIndex).To(Equal(uint(0)))
	})
})
