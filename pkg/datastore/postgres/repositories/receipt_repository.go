// VulcanizeDB
// Copyright © 2019 Vulcanize

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
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type ReceiptRepository struct {
	*postgres.DB
}

func (receiptRepository ReceiptRepository) CreateReceiptsAndLogs(blockId int64, receipts []core.Receipt) error {
	tx, err := receiptRepository.DB.Beginx()
	if err != nil {
		return err
	}
	for _, receipt := range receipts {
		receiptId, err := receiptRepository.CreateReceipt(blockId, receipt, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		if len(receipt.Logs) > 0 {
			err = createLogs(receipt.Logs, receiptId, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func createLogs(logs []core.Log, receiptId int64, tx *sqlx.Tx) error {
	for _, log := range logs {
		_, err := tx.Exec(
			`INSERT INTO logs (block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data, receipt_id)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                `,
			log.BlockNumber, log.Address, log.TxHash, log.Index, log.Topics[0], log.Topics[1], log.Topics[2], log.Topics[3], log.Data, receiptId,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

//TODO: test that creating the address should be in the transaction
func (ReceiptRepository) CreateReceipt(blockId int64, receipt core.Receipt, tx *sqlx.Tx) (int64, error) {
	var receiptId int64
	addressId, getAddressErr := AddressRepository{}.GetOrCreateAddressInTransaction(tx, receipt.ContractAddress)
	if getAddressErr != nil {
		logrus.Error("createReceipt: Error getting address id: ", getAddressErr)
		return receiptId, getAddressErr
	}
	err := tx.QueryRow(
		`INSERT INTO full_sync_receipts
               (contract_address_id, tx_hash, cumulative_gas_used, gas_used, state_root, status, block_id)
               VALUES ($1, $2, $3, $4, $5, $6, $7)
               RETURNING id`,
		addressId, receipt.TxHash, receipt.CumulativeGasUsed, receipt.GasUsed, receipt.StateRoot, receipt.Status, blockId).Scan(&receiptId)
	if err != nil {
		tx.Rollback()
		logrus.Warning("CreateReceipt: error inserting receipt: ", err)
		return receiptId, err
	}
	return receiptId, nil
}

func (receiptRepository ReceiptRepository) GetReceipt(txHash string) (core.Receipt, error) {
	row := receiptRepository.DB.QueryRow(
		`SELECT contract_address_id,
                       tx_hash,
                       cumulative_gas_used,
                       gas_used,
                       state_root,
                       status
                FROM full_sync_receipts
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
	var cumulativeGasUsed uint64
	var gasUsed uint64
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
