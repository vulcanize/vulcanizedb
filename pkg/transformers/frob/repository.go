// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package frob

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerID int64, transactionIndex uint, model FrobModel) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}

type FrobRepository struct {
	db *postgres.DB
}

func NewFrobRepository(db *postgres.DB) FrobRepository {
	return FrobRepository{db: db}
}

func (repository FrobRepository) Create(headerID int64, transactionIndex uint, model FrobModel) error {
	_, err := repository.db.Exec(`INSERT INTO maker.frob (header_id, tx_idx, art, era, gem, ilk, ink, lad)
		VALUES($1, $2, $3::NUMERIC, $4::NUMERIC, $5::NUMERIC, $6, $7::NUMERIC, $8)`,
		headerID, transactionIndex, model.Art, model.Era, model.Gem, model.Ilk, model.Ink, model.Lad)
	return err
}

func (repository FrobRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN maker.frob on headers.id = header_id
               WHERE header_id ISNULL
               AND headers.block_number >= $1
               AND headers.block_number <= $2
               AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		repository.db.Node.ID,
	)
	return result, err
}
