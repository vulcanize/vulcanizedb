package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type ContractRepository interface {
	CreateContract(contract core.Contract) error
	ContractExists(contractHash string) bool
	FindContract(contractHash string) (core.Contract, error)
}

var ErrContractDoesNotExist = func(contractHash string) error {
	return errors.New(fmt.Sprintf("Contract %v does not exist", contractHash))
}

func (repository DB) CreateContract(contract core.Contract) error {
	abi := contract.Abi
	var abiToInsert *string
	if abi != "" {
		abiToInsert = &abi
	}
	_, err := repository.Db.Exec(
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

func (repository DB) ContractExists(contractHash string) bool {
	var exists bool
	repository.Db.QueryRow(
		`SELECT exists(
                   SELECT 1
                   FROM watched_contracts
                   WHERE contract_hash = $1)`, contractHash).Scan(&exists)
	return exists
}

func (repository DB) FindContract(contractHash string) (core.Contract, error) {
	var hash string
	var abi string
	contract := repository.Db.QueryRow(
		`SELECT contract_hash, contract_abi FROM watched_contracts WHERE contract_hash=$1`, contractHash)
	err := contract.Scan(&hash, &abi)
	if err == sql.ErrNoRows {
		return core.Contract{}, ErrContractDoesNotExist(contractHash)
	}
	savedContract := repository.addTransactions(core.Contract{Hash: hash, Abi: abi})
	return savedContract, nil
}

func (repository DB) addTransactions(contract core.Contract) core.Contract {
	transactionRows, _ := repository.Db.Queryx(`
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
	transactions := repository.loadTransactions(transactionRows)
	savedContract := core.Contract{Hash: contract.Hash, Transactions: transactions, Abi: contract.Abi}
	return savedContract
}
