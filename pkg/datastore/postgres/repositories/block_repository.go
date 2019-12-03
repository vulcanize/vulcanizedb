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
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/sirupsen/logrus"
)

const (
	blocksFromHeadBeforeFinal = 20
)

var ErrBlockExists = errors.New("Won't add block that already exists.")

type BlockRepository struct {
	database *postgres.DB
}

func NewBlockRepository(database *postgres.DB) *BlockRepository {
	return &BlockRepository{database: database}
}

func (repository BlockRepository) SetBlocksStatus(chainHead int64) error {
	cutoff := chainHead - blocksFromHeadBeforeFinal
	_, err := repository.database.Exec(`
                  UPDATE blocks SET is_final = TRUE
                  WHERE is_final = FALSE AND number < $1`,
		cutoff)

	return err
}

func (repository BlockRepository) CreateOrUpdateBlock(block core.Block) (int64, error) {
	var err error
	var blockId int64
	retrievedBlockHash, ok := repository.getBlockHash(block)
	if !ok {
		return repository.insertBlock(block)
	}
	if ok && retrievedBlockHash != block.Hash {
		err = repository.removeBlock(block.Number)
		if err != nil {
			return 0, err
		}
		return repository.insertBlock(block)
	}
	return blockId, ErrBlockExists
}

func (repository BlockRepository) MissingBlockNumbers(startingBlockNumber, highestBlockNumber int64) []int64 {
	numbers := make([]int64, 0)
	err := repository.database.Select(&numbers,
		`SELECT all_block_numbers
          FROM (
              SELECT generate_series($1::INT, $2::INT) AS all_block_numbers) series
          WHERE all_block_numbers NOT IN (
              SELECT number FROM blocks WHERE eth_node_id = $3
		  )`,
		startingBlockNumber, highestBlockNumber, repository.database.NodeID)
	if err != nil {
		logrus.Errorf("MissingBlockNumbers: error getting blocks: %s", err.Error())
	}
	return numbers
}

func (repository BlockRepository) GetBlock(blockNumber int64) (core.Block, error) {
	blockRows := repository.database.QueryRowx(
		`SELECT id,
                       number,
                       gas_limit,
                       gas_used,
                       time,
                       difficulty,
                       hash,
                       nonce,
                       parent_hash,
                       size,
                       uncle_hash,
                       is_final,
                       miner,
                       extra_data,
                       reward,
                       uncles_reward
               FROM blocks
               WHERE eth_node_id = $1 AND number = $2`, repository.database.NodeID, blockNumber)
	savedBlock, err := repository.loadBlock(blockRows)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return core.Block{}, datastore.ErrBlockDoesNotExist(blockNumber)
		default:
			logrus.Error("GetBlock: error loading blocks: ", err)
			return savedBlock, err
		}
	}
	return savedBlock, nil
}

func (repository BlockRepository) insertBlock(block core.Block) (int64, error) {
	var blockId int64
	tx, beginErr := repository.database.Beginx()
	if beginErr != nil {
		return 0, postgres.ErrBeginTransactionFailed(beginErr)
	}
	insertBlockErr := tx.QueryRow(
		`INSERT INTO blocks
                (eth_node_id, number, gas_limit, gas_used, time, difficulty, hash, nonce, parent_hash, size, uncle_hash, is_final, miner, extra_data, reward, uncles_reward)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
                RETURNING id `,
		repository.database.NodeID,
		block.Number,
		block.GasLimit,
		block.GasUsed,
		block.Time,
		block.Difficulty,
		block.Hash,
		block.Nonce,
		block.ParentHash,
		block.Size,
		block.UncleHash,
		block.IsFinal,
		block.Miner,
		block.ExtraData,
		nullStringToZero(block.Reward),
		nullStringToZero(block.UnclesReward)).
		Scan(&blockId)
	if insertBlockErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logrus.Error("failed to rollback transaction: ", rollbackErr)
		}
		return 0, postgres.ErrDBInsertFailed(insertBlockErr)
	}
	if len(block.Uncles) > 0 {
		insertUncleErr := repository.createUncles(tx, blockId, block.Hash, block.Uncles)
		if insertUncleErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Errorf("error rolling back transaction: %s", rollbackErr.Error())
			}
			return 0, postgres.ErrDBInsertFailed(insertUncleErr)
		}
	}
	if len(block.Transactions) > 0 {
		insertTxErr := repository.createTransactions(tx, blockId, block.Transactions)
		if insertTxErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Warn("failed to rollback transaction: ", rollbackErr)
			}
			return 0, postgres.ErrDBInsertFailed(insertTxErr)
		}
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logrus.Warn("failed to rollback transaction: ", rollbackErr)
		}
		return 0, commitErr
	}
	return blockId, nil
}

func (repository BlockRepository) createUncles(tx *sqlx.Tx, blockId int64, blockHash string, uncles []core.Uncle) error {
	for _, uncle := range uncles {
		err := repository.createUncle(tx, blockId, uncle)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repository BlockRepository) createUncle(tx *sqlx.Tx, blockId int64, uncle core.Uncle) error {
	_, err := tx.Exec(
		`INSERT INTO uncles
       (hash, block_id, reward, miner, raw, block_timestamp, eth_node_id)
       VALUES ($1, $2, $3, $4, $5, $6, $7::NUMERIC)
       RETURNING id`,
		uncle.Hash, blockId, nullStringToZero(uncle.Reward), uncle.Miner, uncle.Raw, uncle.Timestamp, repository.database.NodeID)
	return err
}

func (repository BlockRepository) createTransactions(tx *sqlx.Tx, blockId int64, transactions []core.TransactionModel) error {
	for _, transaction := range transactions {
		err := repository.createTransaction(tx, blockId, transaction)
		if err != nil {
			return err
		}
	}
	return nil
}

//Fields like value lose precision if converted to
//int64 so convert to string instead. But nil
//big.Int -> string = "" so convert to "0"
func nullStringToZero(s string) string {
	if s == "" {
		return "0"
	}
	return s
}

func (repository BlockRepository) createTransaction(tx *sqlx.Tx, blockId int64, transaction core.TransactionModel) error {
	_, err := tx.Exec(
		`INSERT INTO full_sync_transactions
       (block_id, gas_limit, gas_price, hash, input_data, nonce, raw, tx_from, tx_index, tx_to, "value")
       VALUES ($1, $2::NUMERIC, $3::NUMERIC, $4, $5, $6::NUMERIC, $7,  $8, $9::NUMERIC, $10, $11::NUMERIC)
       RETURNING id`, blockId, transaction.GasLimit, transaction.GasPrice, transaction.Hash, transaction.Data,
		transaction.Nonce, transaction.Raw, transaction.From, transaction.TxIndex, transaction.To, nullStringToZero(transaction.Value))
	if err != nil {
		return err
	}
	if hasReceipt(transaction) {
		receiptRepo := FullSyncReceiptRepository{}
		receiptId, err := receiptRepo.CreateFullSyncReceiptInTx(blockId, transaction.Receipt, tx)
		if err != nil {
			return err
		}
		if hasLogs(transaction) {
			err = repository.createLogs(tx, transaction.Receipt.Logs, receiptId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func hasLogs(transaction core.TransactionModel) bool {
	return len(transaction.Receipt.Logs) > 0
}

func hasReceipt(transaction core.TransactionModel) bool {
	return transaction.Receipt.TxHash != ""
}

func (repository BlockRepository) getBlockHash(block core.Block) (string, bool) {
	var retrievedBlockHash string
	// TODO: handle possible error
	getErr := repository.database.Get(&retrievedBlockHash,
		`SELECT hash
               FROM blocks
               WHERE number = $1 AND eth_node_id = $2`,
		block.Number, repository.database.NodeID)
	if getErr != nil {
		logrus.Errorf("error getting block hash: %s", getErr.Error())
	}
	return retrievedBlockHash, blockExists(retrievedBlockHash)
}

func (repository BlockRepository) createLogs(tx *sqlx.Tx, logs []core.FullSyncLog, receiptId int64) error {
	for _, tlog := range logs {
		_, err := tx.Exec(
			`INSERT INTO full_sync_logs (block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data, receipt_id)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                `,
			tlog.BlockNumber, tlog.Address, tlog.TxHash, tlog.Index, tlog.Topics[0], tlog.Topics[1], tlog.Topics[2], tlog.Topics[3], tlog.Data, receiptId,
		)
		if err != nil {
			return postgres.ErrDBInsertFailed(err)
		}
	}
	return nil
}

func blockExists(retrievedBlockHash string) bool {
	return retrievedBlockHash != ""
}

func (repository BlockRepository) removeBlock(blockNumber int64) error {
	_, err := repository.database.Exec(
		`DELETE FROM blocks WHERE number=$1 AND eth_node_id=$2`,
		blockNumber, repository.database.NodeID)
	if err != nil {
		return postgres.ErrDBDeleteFailed(err)
	}
	return nil
}

func (repository BlockRepository) loadBlock(blockRows *sqlx.Row) (core.Block, error) {
	type b struct {
		ID int
		core.Block
	}
	var block b
	err := blockRows.StructScan(&block)
	if err != nil {
		logrus.Error("loadBlock: error loading block: ", err)
		return core.Block{}, err
	}
	transactionRows, err := repository.database.Queryx(`
		SELECT hash,
			gas_limit,
			gas_price,
			input_data,
			nonce,
			raw,
			tx_from,
			tx_index,
			tx_to,
			value
		FROM full_sync_transactions
		WHERE block_id = $1
		ORDER BY hash`, block.ID)
	if err != nil {
		logrus.Error("loadBlock: error fetting transactions: ", err)
		return core.Block{}, err
	}
	block.Transactions = repository.LoadTransactions(transactionRows)
	return block.Block, nil
}

func (repository BlockRepository) LoadTransactions(transactionRows *sqlx.Rows) []core.TransactionModel {
	var transactions []core.TransactionModel
	for transactionRows.Next() {
		var transaction core.TransactionModel
		err := transactionRows.StructScan(&transaction)
		if err != nil {
			logrus.Fatal(err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}
