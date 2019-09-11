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
	"github.com/ethereum/go-ethereum/log"
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type HeaderSyncReceiptRepository struct{}

func (HeaderSyncReceiptRepository) CreateHeaderSyncReceiptInTx(headerID, transactionID int64, receipt core.Receipt, tx *sqlx.Tx) (int64, error) {
	var receiptId int64
	addressId, getAddressErr := GetOrCreateAddressInTransaction(tx, receipt.ContractAddress)
	if getAddressErr != nil {
		log.Error("createReceipt: Error getting address id: ", getAddressErr)
		return receiptId, getAddressErr
	}
	err := tx.QueryRowx(`INSERT INTO public.header_sync_receipts
               (header_id, transaction_id, contract_address_id, cumulative_gas_used, gas_used, state_root, status, tx_hash, rlp)
               VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			   ON CONFLICT (header_id, transaction_id) DO UPDATE
			   SET (contract_address_id, cumulative_gas_used, gas_used, state_root, status, tx_hash, rlp) = ($3, $4::NUMERIC, $5::NUMERIC, $6, $7, $8, $9)
               RETURNING id`,
		headerID, transactionID, addressId, receipt.CumulativeGasUsed, receipt.GasUsed, receipt.StateRoot, receipt.Status, receipt.TxHash, receipt.Rlp).Scan(&receiptId)
	if err != nil {
		log.Error("header_repository: error inserting receipt: ", err)
		return receiptId, err
	}
	return receiptId, err
}
