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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Frob Transformer", func() {
	It("fetches and transforms a Frob event from Kovan chain", func() {
		blockNumber := int64(8935258)
		config := frob.FrobConfig
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

		initializer := frob.FrobTransformerInitializer{Config: config}
		transformer := initializer.NewFrobTransformer(db, blockchain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResult []frob.FrobModel
		err = db.Select(&dbResult, `SELECT art, dart, dink, iart, ilk, ink, urn from maker.frob`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Art).To(Equal("10000000000000000"))
		Expect(dbResult[0].Dart).To(Equal("0"))
		Expect(dbResult[0].Dink).To(Equal("10000000000000"))
		Expect(dbResult[0].IArt).To(Equal("1495509999999999999992"))
		Expect(dbResult[0].Ilk).To(Equal("ETH"))
		Expect(dbResult[0].Ink).To(Equal("10050100000000000"))
		Expect(dbResult[0].Urn).To(Equal("0xc8E093e5f3F9B5Aa6A6b33ea45960b93C161430C"))
	})

	It("unpacks an event log", func() {
		address := common.HexToAddress(shared.PitContractAddress)
		abi, err := geth.ParseAbi(shared.PitABI)
		Expect(err).NotTo(HaveOccurred())

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		entity := &frob.FrobEntity{}

		var eventLog = test_data.EthFrobLog

		err = contract.UnpackLog(entity, "Frob", eventLog)
		Expect(err).NotTo(HaveOccurred())

		expectedEntity := test_data.FrobEntity
		Expect(entity.Art).To(Equal(expectedEntity.Art))
		Expect(entity.IArt).To(Equal(expectedEntity.IArt))
		Expect(entity.Ilk).To(Equal(expectedEntity.Ilk))
		Expect(entity.Ink).To(Equal(expectedEntity.Ink))
		Expect(entity.Urn).To(Equal(expectedEntity.Urn))
	})
})
