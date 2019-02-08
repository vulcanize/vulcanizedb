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
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flap_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

var _ = Describe("FlapKick Transformer", func() {
	var (
		db         *postgres.DB
		blockChain core.BlockChain
	)

	BeforeEach(func() {
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)
	})

	It("fetches and transforms a FlapKick event from Kovan chain", func() {
		blockNumber := int64(9002933)
		config := shared.TransformerConfig{
			TransformerName:     constants.FlapKickLabel,
			ContractAddresses:   []string{test_data.KovanFlapperContractAddress},
			ContractAbi:         test_data.KovanFlapperABI,
			Topic:               test_data.KovanFlapKickSignature,
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		transformer := factories.Transformer{
			Config:     config,
			Converter:  &flap_kick.FlapKickConverter{},
			Repository: &flap_kick.FlapKickRepository{},
		}.NewTransformer(db)

		fetcher := shared.NewFetcher(blockChain)
		logs, err := fetcher.FetchLogs(
			shared.HexStringsToAddresses(config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []flap_kick.FlapKickModel
		err = db.Select(&dbResult, `SELECT bid, bid_id, "end", gal, lot FROM maker.flap_kick`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult[0].Bid).To(Equal("0"))
		Expect(dbResult[0].BidId).To(Equal("1"))
		Expect(dbResult[0].End.Equal(time.Unix(1539163860, 0))).To(BeTrue())
		Expect(dbResult[0].Gal).To(Equal("0x0000d8b4147eDa80Fec7122AE16DA2479Cbd7ffB"))
		Expect(dbResult[0].Lot).To(Equal("1000000000000000000"))
	})

	It("rechecks flap kick transformer", func() {
		blockNumber := int64(9002933)
		config := shared.TransformerConfig{
			TransformerName:     constants.FlapKickLabel,
			ContractAddresses:   []string{test_data.KovanFlapperContractAddress},
			ContractAbi:         test_data.KovanFlapperABI,
			Topic:               test_data.KovanFlapKickSignature,
			StartingBlockNumber: blockNumber,
			EndingBlockNumber:   blockNumber,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		transformer := factories.Transformer{
			Config:     config,
			Converter:  &flap_kick.FlapKickConverter{},
			Repository: &flap_kick.FlapKickRepository{},
		}.NewTransformer(db)

		fetcher := shared.NewFetcher(blockChain)
		logs, err := fetcher.FetchLogs(
			shared.HexStringsToAddresses(config.ContractAddresses),
			[]common.Hash{common.HexToHash(config.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderMissing)
		Expect(err).NotTo(HaveOccurred())

		err = transformer.Execute(logs, header, constants.HeaderRecheck)
		Expect(err).NotTo(HaveOccurred())

		var headerID int64
		err = db.Get(&headerID, `SELECT id FROM public.headers WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())

		var flapkickChecked []int
		err = db.Select(&flapkickChecked, `SELECT flap_kick_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
		Expect(err).NotTo(HaveOccurred())

		Expect(flapkickChecked[0]).To(Equal(2))
	})
})
