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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vow_flog"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("VowFlog LogNoteTransformer", func() {
	It("transforms VowFlog log events", func() {
		blockNumber := int64(8946819)
		config := vow_flog.VowFlogConfig
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
		Expect(1).To(Equal(1))

		initializer := factories.LogNoteTransformer{
			Config:     config,
			Fetcher:    &shared.Fetcher{},
			Converter:  &vow_flog.VowFlogConverter{},
			Repository: &vow_flog.VowFlogRepository{},
		}
		transformer := initializer.NewLogNoteTransformer(db, blockChain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vow_flog.VowFlogModel
		err = db.Select(&dbResult, `SELECT era, log_idx, tx_idx from maker.vow_flog`)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult[0].Era).To(Equal("1538558052"))
		Expect(dbResult[0].LogIndex).To(Equal(uint(2)))
		Expect(dbResult[0].TransactionIndex).To(Equal(uint(2)))
	})
})
