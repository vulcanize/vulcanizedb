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
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var ErrValidHeaderExists = errors.New("valid header already exists")

type HeaderRepository struct {
	database *postgres.DB
}

func NewHeaderRepository(database *postgres.DB) HeaderRepository {
	return HeaderRepository{database: database}
}

func (repository HeaderRepository) CreateOrUpdateHeader(header core.Header) (int64, error) {
	hash, err := repository.getHeaderHash(header)
	if err != nil {
		if headerDoesNotExist(err) {
			return repository.insertHeader(header)
		}
		log.Error("CreateOrUpdateHeader: error getting header hash: ", err)
		return 0, err
	}
	if headerMustBeReplaced(hash, header) {
		return repository.replaceHeader(header)
	}
	return 0, ErrValidHeaderExists
}

func (repository HeaderRepository) GetHeader(blockNumber int64) (core.Header, error) {
	var header core.Header
	err := repository.database.Get(&header, `SELECT id, block_number, hash, raw, block_timestamp FROM headers WHERE block_number = $1 AND eth_node_fingerprint = $2`,
		blockNumber, repository.database.Node.ID)
	log.Error("GetHeader: error getting headers: ", err)
	return header, err
}

func (repository HeaderRepository) MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64, nodeID string) ([]int64, error) {
	numbers := make([]int64, 0)
	err := repository.database.Select(&numbers,
		`SELECT series.block_number
			FROM (SELECT generate_series($1::INT, $2::INT) AS block_number) AS series
			LEFT OUTER JOIN (SELECT block_number FROM headers
				WHERE eth_node_fingerprint = $3) AS synced
			USING (block_number)
			WHERE  synced.block_number IS NULL`,
		startingBlockNumber, endingBlockNumber, nodeID)
	if err != nil {
		log.Errorf("MissingBlockNumbers failed to get blocks between %v - %v for node %v",
			startingBlockNumber, endingBlockNumber, nodeID)
		return []int64{}, err
	}
	return numbers, nil
}

func headerMustBeReplaced(hash string, header core.Header) bool {
	return hash != header.Hash
}

func headerDoesNotExist(err error) bool {
	return err == sql.ErrNoRows
}

func (repository HeaderRepository) getHeaderHash(header core.Header) (string, error) {
	var hash string
	err := repository.database.Get(&hash, `SELECT hash FROM headers WHERE block_number = $1 AND eth_node_fingerprint = $2`,
		header.BlockNumber, repository.database.Node.ID)
	return hash, err
}

func (repository HeaderRepository) insertHeader(header core.Header) (int64, error) {
	var headerId int64
	err := repository.database.QueryRowx(
		`INSERT INTO public.headers (block_number, hash, block_timestamp, raw, eth_node_id, eth_node_fingerprint) VALUES ($1, $2, $3::NUMERIC, $4, $5, $6) RETURNING id`,
		header.BlockNumber, header.Hash, header.Timestamp, header.Raw, repository.database.NodeID, repository.database.Node.ID).Scan(&headerId)
	if err != nil {
		log.Error("insertHeader: error inserting header: ", err)
	}
	return headerId, err
}

func (repository HeaderRepository) replaceHeader(header core.Header) (int64, error) {
	_, err := repository.database.Exec(`DELETE FROM headers WHERE block_number = $1 AND eth_node_fingerprint = $2`,
		header.BlockNumber, repository.database.Node.ID)
	if err != nil {
		log.Error("replaceHeader: error deleting headers: ", err)
		return 0, err
	}
	return repository.insertHeader(header)
}
