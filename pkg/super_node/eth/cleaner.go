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
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// Cleaner satisfies the shared.Cleaner interface fo ethereum
type Cleaner struct {
	db *postgres.DB
}

// NewCleaner returns a new Cleaner struct that satisfies the shared.Cleaner interface
func NewCleaner(db *postgres.DB) *Cleaner {
	return &Cleaner{
		db: db,
	}
}

// ResetValidation resets the validation level to 0 to enable revalidation
func (c *Cleaner) ResetValidation(rngs [][2]uint64) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}
	for _, rng := range rngs {
		logrus.Infof("eth db cleaner resetting validation level to 0 for block range %d to %d", rng[0], rng[1])
		pgStr := `UPDATE eth.header_cids
				SET times_validated = 0
				WHERE block_number BETWEEN $1 AND $2`
		if _, err := tx.Exec(pgStr, rng[0], rng[1]); err != nil {
			shared.Rollback(tx)
			return err
		}
	}
	return tx.Commit()
}

// Clean removes the specified data from the db within the provided block range
func (c *Cleaner) Clean(rngs [][2]uint64, t shared.DataType) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}
	for _, rng := range rngs {
		logrus.Infof("eth db cleaner cleaning up block range %d to %d", rng[0], rng[1])
		if err := c.clean(tx, rng, t); err != nil {
			shared.Rollback(tx)
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	logrus.Infof("eth db cleaner vacuum analyzing cleaned tables to free up space from deleted rows")
	return c.vacuumAnalyze(t)
}

func (c *Cleaner) clean(tx *sqlx.Tx, rng [2]uint64, t shared.DataType) error {
	switch t {
	case shared.Full, shared.Headers:
		return c.cleanFull(tx, rng)
	case shared.Uncles:
		if err := c.cleanUncleIPLDs(tx, rng); err != nil {
			return err
		}
		return c.cleanUncleMetaData(tx, rng)
	case shared.Transactions:
		if err := c.cleanReceiptIPLDs(tx, rng); err != nil {
			return err
		}
		if err := c.cleanTransactionIPLDs(tx, rng); err != nil {
			return err
		}
		return c.cleanTransactionMetaData(tx, rng)
	case shared.Receipts:
		if err := c.cleanReceiptIPLDs(tx, rng); err != nil {
			return err
		}
		return c.cleanReceiptMetaData(tx, rng)
	case shared.State:
		if err := c.cleanStorageIPLDs(tx, rng); err != nil {
			return err
		}
		if err := c.cleanStateIPLDs(tx, rng); err != nil {
			return err
		}
		return c.cleanStateMetaData(tx, rng)
	case shared.Storage:
		if err := c.cleanStorageIPLDs(tx, rng); err != nil {
			return err
		}
		return c.cleanStorageMetaData(tx, rng)
	default:
		return fmt.Errorf("eth cleaner unrecognized type: %s", t.String())
	}
}

func (c *Cleaner) vacuumAnalyze(t shared.DataType) error {
	switch t {
	case shared.Full, shared.Headers:
		return c.vacuumFull()
	case shared.Uncles:
		if err := c.vacuumUncles(); err != nil {
			return err
		}
	case shared.Transactions:
		if err := c.vacuumTxs(); err != nil {
			return err
		}
		if err := c.vacuumRcts(); err != nil {
			return err
		}
	case shared.Receipts:
		if err := c.vacuumRcts(); err != nil {
			return err
		}
	case shared.State:
		if err := c.vacuumState(); err != nil {
			return err
		}
		if err := c.vacuumAccounts(); err != nil {
			return err
		}
		if err := c.vacuumStorage(); err != nil {
			return err
		}
	case shared.Storage:
		if err := c.vacuumStorage(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("eth cleaner unrecognized type: %s", t.String())
	}
	return c.vacuumIPLDs()
}

func (c *Cleaner) vacuumFull() error {
	if err := c.vacuumHeaders(); err != nil {
		return err
	}
	if err := c.vacuumUncles(); err != nil {
		return err
	}
	if err := c.vacuumTxs(); err != nil {
		return err
	}
	if err := c.vacuumRcts(); err != nil {
		return err
	}
	if err := c.vacuumState(); err != nil {
		return err
	}
	if err := c.vacuumAccounts(); err != nil {
		return err
	}
	return c.vacuumStorage()
}

func (c *Cleaner) vacuumHeaders() error {
	_, err := c.db.Exec(`VACUUM ANALYZE eth.header_cids`)
	return err
}

func (c *Cleaner) vacuumUncles() error {
	_, err := c.db.Exec(`VACUUM ANALYZE eth.uncle_cids`)
	return err
}

func (c *Cleaner) vacuumTxs() error {
	_, err := c.db.Exec(`VACUUM ANALYZE eth.transaction_cids`)
	return err
}

func (c *Cleaner) vacuumRcts() error {
	_, err := c.db.Exec(`VACUUM ANALYZE eth.receipt_cids`)
	return err
}

func (c *Cleaner) vacuumState() error {
	_, err := c.db.Exec(`VACUUM ANALYZE eth.state_cids`)
	return err
}

func (c *Cleaner) vacuumAccounts() error {
	_, err := c.db.Exec(`VACUUM ANALYZE eth.state_accounts`)
	return err
}

func (c *Cleaner) vacuumStorage() error {
	_, err := c.db.Exec(`VACUUM ANALYZE eth.storage_cids`)
	return err
}

func (c *Cleaner) vacuumIPLDs() error {
	_, err := c.db.Exec(`VACUUM ANALYZE public.blocks`)
	return err
}

func (c *Cleaner) cleanFull(tx *sqlx.Tx, rng [2]uint64) error {
	if err := c.cleanStorageIPLDs(tx, rng); err != nil {
		return err
	}
	if err := c.cleanStateIPLDs(tx, rng); err != nil {
		return err
	}
	if err := c.cleanReceiptIPLDs(tx, rng); err != nil {
		return err
	}
	if err := c.cleanTransactionIPLDs(tx, rng); err != nil {
		return err
	}
	if err := c.cleanUncleIPLDs(tx, rng); err != nil {
		return err
	}
	if err := c.cleanHeaderIPLDs(tx, rng); err != nil {
		return err
	}
	return c.cleanHeaderMetaData(tx, rng)
}

func (c *Cleaner) cleanStorageIPLDs(tx *sqlx.Tx, rng [2]uint64) error {
	pgStr := `DELETE FROM public.blocks A
			USING eth.storage_cids B, eth.state_cids C, eth.header_cids D
			WHERE A.key = B.cid
			AND B.state_id = C.id
			AND C.header_id = D.id
			AND D.block_number BETWEEN $1 AND $2`
	_, err := tx.Exec(pgStr, rng[0], rng[1])
	return err
}

func (c *Cleaner) cleanStorageMetaData(tx *sqlx.Tx, rng [2]uint64) error {
	pgStr := `DELETE FROM eth.storage_cids A
			USING eth.state_cids B, eth.header_cids C
			WHERE A.state_id = B.id
			AND B.header_id = C.id
			AND C.block_number BETWEEN $1 AND $2`
	_, err := tx.Exec(pgStr, rng[0], rng[1])
	return err
}

func (c *Cleaner) cleanStateIPLDs(tx *sqlx.Tx, rng [2]uint64) error {
	pgStr := `DELETE FROM public.blocks A
			USING eth.state_cids B, eth.header_cids C
			WHERE A.key = B.cid
			AND B.header_id = C.id
			AND C.block_number BETWEEN $1 AND $2`
	_, err := tx.Exec(pgStr, rng[0], rng[1])
	return err
}

func (c *Cleaner) cleanStateMetaData(tx *sqlx.Tx, rng [2]uint64) error {
	pgStr := `DELETE FROM eth.state_cids A
			USING eth.header_cids B
			WHERE A.header_id = B.id
			AND B.block_number BETWEEN $1 AND $2`
	_, err := tx.Exec(pgStr, rng[0], rng[1])
	return err
}

func (c *Cleaner) cleanReceiptIPLDs(tx *sqlx.Tx, rng [2]uint64) error {
	pgStr := `DELETE FROM public.blocks A
			USING eth.receipt_cids B, eth.transaction_cids C, eth.header_cids D
			WHERE A.key = B.cid
			AND B.tx_id = C.id
			AND C.header_id = D.id
			AND D.block_number BETWEEN $1 AND $2`
	_, err := tx.Exec(pgStr, rng[0], rng[1])
	return err
}

func (c *Cleaner) cleanReceiptMetaData(tx *sqlx.Tx, rng [2]uint64) error {
	pgStr := `DELETE FROM eth.receipt_cids A
			USING eth.transaction_cids B, eth.header_cids C
			WHERE A.tx_id = B.id
			AND B.header_id = C.id
			AND C.block_number BETWEEN $1 AND $2`
	_, err := tx.Exec(pgStr, rng[0], rng[1])
	return err
}

func (c *Cleaner) cleanTransactionIPLDs(tx *sqlx.Tx, rng [2]uint64) error {
	pgStr := `DELETE FROM public.blocks A
			USING eth.transaction_cids B, eth.header_cids C
			WHERE A.key = B.cid
			AND B.header_id = C.id
			AND C.block_number BETWEEN $1 AND $2`
	_, err := tx.Exec(pgStr, rng[0], rng[1])
	return err
}

func (c *Cleaner) cleanTransactionMetaData(tx *sqlx.Tx, rng [2]uint64) error {
	pgStr := `DELETE FROM eth.transaction_cids A
			USING eth.header_cids B
			WHERE A.header_id = B.id
			AND B.block_number BETWEEN $1 AND $2`
	_, err := tx.Exec(pgStr, rng[0], rng[1])
	return err
}

func (c *Cleaner) cleanUncleIPLDs(tx *sqlx.Tx, rng [2]uint64) error {
	pgStr := `DELETE FROM public.blocks A
			USING eth.uncle_cids B, eth.header_cids C
			WHERE A.key = B.cid
			AND B.header_id = C.id
			AND C.block_number BETWEEN $1 AND $2`
	_, err := tx.Exec(pgStr, rng[0], rng[1])
	return err
}

func (c *Cleaner) cleanUncleMetaData(tx *sqlx.Tx, rng [2]uint64) error {
	pgStr := `DELETE FROM eth.uncle_cids A
			USING eth.header_cids B
			WHERE A.header_id = B.id
			AND B.block_number BETWEEN $1 AND $2`
	_, err := tx.Exec(pgStr, rng[0], rng[1])
	return err
}

func (c *Cleaner) cleanHeaderIPLDs(tx *sqlx.Tx, rng [2]uint64) error {
	pgStr := `DELETE FROM public.blocks A
			USING eth.header_cids B
			WHERE A.key = B.cid
			AND B.block_number BETWEEN $1 AND $2`
	_, err := tx.Exec(pgStr, rng[0], rng[1])
	return err
}

func (c *Cleaner) cleanHeaderMetaData(tx *sqlx.Tx, rng [2]uint64) error {
	pgStr := `DELETE FROM eth.header_cids
			WHERE block_number BETWEEN $1 AND $2`
	_, err := tx.Exec(pgStr, rng[0], rng[1])
	return err
}
