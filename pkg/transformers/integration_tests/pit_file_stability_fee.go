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

	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("PitFileStabilityFee LogNoteTransformer", func() {
	It("fetches and transforms a PitFileStabilityFee event from Kovan chain", func() {
		blockNumber := int64(8535544)
		config := stability_fee.StabilityFeeFileConfig
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
			Converter:  &stability_fee.PitFileStabilityFeeConverter{},
			Repository: &stability_fee.PitFileStabilityFeeRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db, blockchain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResult []stability_fee.PitFileStabilityFeeModel
		err = db.Select(&dbResult, `SELECT what, data from maker.pit_file_stability_fee`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].What).To(Equal("drip"))
		Expect(dbResult[0].Data).To(Equal("0xea29Db06E0Aa791E8ca2330D8cd9073E0760b3F1"))
	})
})
