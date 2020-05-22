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
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/sirupsen/logrus"
)

type headerRepository struct {
	db *postgres.DB
}

func NewHeaderRepository(database *postgres.DB) headerRepository {
	return headerRepository{db: database}
}

func (repo headerRepository) CreateOrUpdateHeader(header core.Header) (int64, error) {
	var headerID int64
	err := repo.db.QueryRowx("SELECT * FROM public.get_or_create_header($1, $2, $3, $4, $5)",
		header.BlockNumber, header.Hash, header.Raw, header.Timestamp, repo.db.NodeID).Scan(&headerID)
	if err != nil {
		return headerID, fmt.Errorf("error inserting header for block %d: %w", header.BlockNumber, err)
	}
	return headerID, nil
}

func (repo headerRepository) CreateTransactions(headerID int64, transactions []core.TransactionModel) error {
	for _, transaction := range transactions {
		_, err := repo.db.Exec(`INSERT INTO public.transactions
		(header_id, hash, gas_limit, gas_price, input_data, nonce, raw, tx_from, tx_index, tx_to, "value") 
		VALUES ($1, $2, $3::NUMERIC, $4::NUMERIC, $5, $6::NUMERIC, $7, $8, $9::NUMERIC, $10, $11::NUMERIC)
		ON CONFLICT DO NOTHING`, headerID, transaction.Hash, transaction.GasLimit, transaction.GasPrice,
			transaction.Data, transaction.Nonce, transaction.Raw, transaction.From, transaction.TxIndex, transaction.To,
			transaction.Value)
		if err != nil {
			return fmt.Errorf("error creating transactions: %w", err)
		}
	}
	return nil
}

func (repo headerRepository) CreateTransactionInTx(tx *sqlx.Tx, headerID int64, transaction core.TransactionModel) (int64, error) {
	var txId int64
	err := tx.QueryRowx(`INSERT INTO public.transactions
		(header_id, hash, gas_limit, gas_price, input_data, nonce, raw, tx_from, tx_index, tx_to, "value")
		VALUES ($1, $2, $3::NUMERIC, $4::NUMERIC, $5, $6::NUMERIC, $7, $8, $9::NUMERIC, $10, $11::NUMERIC)
		ON CONFLICT (hash) DO UPDATE
		SET (gas_limit, gas_price, input_data, nonce, raw, tx_from, tx_index, tx_to, "value") = ($3::NUMERIC, $4::NUMERIC, $5, $6::NUMERIC, $7, $8, $9::NUMERIC, $10, $11::NUMERIC)
		RETURNING id`,
		headerID, transaction.Hash, transaction.GasLimit, transaction.GasPrice,
		transaction.Data, transaction.Nonce, transaction.Raw, transaction.From,
		transaction.TxIndex, transaction.To, transaction.Value).Scan(&txId)
	if err != nil {
		logrus.Error("header_repository: error inserting transaction: ", err)
	}
	return txId, err
}

func (repo headerRepository) GetHeaderByBlockNumber(blockNumber int64) (core.Header, error) {
	var header core.Header
	err := repo.db.Get(&header,
		`SELECT id, block_number, hash, raw, block_timestamp FROM headers WHERE block_number = $1`, blockNumber)
	return header, err
}

func (repo headerRepository) GetHeaderByID(id int64) (core.Header, error) {
	var header core.Header
	headerErr := repo.db.Get(&header, `SELECT id, block_number, hash, raw, block_timestamp FROM headers WHERE id = $1`, id)
	return header, headerErr
}

func (repo headerRepository) GetHeadersInRange(startingBlock, endingBlock int64) ([]core.Header, error) {
	var headers []core.Header
	err := repo.db.Select(&headers,
		`SELECT id, block_number, hash, raw, block_timestamp FROM headers WHERE block_number BETWEEN $1 AND $2 ORDER BY block_number ASC`,
		startingBlock, endingBlock)
	return headers, err
}

func (repo headerRepository) MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64) ([]int64, error) {
	numbers := make([]int64, 0)
	err := repo.db.Select(&numbers,
		`SELECT series.block_number
			FROM (SELECT generate_series($1::INT, $2::INT) AS block_number) AS series
			LEFT OUTER JOIN (SELECT block_number FROM public.headers) AS synced
			USING (block_number)
			WHERE  synced.block_number IS NULL`,
		startingBlockNumber, endingBlockNumber)
	if err != nil {
		logrus.Errorf("MissingBlockNumbers failed to get blocks between %v - %v",
			startingBlockNumber, endingBlockNumber)
	}
	return numbers, err
}

func (repo headerRepository) GetMostRecentHeaderBlockNumber() (int64, error) {
	var blockNumber int64
	err := repo.db.Get(&blockNumber,
		`SELECT block_number FROM headers ORDER BY block_number DESC LIMIT 1`)
	return blockNumber, err
}
