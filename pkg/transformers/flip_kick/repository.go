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

package flip_kick

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerId int64, flipKick FlipKickModel) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}

type FlipKickRepository struct {
	DB *postgres.DB
}

func NewFlipKickRepository(db *postgres.DB) FlipKickRepository {
	return FlipKickRepository{DB: db}
}
func (fkr FlipKickRepository) Create(headerId int64, flipKick FlipKickModel) error {
	_, err := fkr.DB.Exec(
		`INSERT into maker.flip_kick (header_id, id, lot, bid, gal, "end", urn, tab, raw_log)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		headerId, flipKick.Id, flipKick.Lot, flipKick.Bid, flipKick.Gal, flipKick.End, flipKick.Urn, flipKick.Tab, flipKick.Raw,
	)

	if err != nil {
		return err
	}

	return nil
}

func (fkr FlipKickRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := fkr.DB.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN maker.flip_kick on headers.id = header_id
               WHERE header_id ISNULL
               AND headers.block_number >= $1
               AND headers.block_number <= $2
               AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		fkr.DB.Node.ID,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return result, err
	}

	return result, nil
}
