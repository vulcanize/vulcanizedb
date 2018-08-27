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

type Datastore interface {
	CreateTransfer(model TransferModel, vulcanizeLogId int64) error
	CreateApproval(model ApprovalModel, vulcanizeLogId int64) error
}

type Repository struct {
	*postgres.DB
}

func (repository Repository) CreateTransfer(transferModel TransferModel, vulcanizeLogId int64) error {
	_, err := repository.DB.Exec(

		`INSERT INTO TRANSFERS (vulcanize_log_id, token_address, to, from, tokens, block, tx)
               VALUES ($1, $2, $3, $4, $5, $6, $7)
                ON CONFLICT (vulcanize_log_id) DO NOTHING`,
		vulcanizeLogId, transferModel.TokenAddress, transferModel.To, transferModel.From, transferModel.Tokens, transferModel.Block, transferModel.TxHash)

	if err != nil {
		return err
	}

	return nil
}

func (repository Repository) CreateApproval(approvalModel ApprovalModel, vulcanizeLogId int64) error {
	_, err := repository.DB.Exec(

		`INSERT INTO APPROVALS (vulcanize_log_id, token_address, token_owner, token_spender, tokens, block, tx)
               VALUES ($1, $2, $3, $4, $5, $6, $7)
                ON CONFLICT (vulcanize_log_id) DO NOTHING`,
		vulcanizeLogId, approvalModel.TokenAddress, approvalModel.Owner, approvalModel.Spender, approvalModel.Tokens, approvalModel.Block, approvalModel.TxHash)

	if err != nil {
		return err
	}

	return nil
}
