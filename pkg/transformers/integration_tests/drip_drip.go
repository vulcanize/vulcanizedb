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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_drip"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("DripDrip Transformer", func() {
	It("transforms DripDrip log events", func() {
		blockNumber := int64(8934775)
		config := drip_drip.DripDripConfig
		config.StartingBlockNumber = blockNumber
		config.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		err = persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Converter:  &drip_drip.DripDripConverter{},
			Repository: &drip_drip.DripDripRepository{},
			Fetcher:    &shared.Fetcher{},
		}
		transformer := initializer.NewLogNoteTransformer(db, blockChain)

		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResults []drip_drip.DripDripModel
		err = db.Select(&dbResults, `SELECT ilk from maker.drip_drip`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResults)).To(Equal(1))
		dbResult := dbResults[0]
		Expect(dbResult.Ilk).To(Equal("ETH"))
	})
})
