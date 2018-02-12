package postgres

import (
	"database/sql"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

type ContractRepository struct {
	*DB
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
		return ErrDBInsertFailed
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
		return core.Contract{}, repositories.ErrContractDoesNotExist(contractHash)
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
