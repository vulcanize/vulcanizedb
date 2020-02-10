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

package btc

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/postgres"
)

// TearDownDB is used to tear down the super node dbs after tests
func TearDownDB(db *postgres.DB) {
	tx, err := db.Beginx()
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM btc.header_cids`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM btc.transaction_cids`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM btc.tx_inputs`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM btc.tx_outputs`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM blocks`)
	Expect(err).NotTo(HaveOccurred())

	err = tx.Commit()
	Expect(err).NotTo(HaveOccurred())
}
