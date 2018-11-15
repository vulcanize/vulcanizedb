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

type ERC20EventDatastore interface {
	CreateTransfer(model *TransferModel, vulcanizeLogId int64) error
	CreateApproval(model *ApprovalModel, vulcanizeLogId int64) error
}

type ERC20EventRepository struct {
	*postgres.DB
}

func (repository ERC20EventRepository) CreateTransfer(transferModel *TransferModel, vulcanizeLogId int64) error {
	_, err := repository.DB.Exec(

		`INSERT INTO token_transfers (vulcanize_log_id, token_name, token_address, to_address, from_address, tokens, block, tx)
               VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                ON CONFLICT (vulcanize_log_id) DO NOTHING`,
		vulcanizeLogId, transferModel.TokenName, transferModel.TokenAddress, transferModel.To, transferModel.From, transferModel.Tokens, transferModel.Block, transferModel.TxHash)

	return err
}

func (repository ERC20EventRepository) CreateApproval(approvalModel *ApprovalModel, vulcanizeLogId int64) error {
	_, err := repository.DB.Exec(

		`INSERT INTO token_approvals (vulcanize_log_id, token_name, token_address, owner, spender, tokens, block, tx)
               VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                ON CONFLICT (vulcanize_log_id) DO NOTHING`,
		vulcanizeLogId, approvalModel.TokenName, approvalModel.TokenAddress, approvalModel.Owner, approvalModel.Spender, approvalModel.Tokens, approvalModel.Block, approvalModel.TxHash)

	return err
}
