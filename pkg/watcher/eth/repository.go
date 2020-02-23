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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node"
	"github.com/vulcanize/vulcanizedb/pkg/watcher/shared"
)

var (
	vacuumThreshold int64 = 5000 // dont know how to decided what this should be set to
)

// Repository is the underlying struct for satisfying the shared.Repository interface for eth
type Repository struct {
	db               *postgres.DB
	triggerFunctions [][2]string
	deleteCalls      int64
}

// NewRepository returns a new eth.Repository that satisfies the shared.Repository interface
func NewRepository(db *postgres.DB, triggerFunctions [][2]string) shared.Repository {
	return &Repository{
		db:               db,
		triggerFunctions: triggerFunctions,
		deleteCalls:      0,
	}
}

// LoadTriggers is used to initialize Postgres trigger function
// this needs to be called after the wasm functions these triggers invoke have been instantiated
func (r *Repository) LoadTriggers() error {
	panic("implement me")
}

// QueueData puts super node payload data into the db queue
func (r *Repository) QueueData(payload super_node.SubscriptionPayload) error {
	pgStr := `INSERT INTO eth.queued_data (data, height) VALUES ($1, $2)
			ON CONFLICT (height) DO UPDATE SET (data) VALUES ($1)`
	_, err := r.db.Exec(pgStr, payload.Data, payload.Height)
	return err
}

// GetQueueData grabs a chunk super node payload data from the queue table so that it can
// be forwarded to the ready table
// this is used to make sure we enter data into the ready table in sequential order
// even if we receive data out-of-order
// it returns the new index
// delete the data it retrieves so as to clear the queue
func (r *Repository) GetQueueData(height int64) (super_node.SubscriptionPayload, int64, error) {
	r.deleteCalls++
	pgStr := `DELETE FROM eth.queued_data
			WHERE height = $1
			RETURNING *`
	var res shared.QueuedData
	if err := r.db.Get(&res, pgStr, height); err != nil {
		return super_node.SubscriptionPayload{}, height, err
	}
	payload := super_node.SubscriptionPayload{
		Data:   res.Data,
		Height: res.Height,
		Flag:   super_node.EmptyFlag,
	}
	height++
	// Periodically clean up space in the queue table
	if r.deleteCalls >= vacuumThreshold {
		_, err := r.db.Exec(`VACUUM ANALYZE eth.queued_data`)
		if err != nil {
			logrus.Error(err)
		}
		r.deleteCalls = 0
	}
	return payload, height, nil
}

// ReadyData puts super node payload data in the tables ready for processing by trigger functions
func (r *Repository) ReadyData(payload super_node.SubscriptionPayload) error {
	panic("implement me")
}

func (r *Repository) readyHeader(header *types.Header) error {
	panic("implement me")
}

func (r *Repository) readyUncle(uncle *types.Header) error {
	panic("implement me")
}

func (r *Repository) readyTxs(transactions types.Transactions) error {
	panic("implement me")
}

func (r *Repository) readyRcts(receipts types.Receipts) error {
	panic("implement me")
}

func (r *Repository) readyState(stateNodes map[common.Address][]byte) error {
	panic("implement me")
}

func (r *Repository) readyStorage(storageNodes map[common.Address]map[common.Address][]byte) error {
	panic("implement me")
}
