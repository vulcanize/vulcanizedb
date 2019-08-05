// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repositories

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type HeaderSyncReceiptRepository struct{}

func (HeaderSyncReceiptRepository) CreateHeaderSyncReceiptInTx(headerID, transactionID int64, receipt core.Receipt, tx *sqlx.Tx) (int64, error) {
	var receiptId int64
	addressId, getAddressErr := AddressRepository{}.GetOrCreateAddressInTransaction(tx, receipt.ContractAddress)
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
