// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package integration_tests

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("FlopKick Transformer", func() {
	var (
		db          *postgres.DB
		blockChain  core.BlockChain
		config      shared.TransformerConfig
		initializer factories.Transformer
		fetcher     shared.LogFetcher
		addresses   []common.Address
		topics      []common.Hash
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		config = shared.TransformerConfig{
			TransformerName:     constants.FlopKickLabel,
			ContractAddresses:   []string{test_data.KovanFlopperContractAddress},
			ContractAbi:         test_data.KovanFlopperABI,
			Topic:               test_data.KovanFlopKickSignature,
			StartingBlockNumber: 0,
			EndingBlockNumber:   -1,
		}

		initializer = factories.Transformer{
			Config:     config,
			Converter:  &flop_kick.FlopKickConverter{},
			Repository: &flop_kick.FlopKickRepository{},
		}

		fetcher = shared.NewFetcher(blockChain)
		addresses = shared.HexStringsToAddresses(config.ContractAddresses)
		topics = []common.Hash{common.HexToHash(config.Topic)}
	})

	It("fetches and transforms a FlopKick event from Kovan chain", func() {
		blockNumber := int64(8672119)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := fetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []flop_kick.Model
		err = db.Select(&dbResult, `SELECT bid, bid_id, "end", gal, lot FROM maker.flop_kick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("0"))
		Expect(dbResult[0].BidId).To(Equal("1"))
		Expect(dbResult[0].End.Equal(time.Unix(1536726768, 0))).To(BeTrue())
		Expect(dbResult[0].Gal).To(Equal("0x9B870D55BaAEa9119dBFa71A92c5E26E79C4726d"))
		// this very large number appears to be derived from the data including: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		Expect(dbResult[0].Lot).To(Equal("115792089237316195423570985008687907853269984665640564039457584007913129639935"))
	})

	It("rechecks flop kick event", func() {
		blockNumber := int64(8672119)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := fetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []flop_kick.Model
		err = db.Select(&dbResult, `SELECT bid, bid_id, "end", gal, lot FROM maker.flop_kick`)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var flopKickChecked []int
		err = db.Select(&flopKickChecked, `SELECT flop_kick_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(flopKickChecked[0]).To(Equal(2))

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("0"))
		Expect(dbResult[0].BidId).To(Equal("1"))
		Expect(dbResult[0].End.Equal(time.Unix(1536726768, 0))).To(BeTrue())
		Expect(dbResult[0].Gal).To(Equal("0x9B870D55BaAEa9119dBFa71A92c5E26E79C4726d"))
		// this very large number appears to be derived from the data including: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		Expect(dbResult[0].Lot).To(Equal("115792089237316195423570985008687907853269984665640564039457584007913129639935"))
	})

	It("fetches and transforms another FlopKick event from Kovan chain", func() {
		blockNumber := int64(8955611)
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logs, err := fetcher.FetchLogs(addresses, topics, header)
		Expect(err).NotTo(HaveOccurred())

		transformer := initializer.NewTransformer(db)
		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []flop_kick.Model
		err = db.Select(&dbResult, `SELECT bid, bid_id, "end", gal, lot FROM maker.flop_kick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("10000000000000000000000"))
		Expect(dbResult[0].BidId).To(Equal("2"))
		Expect(dbResult[0].End.Equal(time.Unix(1538810564, 0))).To(BeTrue())
		Expect(dbResult[0].Gal).To(Equal("0x3728e9777B2a0a611ee0F89e00E01044ce4736d1"))
		Expect(dbResult[0].Lot).To(Equal("115792089237316195423570985008687907853269984665640564039457584007913129639935"))
	})

	It("unpacks an flop kick event log", func() {
		address := common.HexToAddress(test_data.KovanFlopperContractAddress)
		abi, err := geth.ParseAbi(test_data.KovanFlopperABI)
		Expect(err).NotTo(HaveOccurred())

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		entity := &flop_kick.Entity{}

		var eventLog = test_data.FlopKickLog

		err = contract.UnpackLog(entity, "Kick", eventLog)
		Expect(err).NotTo(HaveOccurred())

		expectedEntity := test_data.FlopKickEntity
		Expect(entity.Id).To(Equal(expectedEntity.Id))
		Expect(entity.Lot).To(Equal(expectedEntity.Lot))
		Expect(entity.Bid).To(Equal(expectedEntity.Bid))
		Expect(entity.Gal).To(Equal(expectedEntity.Gal))
		Expect(entity.End).To(Equal(expectedEntity.End))
	})
})
