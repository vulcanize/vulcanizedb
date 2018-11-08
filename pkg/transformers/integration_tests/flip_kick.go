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
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("FlipKick Transformer", func() {
	It("unpacks an event log", func() {
		address := common.HexToAddress(constants.FlipperContractAddress)
		abi, err := geth.ParseAbi(constants.FlipperABI)
		Expect(err).NotTo(HaveOccurred())

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		entity := &flip_kick.FlipKickEntity{}

		var eventLog = test_data.EthFlipKickLog

		err = contract.UnpackLog(entity, "Kick", eventLog)
		Expect(err).NotTo(HaveOccurred())

		expectedEntity := test_data.FlipKickEntity
		Expect(entity.Id).To(Equal(expectedEntity.Id))
		Expect(entity.Lot).To(Equal(expectedEntity.Lot))
		Expect(entity.Bid).To(Equal(expectedEntity.Bid))
		Expect(entity.Gal).To(Equal(expectedEntity.Gal))
		Expect(entity.End).To(Equal(expectedEntity.End))
		Expect(entity.Urn).To(Equal(expectedEntity.Urn))
		Expect(entity.Tab).To(Equal(expectedEntity.Tab))
	})

	It("fetches and transforms a FlipKick event from Kovan chain", func() {
		blockNumber := int64(8956476)
		config := flip_kick.FlipKickConfig
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

		initializer := factories.Transformer{
			Config:     flip_kick.FlipKickConfig,
			Converter:  &flip_kick.FlipKickConverter{},
			Repository: &flip_kick.FlipKickRepository{},
			Fetcher:    &shared.Fetcher{},
		}
		transformer := initializer.NewTransformer(db, blockchain)
		err = transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		var dbResult []flip_kick.FlipKickModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, "end", gal, lot FROM maker.flip_kick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("0"))
		Expect(dbResult[0].BidId).To(Equal("6"))
		Expect(dbResult[0].End.Equal(time.Unix(1538816904, 0))).To(BeTrue())
		Expect(dbResult[0].Gal).To(Equal("0x3728e9777B2a0a611ee0F89e00E01044ce4736d1"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))
	})
})
