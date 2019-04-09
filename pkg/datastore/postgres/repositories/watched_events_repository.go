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
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type WatchedEventRepository struct {
	*postgres.DB
}

func (watchedEventRepository WatchedEventRepository) GetWatchedEvents(name string) ([]*core.WatchedEvent, error) {
	rows, err := watchedEventRepository.DB.Queryx(`SELECT id, name, block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data FROM watched_event_logs where name=$1`, name)
	if err != nil {
		logrus.Error("GetWatchedEvents: error getting watched events: ", err)
		return nil, err
	}
	defer rows.Close()

	lgs := make([]*core.WatchedEvent, 0)
	for rows.Next() {
		lg := new(core.WatchedEvent)
		err = rows.StructScan(lg)
		if err != nil {
			logrus.Warn("GetWatchedEvents: error scanning log: ", err)
			return nil, err
		}
		lgs = append(lgs, lg)
	}
	if err = rows.Err(); err != nil {
		logrus.Warn("GetWatchedEvents: error scanning logs: ", err)
		return nil, err
	}
	return lgs, nil
}
