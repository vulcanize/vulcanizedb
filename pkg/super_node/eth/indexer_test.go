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

package eth_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

var _ = Describe("Indexer", func() {
	var (
		db   *postgres.DB
		err  error
		repo *eth.CIDIndexer
	)
	BeforeEach(func() {
		db, err = shared.SetupDB()
		Expect(err).ToNot(HaveOccurred())
		repo = eth.NewCIDIndexer(db)
	})
	AfterEach(func() {
		eth.TearDownDB(db)
	})

	Describe("Index", func() {
		It("Indexes CIDs and related metadata into vulcanizedb", func() {
			err = repo.Index(mocks.MockCIDPayload)
			Expect(err).ToNot(HaveOccurred())
			pgStr := `SELECT cid, td, reward, id
				FROM eth.header_cids
				WHERE block_number = $1`
			// check header was properly indexed
			type res struct {
				CID    string
				TD     string
				Reward string
				ID     int
			}
			headers := new(res)
			err = db.QueryRowx(pgStr, 1).StructScan(headers)
			Expect(err).ToNot(HaveOccurred())
			Expect(headers.CID).To(Equal(mocks.HeaderCID.String()))
			Expect(headers.TD).To(Equal(mocks.MockBlock.Difficulty().String()))
			Expect(headers.Reward).To(Equal("5000000000000000000"))
			// check trxs were properly indexed
			trxs := make([]string, 0)
			pgStr = `SELECT transaction_cids.cid FROM eth.transaction_cids INNER JOIN eth.header_cids ON (transaction_cids.header_id = header_cids.id)
				WHERE header_cids.block_number = $1`
			err = db.Select(&trxs, pgStr, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(trxs)).To(Equal(2))
			Expect(shared.ListContainsString(trxs, mocks.Trx1CID.String())).To(BeTrue())
			Expect(shared.ListContainsString(trxs, mocks.Trx2CID.String())).To(BeTrue())
			// check receipts were properly indexed
			rcts := make([]string, 0)
			pgStr = `SELECT receipt_cids.cid FROM eth.receipt_cids, eth.transaction_cids, eth.header_cids
				WHERE receipt_cids.tx_id = transaction_cids.id 
				AND transaction_cids.header_id = header_cids.id
				AND header_cids.block_number = $1`
			err = db.Select(&rcts, pgStr, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(rcts)).To(Equal(2))
			Expect(shared.ListContainsString(rcts, mocks.Rct1CID.String())).To(BeTrue())
			Expect(shared.ListContainsString(rcts, mocks.Rct2CID.String())).To(BeTrue())
			// check that state nodes were properly indexed
			stateNodes := make([]eth.StateNodeModel, 0)
			pgStr = `SELECT state_cids.cid, state_cids.state_leaf_key, state_cids.node_type, state_cids.state_path, state_cids.header_id
				FROM eth.state_cids INNER JOIN eth.header_cids ON (state_cids.header_id = header_cids.id)
				WHERE header_cids.block_number = $1`
			err = db.Select(&stateNodes, pgStr, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(stateNodes)).To(Equal(2))
			for _, stateNode := range stateNodes {
				if stateNode.CID == mocks.State1CID.String() {
					Expect(stateNode.NodeType).To(Equal(2))
					Expect(stateNode.StateKey).To(Equal(common.BytesToHash(mocks.ContractLeafKey).Hex()))
					Expect(stateNode.Path).To(Equal([]byte{'\x06'}))
				}
				if stateNode.CID == mocks.State2CID.String() {
					Expect(stateNode.NodeType).To(Equal(2))
					Expect(stateNode.StateKey).To(Equal(common.BytesToHash(mocks.AccountLeafKey).Hex()))
					Expect(stateNode.Path).To(Equal([]byte{'\x0c'}))
				}
			}
			// check that storage nodes were properly indexed
			storageNodes := make([]eth.StorageNodeWithStateKeyModel, 0)
			pgStr = `SELECT storage_cids.cid, state_cids.state_leaf_key, storage_cids.storage_leaf_key, storage_cids.node_type, storage_cids.storage_path 
				FROM eth.storage_cids, eth.state_cids, eth.header_cids
				WHERE storage_cids.state_id = state_cids.id 
				AND state_cids.header_id = header_cids.id
				AND header_cids.block_number = $1`
			err = db.Select(&storageNodes, pgStr, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(storageNodes)).To(Equal(1))
			Expect(storageNodes[0]).To(Equal(eth.StorageNodeWithStateKeyModel{
				CID:        mocks.StorageCID.String(),
				NodeType:   2,
				StorageKey: common.BytesToHash(mocks.StorageLeafKey).Hex(),
				StateKey:   common.BytesToHash(mocks.ContractLeafKey).Hex(),
				Path:       []byte{},
			}))
		})
	})
})
