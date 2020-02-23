// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package btc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

var _ = Describe("Indexer", func() {
	var (
		db   *postgres.DB
		err  error
		repo *btc.CIDIndexer
	)
	BeforeEach(func() {
		db, err = shared.SetupDB()
		Expect(err).ToNot(HaveOccurred())
		repo = btc.NewCIDIndexer(db)
	})
	AfterEach(func() {
		btc.TearDownDB(db)
	})

	Describe("Index", func() {
		It("Indexes CIDs and related metadata into vulcanizedb", func() {
			err = repo.Index(&mocks.MockCIDPayload)
			Expect(err).ToNot(HaveOccurred())
			pgStr := `SELECT * FROM btc.header_cids
				WHERE block_number = $1`
			// check header was properly indexed
			header := new(btc.HeaderModel)
			err = db.Get(header, pgStr, mocks.MockHeaderMetaData.BlockNumber)
			Expect(err).ToNot(HaveOccurred())
			Expect(header.CID).To(Equal(mocks.MockHeaderMetaData.CID))
			Expect(header.BlockNumber).To(Equal(mocks.MockHeaderMetaData.BlockNumber))
			Expect(header.Bits).To(Equal(mocks.MockHeaderMetaData.Bits))
			Expect(header.Timestamp).To(Equal(mocks.MockHeaderMetaData.Timestamp))
			Expect(header.BlockHash).To(Equal(mocks.MockHeaderMetaData.BlockHash))
			Expect(header.ParentHash).To(Equal(mocks.MockHeaderMetaData.ParentHash))
			// check trxs were properly indexed
			trxs := make([]btc.TxModel, 0)
			pgStr = `SELECT transaction_cids.id, transaction_cids.header_id, transaction_cids.index,
				transaction_cids.tx_hash, transaction_cids.cid, transaction_cids.segwit, transaction_cids.witness_hash
				FROM btc.transaction_cids INNER JOIN btc.header_cids ON (transaction_cids.header_id = header_cids.id)
				WHERE header_cids.block_number = $1`
			err = db.Select(&trxs, pgStr, mocks.MockHeaderMetaData.BlockNumber)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(trxs)).To(Equal(3))
			for _, tx := range trxs {
				Expect(tx.SegWit).To(Equal(false))
				Expect(tx.HeaderID).To(Equal(header.ID))
				Expect(tx.WitnessHash).To(Equal(""))
				switch tx.Index {
				case 0:
					Expect(tx.CID).To(Equal("mockTrxCID1"))
					Expect(tx.TxHash).To(Equal(mocks.MockBlock.Transactions[0].TxHash().String()))
				case 1:
					Expect(tx.CID).To(Equal("mockTrxCID2"))
					Expect(tx.TxHash).To(Equal(mocks.MockBlock.Transactions[1].TxHash().String()))
				case 2:
					Expect(tx.CID).To(Equal("mockTrxCID3"))
					Expect(tx.TxHash).To(Equal(mocks.MockBlock.Transactions[2].TxHash().String()))
				}
			}
		})
	})
})
