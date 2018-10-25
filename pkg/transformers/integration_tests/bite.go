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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/bite"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Bite Transformer", func() {
	It("fetches and transforms a Bite event from Kovan chain", func() {
		blockNumber := int64(8956422)
		config := bite.BiteConfig
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

		initializer := bite.BiteTransformerInitializer{Config: config}
		transformer := initializer.NewBiteTransformer(db, blockchain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResult []bite.BiteModel
		err = db.Select(&dbResult, `SELECT art, iart, ilk, ink, nflip, tab, urn from maker.bite`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Art).To(Equal("149846666666666655744"))
		Expect(dbResult[0].IArt).To(Equal("1645356666666666655736"))
		Expect(dbResult[0].Ilk).To(Equal("ETH"))
		Expect(dbResult[0].Ink).To(Equal("1000000000000000000"))
		Expect(dbResult[0].NFlip).To(Equal("2"))
		Expect(dbResult[0].Tab).To(Equal("149846666666666655744"))
		Expect(dbResult[0].Urn).To(Equal("0x0000d8b4147eDa80Fec7122AE16DA2479Cbd7ffB"))
	})
})
