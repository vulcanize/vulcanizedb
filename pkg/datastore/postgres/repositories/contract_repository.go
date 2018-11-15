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

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type ContractRepository struct {
	*postgres.DB
}

func (contractRepository ContractRepository) CreateContract(contract core.Contract) error {
	abi := contract.Abi
	var abiToInsert *string
	if abi != "" {
		abiToInsert = &abi
	}
	_, err := contractRepository.DB.Exec(
		`INSERT INTO watched_contracts (contract_hash, contract_abi)
				VALUES ($1, $2)
				ON CONFLICT (contract_hash)
				  DO UPDATE
					SET contract_hash = $1, contract_abi = $2
				`, contract.Hash, abiToInsert)
	if err != nil {
		return postgres.ErrDBInsertFailed
	}
	return nil
}

func (contractRepository ContractRepository) ContractExists(contractHash string) bool {
	var exists bool
	contractRepository.DB.QueryRow(
		`SELECT exists(
                   SELECT 1
                   FROM watched_contracts
                   WHERE contract_hash = $1)`, contractHash).Scan(&exists)
	return exists
}

func (contractRepository ContractRepository) GetContract(contractHash string) (core.Contract, error) {
	var hash string
	var abi string
	contract := contractRepository.DB.QueryRow(
		`SELECT contract_hash, contract_abi FROM watched_contracts WHERE contract_hash=$1`, contractHash)
	err := contract.Scan(&hash, &abi)
	if err == sql.ErrNoRows {
		return core.Contract{}, datastore.ErrContractDoesNotExist(contractHash)
	}
	savedContract := contractRepository.addTransactions(core.Contract{Hash: hash, Abi: abi})
	return savedContract, nil
}

func (contractRepository ContractRepository) addTransactions(contract core.Contract) core.Contract {
	transactionRows, _ := contractRepository.DB.Queryx(`
            SELECT hash,
                   nonce,
                   tx_to,
                   tx_from,
                   gaslimit,
                   gasprice,
                   value,
                   input_data
            FROM transactions
            WHERE tx_to = $1
            ORDER BY block_id DESC`, contract.Hash)
	blockRepository := &BlockRepository{contractRepository.DB}
	transactions := blockRepository.LoadTransactions(transactionRows)
	savedContract := core.Contract{Hash: contract.Hash, Transactions: transactions, Abi: contract.Abi}
	return savedContract
}
