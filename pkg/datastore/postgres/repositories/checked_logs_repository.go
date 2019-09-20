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

package repositories

import (
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type CheckedLogsRepository struct {
	db *postgres.DB
}

func NewCheckedLogsRepository(db *postgres.DB) CheckedLogsRepository {
	return CheckedLogsRepository{db: db}
}

// Return whether a given address + topic0 has been fetched on a previous run of vDB
func (repository CheckedLogsRepository) AlreadyWatchingLog(addresses []string, topic0 string) (bool, error) {
	for _, address := range addresses {
		var addressExists bool
		getAddressExistsErr := repository.db.Get(&addressExists, `SELECT EXISTS(SELECT 1 FROM public.watched_logs WHERE contract_address = $1)`, address)
		if getAddressExistsErr != nil {
			return false, getAddressExistsErr
		}
		if !addressExists {
			return false, nil
		}
	}
	var topicZeroExists bool
	getTopicZeroExistsErr := repository.db.Get(&topicZeroExists, `SELECT EXISTS(SELECT 1 FROM public.watched_logs WHERE topic_zero = $1)`, topic0)
	if getTopicZeroExistsErr != nil {
		return false, getTopicZeroExistsErr
	}
	return topicZeroExists, nil
}

// Persist that a given address + topic0 has is being fetched on this run of vDB
func (repository CheckedLogsRepository) MarkLogWatched(addresses []string, topic0 string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	for _, address := range addresses {
		_, insertErr := tx.Exec(`INSERT INTO public.watched_logs (contract_address, topic_zero) VALUES ($1, $2)`, address, topic0)
		if insertErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Errorf("error rolling back transaction inserting checked logs: %s", rollbackErr.Error())
			}
			return insertErr
		}
	}
	return tx.Commit()
}
