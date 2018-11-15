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
