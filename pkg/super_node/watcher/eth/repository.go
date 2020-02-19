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
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/watcher/shared"
)

// Repository is the underlying struct for satisfying the shared.Repository interface for eth
type Repository struct {
	db               *postgres.DB
	triggerFunctions [][2]string
}

// NewRepository returns a new eth.Repository that satisfies the shared.Repository interface
func NewRepository(db *postgres.DB, triggerFunctions [][2]string) shared.Repository {
	return &Repository{
		db:               db,
		triggerFunctions: triggerFunctions,
	}
}

// LoadTriggers is used to initialize Postgres trigger function
// this needs to be called after the wasm functions these triggers invoke have been instantiated
func (r *Repository) LoadTriggers() error {
	panic("implement me")
}

// QueueData puts super node payload data into the db queue
func (r *Repository) QueueData(payload super_node.SubscriptionPayload) error {
	panic("implement me")
}

// GetQueueData grabs super node payload data from the db queue
func (r *Repository) GetQueueData(height int64, hash string) (super_node.SubscriptionPayload, error) {
	panic("implement me")
}

// ReadyData puts super node payload data in the tables ready for processing by trigger functions
func (r *Repository) ReadyData(payload super_node.SubscriptionPayload) error {
	panic("implement me")
}
