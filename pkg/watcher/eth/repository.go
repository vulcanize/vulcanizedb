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
	"io/ioutil"

	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/watcher/shared"
)

var (
	vacuumThreshold int64 = 5000
)

// Repository is the underlying struct for satisfying the shared.Repository interface for eth
type Repository struct {
	cidIndexer       *eth.CIDIndexer
	converter        *WatcherConverter
	db               *postgres.DB
	triggerFunctions []string
	deleteCalls      int64
}

// NewRepository returns a new eth.Repository that satisfies the shared.Repository interface
func NewRepository(db *postgres.DB, triggerFunctions []string) shared.Repository {
	return &Repository{
		cidIndexer:       eth.NewCIDIndexer(db),
		converter:        NewWatcherConverter(params.MainnetChainConfig),
		db:               db,
		triggerFunctions: triggerFunctions,
		deleteCalls:      0,
	}
}

// LoadTriggers is used to initialize Postgres trigger function
// this needs to be called after the wasm functions these triggers invoke have been instantiated in Postgres
func (r *Repository) LoadTriggers() error {
	// TODO: enable loading of triggers from IPFS
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	for _, funcPath := range r.triggerFunctions {
		sqlFile, err := ioutil.ReadFile(funcPath)
		if err != nil {
			return err
		}
		sqlString := string(sqlFile)
		if _, err := tx.Exec(sqlString); err != nil {
			return err
		}

	}
	return tx.Commit()
}

// QueueData puts super node payload data into the db queue
func (r *Repository) QueueData(payload super_node.SubscriptionPayload) error {
	pgStr := `INSERT INTO eth.queued_data (data, height) VALUES ($1, $2)
			ON CONFLICT (height) DO UPDATE SET (data) VALUES ($1)`
	_, err := r.db.Exec(pgStr, payload.Data, payload.Height)
	return err
}

// GetQueueData grabs payload data from the queue table so that it can be readied
// Used ensure we enter data into the tables that triggers act on in sequential order, even if we receive data out-of-order
// Returns the queued data, the new index, and err
// Deletes from the queue the data it retrieves
// Periodically vacuum's the table to free up space from the deleted rows
func (r *Repository) GetQueueData(height int64) (super_node.SubscriptionPayload, int64, error) {
	pgStr := `DELETE FROM eth.queued_data
			WHERE height = $1
			RETURNING *`
	var res shared.QueuedData
	if err := r.db.Get(&res, pgStr, height); err != nil {
		return super_node.SubscriptionPayload{}, height, err
	}
	// If the delete get query succeeded, increment deleteCalls and height and prep payload to return
	r.deleteCalls++
	height++
	payload := super_node.SubscriptionPayload{
		Data:   res.Data,
		Height: res.Height,
		Flag:   super_node.EmptyFlag,
	}
	// Periodically clean up space in the queued data table
	if r.deleteCalls >= vacuumThreshold {
		_, err := r.db.Exec(`VACUUM ANALYZE eth.queued_data`)
		if err != nil {
			logrus.Error(err)
		}
		r.deleteCalls = 0
	}
	return payload, height, nil
}

// ReadyData puts data in the tables ready for processing by trigger functions
func (r *Repository) ReadyData(payload super_node.SubscriptionPayload) error {
	var ethIPLDs eth.IPLDs
	if err := rlp.DecodeBytes(payload.Data, &ethIPLDs); err != nil {
		return err
	}
	if err := r.readyIPLDs(ethIPLDs); err != nil {
		return err
	}
	cids, err := r.converter.Convert(ethIPLDs)
	if err != nil {
		return err
	}
	// Use indexer to persist all of the cid meta data
	// trigger functions will act on these tables
	return r.cidIndexer.Index(cids)
}

// readyIPLDs adds IPLDs directly to the Postgres `blocks` table, rather than going through an IPFS node
func (r *Repository) readyIPLDs(ethIPLDs eth.IPLDs) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	pgStr := `INSERT INTO blocks (key, data) VALUES ($1, $2) 
			ON CONFLICT (key) DO UPDATE SET (data) = ($2)`
	if _, err := tx.Exec(pgStr, ethIPLDs.Header.CID, ethIPLDs.Header.Data); err != nil {
		if err := tx.Rollback(); err != nil {
			logrus.Error(err)
		}
		return err
	}
	for _, uncle := range ethIPLDs.Uncles {
		if _, err := tx.Exec(pgStr, uncle.CID, uncle.Data); err != nil {
			if err := tx.Rollback(); err != nil {
				logrus.Error(err)
			}
			return err
		}
	}
	for _, trx := range ethIPLDs.Transactions {
		if _, err := tx.Exec(pgStr, trx.CID, trx.Data); err != nil {
			if err := tx.Rollback(); err != nil {
				logrus.Error(err)
			}
			return err
		}
	}
	for _, rct := range ethIPLDs.Receipts {
		if _, err := tx.Exec(pgStr, rct.CID, rct.Data); err != nil {
			if err := tx.Rollback(); err != nil {
				logrus.Error(err)
			}
			return err
		}
	}
	for _, state := range ethIPLDs.StateNodes {
		if _, err := tx.Exec(pgStr, state.IPLD.CID, state.IPLD.Data); err != nil {
			if err := tx.Rollback(); err != nil {
				logrus.Error(err)
			}
			return err
		}
	}
	for _, storage := range ethIPLDs.StorageNodes {
		if _, err := tx.Exec(pgStr, storage.IPLD.CID, storage.IPLD.Data); err != nil {
			if err := tx.Rollback(); err != nil {
				logrus.Error(err)
			}
			return err
		}
	}
	return nil
}
