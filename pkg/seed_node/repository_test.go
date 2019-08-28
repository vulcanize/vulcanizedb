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

package seed_node_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/seed_node"
)

var (
	db   *postgres.DB
	err  error
	repo seed_node.CIDRepository
)

var _ = Describe("Repository", func() {
	BeforeEach(func() {
		db, err = seed_node.SetupDB()
		Expect(err).ToNot(HaveOccurred())
		repo = seed_node.NewCIDRepository(db)
	})
	AfterEach(func() {
		seed_node.TearDownDB(db)
	})

	Describe("Index", func() {
		It("Indexes CIDs and related metadata into vulcanizedb", func() {
			err = repo.Index(mocks.MockCIDPayload)
			Expect(err).ToNot(HaveOccurred())
			pgStr := `SELECT cid FROM header_cids
				WHERE block_number = $1 AND final IS TRUE`
			// check header was properly indexed
			headers := make([]string, 0)
			err = db.Select(&headers, pgStr, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(headers)).To(Equal(1))
			Expect(headers[0]).To(Equal("mockHeaderCID"))
			// check trxs were properly indexed
			trxs := make([]string, 0)
			pgStr = `SELECT transaction_cids.cid FROM transaction_cids INNER JOIN header_cids ON (transaction_cids.header_id = header_cids.id)
				WHERE header_cids.block_number = $1`
			err = db.Select(&trxs, pgStr, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(trxs)).To(Equal(2))
			Expect(seed_node.ListContainsString(trxs, "mockTrxCID1")).To(BeTrue())
			Expect(seed_node.ListContainsString(trxs, "mockTrxCID2")).To(BeTrue())
			// check receipts were properly indexed
			rcts := make([]string, 0)
			pgStr = `SELECT receipt_cids.cid FROM receipt_cids, transaction_cids, header_cids
				WHERE receipt_cids.tx_id = transaction_cids.id 
				AND transaction_cids.header_id = header_cids.id
				AND header_cids.block_number = $1`
			err = db.Select(&rcts, pgStr, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(rcts)).To(Equal(2))
			Expect(seed_node.ListContainsString(rcts, "mockRctCID1")).To(BeTrue())
			Expect(seed_node.ListContainsString(rcts, "mockRctCID2")).To(BeTrue())
			// check that state nodes were properly indexed
			stateNodes := make([]ipfs.StateNodeCID, 0)
			pgStr = `SELECT state_cids.cid, state_cids.state_key, state_cids.leaf FROM state_cids INNER JOIN header_cids ON (state_cids.header_id = header_cids.id)
				WHERE header_cids.block_number = $1`
			err = db.Select(&stateNodes, pgStr, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(stateNodes)).To(Equal(2))
			for _, stateNode := range stateNodes {
				if stateNode.CID == "mockStateCID1" {
					Expect(stateNode.Leaf).To(Equal(true))
					Expect(stateNode.Key).To(Equal(mocks.ContractLeafKey.Hex()))
				}
				if stateNode.CID == "mockStateCID2" {
					Expect(stateNode.Leaf).To(Equal(true))
					Expect(stateNode.Key).To(Equal(mocks.AnotherContractLeafKey.Hex()))
				}
			}
			// check that storage nodes were properly indexed
			storageNodes := make([]ipfs.StorageNodeCID, 0)
			pgStr = `SELECT storage_cids.cid, state_cids.state_key, storage_cids.storage_key, storage_cids.leaf FROM storage_cids, state_cids, header_cids
				WHERE storage_cids.state_id = state_cids.id 
				AND state_cids.header_id = header_cids.id
				AND header_cids.block_number = $1`
			err = db.Select(&storageNodes, pgStr, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(storageNodes)).To(Equal(1))
			Expect(storageNodes[0]).To(Equal(ipfs.StorageNodeCID{
				CID:      "mockStorageCID",
				Leaf:     true,
				Key:      "0x0000000000000000000000000000000000000000000000000000000000000001",
				StateKey: mocks.ContractLeafKey.Hex(),
			}))
		})
	})
})
