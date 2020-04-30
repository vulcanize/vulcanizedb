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
	"bytes"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-ds-help"
	"github.com/multiformats/go-multihash"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"

	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

var _ = Describe("PublishAndIndexer", func() {
	var (
		db        *postgres.DB
		err       error
		repo      *btc.IPLDPublisherAndIndexer
		ipfsPgGet = `SELECT data FROM public.blocks
					WHERE key = $1`
	)
	BeforeEach(func() {
		db, err = shared.SetupDB()
		Expect(err).ToNot(HaveOccurred())
		repo = btc.NewIPLDPublisherAndIndexer(db)
	})
	AfterEach(func() {
		btc.TearDownDB(db)
	})

	Describe("Publish", func() {
		It("Published and indexes header and transaction IPLDs in a single tx", func() {
			emptyReturn, err := repo.Publish(mocks.MockConvertedPayload)
			Expect(emptyReturn).To(BeNil())
			Expect(err).ToNot(HaveOccurred())
			pgStr := `SELECT * FROM btc.header_cids
				WHERE block_number = $1`
			// check header was properly indexed
			buf := bytes.NewBuffer(make([]byte, 0, 80))
			err = mocks.MockBlock.Header.Serialize(buf)
			Expect(err).ToNot(HaveOccurred())
			headerBytes := buf.Bytes()
			c, _ := ipld.RawdataToCid(ipld.MBitcoinHeader, headerBytes, multihash.DBL_SHA2_256)
			header := new(btc.HeaderModel)
			err = db.Get(header, pgStr, mocks.MockHeaderMetaData.BlockNumber)
			Expect(err).ToNot(HaveOccurred())
			Expect(header.CID).To(Equal(c.String()))
			Expect(header.BlockNumber).To(Equal(mocks.MockHeaderMetaData.BlockNumber))
			Expect(header.Bits).To(Equal(mocks.MockHeaderMetaData.Bits))
			Expect(header.Timestamp).To(Equal(mocks.MockHeaderMetaData.Timestamp))
			Expect(header.BlockHash).To(Equal(mocks.MockHeaderMetaData.BlockHash))
			Expect(header.ParentHash).To(Equal(mocks.MockHeaderMetaData.ParentHash))
			dc, err := cid.Decode(header.CID)
			Expect(err).ToNot(HaveOccurred())
			mhKey := dshelp.CidToDsKey(dc)
			prefixedKey := blockstore.BlockPrefix.String() + mhKey.String()
			var data []byte
			err = db.Get(&data, ipfsPgGet, prefixedKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(data).To(Equal(headerBytes))

			// check that txs were properly indexed
			trxs := make([]btc.TxModel, 0)
			pgStr = `SELECT transaction_cids.id, transaction_cids.header_id, transaction_cids.index,
				transaction_cids.tx_hash, transaction_cids.cid, transaction_cids.segwit, transaction_cids.witness_hash
				FROM btc.transaction_cids INNER JOIN btc.header_cids ON (transaction_cids.header_id = header_cids.id)
				WHERE header_cids.block_number = $1`
			err = db.Select(&trxs, pgStr, mocks.MockHeaderMetaData.BlockNumber)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(trxs)).To(Equal(3))
			txData := make([][]byte, len(mocks.MockTransactions))
			txCIDs := make([]string, len(mocks.MockTransactions))
			for i, m := range mocks.MockTransactions {
				buf := bytes.NewBuffer(make([]byte, 0))
				err = m.MsgTx().Serialize(buf)
				Expect(err).ToNot(HaveOccurred())
				tx := buf.Bytes()
				txData[i] = tx
				c, _ := ipld.RawdataToCid(ipld.MBitcoinTx, tx, multihash.DBL_SHA2_256)
				txCIDs[i] = c.String()
			}
			for _, tx := range trxs {
				Expect(tx.SegWit).To(Equal(false))
				Expect(tx.HeaderID).To(Equal(header.ID))
				Expect(tx.WitnessHash).To(Equal(""))
				Expect(tx.CID).To(Equal(txCIDs[tx.Index]))
				Expect(tx.TxHash).To(Equal(mocks.MockBlock.Transactions[tx.Index].TxHash().String()))
				dc, err := cid.Decode(tx.CID)
				Expect(err).ToNot(HaveOccurred())
				mhKey := dshelp.CidToDsKey(dc)
				prefixedKey := blockstore.BlockPrefix.String() + mhKey.String()
				fmt.Printf("mhKey: %s\r\n", prefixedKey)
				var data []byte
				err = db.Get(&data, ipfsPgGet, prefixedKey)
				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(Equal(txData[tx.Index]))
			}
		})
	})
})
