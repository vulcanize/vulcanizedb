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
	"database/sql"
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
)

type Validator struct {
	db *postgres.DB
}

func NewValidator(db *postgres.DB) *Validator {
	return &Validator{
		db: db,
	}
}

// Validate is the top-level method for the Ethereum Validator
// It validates sections of headers
// And generates state caches every 2048 blocks using these validated headers
func (v *Validator) Validate(errChan chan error) error {
	// Algo:
	/*
		1. Find the latest state cache (n)
		2. If there are none => start at block 2048 (n+2048)
		3. Validate headers between n and n+2048
			- Perform +256 validation of the header at n+2048
		  	- Must wait for n+2048 to be at least 256 blocks behind the head to do this
		4. If the section is valid and complete, find the latest state diffs for every state_path in the range
		5. Apply these diffs ontop of the cache at n
		NOTE: perform as much of this in the DB and in as few txs as possible
	*/
	panic("implement me")
}

func (v *Validator) getLatestStateCache() (int64, error) {
	pgStr := `SELECT block_number FROM eth.state_cache 
			INNER JOIN eth.header_cids ON (state_cache.header_id = header_cids.id)
			ORDER BY block_number DESC LIMIT 1`
	var latestCache int64
	err := v.db.Get(&latestCache, pgStr)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return latestCache, err
}

func (v *Validator) getHeadersAt(height int) ([]HeaderModel, error) {
	pgStr := `SELECT * FROM eth.header_cids WHERE block_number = $1 AND times_validated > 0`
	var headers []HeaderModel
	return headers, v.db.Select(&headers, pgStr, height)
}

func (v *Validator) validate(height, weight int, headers []HeaderModel) (HeaderModel, bool, error) {
	headers256, err := v.getHeadersAt(height + weight)
	if err != nil {
		if err == sql.ErrNoRows {
			return HeaderModel{}, false, fmt.Errorf("need a header at %d to validate header at %d", height+weight, height)
		}
		return HeaderModel{}, false, err
	}
	header256 := headers256[0]
	if len(headers256) > 1 {
		// We need to recursively 256 validate one of these headers
		validHeader, valid, err := v.validate(height+weight, weight, headers256)
		if !valid || err != nil {
			return HeaderModel{}, false, fmt.Errorf("unable to validate a header at %d from which to validate header at %d; err: %v", height+weight, height, err)
		}
		header256 = validHeader
	}
	validHeaderIDs := make([]int64, 0, weight+1)
	pgStr := `SELECT * FROM eth_valid_section($1, $2)`
	if err := v.db.Select(&validHeaderIDs, pgStr, header256.BlockHash, weight); err != nil {
		return HeaderModel{}, false, err
	}
	// If we don't have weight+1 headers, the section was incomplete
	if len(validHeaderIDs) != weight+1 {
		// TODO: mark the missing header heights as needing validated (set times_validated = 0)
		return HeaderModel{}, false, fmt.Errorf("section of headers between %d and %d is incomplete", height, height+weight)
	}
	// Which ever header is equal to the valid ID at the end of range is the valid header
	validID := validHeaderIDs[256]
	for _, header := range headers {
		if header.ID == validID {
			return header, true, nil
		}
	}
	return HeaderModel{}, false, nil
}
