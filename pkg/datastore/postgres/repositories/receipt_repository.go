package repositories

import (
	"database/sql"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type ReceiptRepository struct {
	*postgres.DB
}

func (receiptRepository ReceiptRepository) GetReceipt(txHash string) (core.Receipt, error) {
	row := receiptRepository.DB.QueryRow(
		`SELECT contract_address,
                       tx_hash,
                       cumulative_gas_used,
                       gas_used,
                       state_root,
                       status
                FROM receipts
                WHERE tx_hash = $1`, txHash)
	receipt, err := loadReceipt(row)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return core.Receipt{}, datastore.ErrReceiptDoesNotExist(txHash)
		default:
			return core.Receipt{}, err
		}
	}
	return receipt, nil
}

func loadReceipt(receiptsRow *sql.Row) (core.Receipt, error) {
	var contractAddress string
	var txHash string
	var cumulativeGasUsed int64
	var gasUsed int64
	var stateRoot string
	var status int

	err := receiptsRow.Scan(&contractAddress, &txHash, &cumulativeGasUsed, &gasUsed, &stateRoot, &status)
	return core.Receipt{
		TxHash:            txHash,
		ContractAddress:   contractAddress,
		CumulativeGasUsed: cumulativeGasUsed,
		GasUsed:           gasUsed,
		StateRoot:         stateRoot,
		Status:            status,
	}, err
}
