// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event_triggered

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type GenericEventDatastore interface {
	CreateBurn(model *BurnModel, vulcanizeLogId int64) error
	CreateMint(model *MintModel, vulcanizeLogId int64) error
}

type GenericEventRepository struct {
	*postgres.DB
}

func (repository GenericEventRepository) CreateBurn(burnModel *BurnModel, vulcanizeLogId int64) error {
	_, err := repository.DB.Exec(

		`INSERT INTO token_burns (vulcanize_log_id, token_name, token_address, burner, tokens, block, tx)
               VALUES ($1, $2, $3, $4, $5, $6, $7)
                ON CONFLICT (vulcanize_log_id) DO NOTHING`,
		vulcanizeLogId, burnModel.TokenName, burnModel.TokenAddress, burnModel.Burner, burnModel.Tokens, burnModel.Block, burnModel.TxHash)

	return err
}

func (repository GenericEventRepository) CreateMint(mintModel *MintModel, vulcanizeLogId int64) error {
	_, err := repository.DB.Exec(

		`INSERT INTO token_mints (vulcanize_log_id, token_name, token_address, minter, mintee, tokens, block, tx)
               VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                ON CONFLICT (vulcanize_log_id) DO NOTHING`,
		vulcanizeLogId, mintModel.TokenName, mintModel.TokenAddress, mintModel.Minter, mintModel.Mintee, mintModel.Tokens, mintModel.Block, mintModel.TxHash)

	return err
}
