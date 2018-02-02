package postgres

import (
	"database/sql"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

func (db DB) CreateContract(contract core.Contract) error {
	abi := contract.Abi
	var abiToInsert *string
	if abi != "" {
		abiToInsert = &abi
	}
	_, err := db.DB.Exec(
		`INSERT INTO watched_contracts (contract_hash, contract_abi)
				VALUES ($1, $2)
				ON CONFLICT (contract_hash)
				  DO UPDATE
					SET contract_hash = $1, contract_abi = $2
				`, contract.Hash, abiToInsert)
	if err != nil {
		return ErrDBInsertFailed
	}
	return nil
}

func (db DB) ContractExists(contractHash string) bool {
	var exists bool
	db.DB.QueryRow(
		`SELECT exists(
                   SELECT 1
                   FROM watched_contracts
                   WHERE contract_hash = $1)`, contractHash).Scan(&exists)
	return exists
}

func (db DB) FindContract(contractHash string) (core.Contract, error) {
	var hash string
	var abi string
	contract := db.DB.QueryRow(
		`SELECT contract_hash, contract_abi FROM watched_contracts WHERE contract_hash=$1`, contractHash)
	err := contract.Scan(&hash, &abi)
	if err == sql.ErrNoRows {
		return core.Contract{}, repositories.ErrContractDoesNotExist(contractHash)
	}
	savedContract := db.addTransactions(core.Contract{Hash: hash, Abi: abi})
	return savedContract, nil
}

func (db DB) addTransactions(contract core.Contract) core.Contract {
	transactionRows, _ := db.DB.Queryx(`
            SELECT tx_hash,
                   tx_nonce,
                   tx_to,
                   tx_from,
                   tx_gaslimit,
                   tx_gasprice,
                   tx_value,
                   tx_input_data
            FROM transactions
            WHERE tx_to = $1
            ORDER BY block_id DESC`, contract.Hash)
	transactions := db.loadTransactions(transactionRows)
	savedContract := core.Contract{Hash: contract.Hash, Transactions: transactions, Abi: contract.Abi}
	return savedContract
}
