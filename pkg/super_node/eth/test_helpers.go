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

package eth

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/eth/datastore/postgres"
)

// TearDownDB is used to tear down the super node dbs after tests
func TearDownDB(db *postgres.DB) {
	tx, err := db.Beginx()
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM eth.header_cids`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM eth.transaction_cids`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM eth.receipt_cids`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM eth.state_cids`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM eth.storage_cids`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM blocks`)
	Expect(err).NotTo(HaveOccurred())

	err = tx.Commit()
	Expect(err).NotTo(HaveOccurred())
}

// TxModelsContainsCID used to check if a list of TxModels contains a specific cid string
func TxModelsContainsCID(txs []TxModel, cid string) bool {
	for _, tx := range txs {
		if tx.CID == cid {
			return true
		}
	}
	return false
}

// ListContainsBytes used to check if a list of byte arrays contains a particular byte array
func ReceiptModelsContainsCID(rcts []ReceiptModel, cid string) bool {
	for _, rct := range rcts {
		if rct.CID == cid {
			return true
		}
	}
	return false
}
