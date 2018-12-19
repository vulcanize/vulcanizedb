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

package repo

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"log"
)

type DripFileRepoRepository struct {
	db *postgres.DB
}

func (repository DripFileRepoRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Begin()
	if dBaseErr != nil {
		return dBaseErr
	}

	for _, model := range models {
		repo, ok := model.(DripFileRepoModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Println("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, DripFileRepoModel{})
		}

		_, execErr := tx.Exec(
			`INSERT into maker.drip_file_repo (header_id, what, data, log_idx, tx_idx, raw_log)
        	VALUES($1, $2, $3::NUMERIC, $4, $5, $6)`,
			headerID, repo.What, repo.Data, repo.LogIndex, repo.TransactionIndex, repo.Raw,
		)

		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Println("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.DripFileRepoChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Println("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}

	return tx.Commit()
}

func (repository DripFileRepoRepository) MarkHeaderChecked(headerID int64) error {
	return shared.MarkHeaderChecked(headerID, repository.db, constants.DripFileRepoChecked)
}

func (repository *DripFileRepoRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
