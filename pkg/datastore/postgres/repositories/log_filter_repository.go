// VulcanizeDB
// Copyright © 2018 Vulcanize

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
	"database/sql"

	"encoding/json"
	"errors"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type FilterRepository struct {
	*postgres.DB
}

func (filterRepository FilterRepository) CreateFilter(query filters.LogFilter) error {
	_, err := filterRepository.DB.Exec(
		`INSERT INTO log_filters 
        (name, from_block, to_block, address, topic0, topic1, topic2, topic3)
        VALUES ($1, NULLIF($2, -1), NULLIF($3, -1), $4, NULLIF($5, ''), NULLIF($6, ''), NULLIF($7, ''), NULLIF($8, ''))`,
		query.Name, query.FromBlock, query.ToBlock, query.Address, query.Topics[0], query.Topics[1], query.Topics[2], query.Topics[3])
	if err != nil {
		return err
	}
	return nil
}

func (filterRepository FilterRepository) GetFilter(name string) (filters.LogFilter, error) {
	lf := DBLogFilter{}
	err := filterRepository.DB.Get(&lf,
		`SELECT
                  id,
                  name,
                  from_block,
                  to_block,
                  address,
                  json_build_array(topic0, topic1, topic2, topic3) AS topics
                FROM log_filters
                WHERE name = $1`, name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return filters.LogFilter{}, datastore.ErrFilterDoesNotExist(name)
		default:
			return filters.LogFilter{}, err
		}
	}
	dbLogFilterToCoreLogFilter(lf)
	return *lf.LogFilter, nil
}

type DBTopics []*string

func (t *DBTopics) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("scan source was not []byte"))
	}
	return json.Unmarshal(asBytes, &t)
}

type DBLogFilter struct {
	ID int
	*filters.LogFilter
	Topics DBTopics
}

func dbLogFilterToCoreLogFilter(lf DBLogFilter) {
	for i, v := range lf.Topics {
		if v != nil {
			lf.LogFilter.Topics[i] = *v
		}
	}
}
