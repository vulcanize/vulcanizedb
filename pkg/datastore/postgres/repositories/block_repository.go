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
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
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

func (blockRepository BlockRepository) SetBlocksStatus(chainHead int64) error {
	cutoff := chainHead - blocksFromHeadBeforeFinal
	_, err := blockRepository.database.Exec(`
                  UPDATE eth_blocks SET is_final = TRUE
                  WHERE is_final = FALSE AND number < $1`,
		cutoff)

	return err
}

func (blockRepository BlockRepository) CreateOrUpdateBlock(block core.Block) (int64, error) {
	var err error
	var blockID int64
	retrievedBlockHash, ok := blockRepository.getBlockHash(block)
	if !ok {
		return blockRepository.insertBlock(block)
	}
	if ok && retrievedBlockHash != block.Hash {
		err = blockRepository.removeBlock(block.Number)
		if err != nil {
			return 0, err
		}
		return blockRepository.insertBlock(block)
	}
	return blockID, ErrBlockExists
}

func (blockRepository BlockRepository) MissingBlockNumbers(startingBlockNumber int64, highestBlockNumber int64, nodeID string) []int64 {
	numbers := make([]int64, 0)
	err := blockRepository.database.Select(&numbers,
		`SELECT all_block_numbers
          FROM (
              SELECT generate_series($1::INT, $2::INT) AS all_block_numbers) series
          WHERE all_block_numbers NOT IN (
			  SELECT number FROM eth_blocks WHERE eth_node_fingerprint = $3
		  ) `,
		startingBlockNumber,
		highestBlockNumber, nodeID)
	if err != nil {
		logrus.Error("MissingBlockNumbers: error getting blocks: ", err)
	}
	return numbers
}

func (blockRepository BlockRepository) GetBlock(blockNumber int64) (core.Block, error) {
	blockRows := blockRepository.database.QueryRowx(
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
               FROM eth_blocks
               WHERE eth_node_id = $1 AND number = $2`, blockRepository.database.NodeID, blockNumber)
	savedBlock, err := blockRepository.loadBlock(blockRows)
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

func (blockRepository BlockRepository) insertBlock(block core.Block) (int64, error) {
	var blockID int64
	tx, beginErr := blockRepository.database.Beginx()
	if beginErr != nil {
		return 0, postgres.ErrBeginTransactionFailed(beginErr)
	}
	insertBlockErr := tx.QueryRow(
		`INSERT INTO eth_blocks
                (eth_node_id, number, gas_limit, gas_used, time, difficulty, hash, nonce, parent_hash, size, uncle_hash, is_final, miner, extra_data, reward, uncles_reward, eth_node_fingerprint)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
                RETURNING id `,
		blockRepository.database.NodeID,
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
		nullStringToZero(block.UnclesReward),
		blockRepository.database.Node.ID).
		Scan(&blockID)
	if insertBlockErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logrus.Error("failed to rollback transaction: ", rollbackErr)
		}
		return 0, postgres.ErrDBInsertFailed(insertBlockErr)
	}
	if len(block.Uncles) > 0 {
		insertUncleErr := blockRepository.createUncles(tx, blockID, block.Hash, block.Uncles)
		if insertUncleErr != nil {
			tx.Rollback()
			return 0, postgres.ErrDBInsertFailed(insertUncleErr)
		}
	}
	if len(block.Transactions) > 0 {
		insertTxErr := blockRepository.createTransactions(tx, blockID, block.Transactions)
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
	return blockID, nil
}

func (blockRepository BlockRepository) createUncles(tx *sqlx.Tx, blockID int64, blockHash string, uncles []core.Uncle) error {
	for _, uncle := range uncles {
		err := blockRepository.createUncle(tx, blockID, uncle)
		if err != nil {
			return err
		}
	}
	return nil
}

func (blockRepository BlockRepository) createUncle(tx *sqlx.Tx, blockID int64, uncle core.Uncle) error {
	_, err := tx.Exec(
		`INSERT INTO uncles
       (hash, block_id, reward, miner, raw, block_timestamp, eth_node_id, eth_node_fingerprint)
       VALUES ($1, $2, $3, $4, $5, $6, $7::NUMERIC, $8)
       RETURNING id`,
		uncle.Hash, blockID, nullStringToZero(uncle.Reward), uncle.Miner, uncle.Raw, uncle.Timestamp, blockRepository.database.NodeID, blockRepository.database.Node.ID)
	return err
}

func (blockRepository BlockRepository) createTransactions(tx *sqlx.Tx, blockID int64, transactions []core.TransactionModel) error {
	for _, transaction := range transactions {
		err := blockRepository.createTransaction(tx, blockID, transaction)
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

func (blockRepository BlockRepository) createTransaction(tx *sqlx.Tx, blockID int64, transaction core.TransactionModel) error {
	_, err := tx.Exec(
		`INSERT INTO full_sync_transactions
       (block_id, gas_limit, gas_price, hash, input_data, nonce, raw, tx_from, tx_index, tx_to, "value")
       VALUES ($1, $2::NUMERIC, $3::NUMERIC, $4, $5, $6::NUMERIC, $7,  $8, $9::NUMERIC, $10, $11::NUMERIC)
       RETURNING id`, blockID, transaction.GasLimit, transaction.GasPrice, transaction.Hash, transaction.Data,
		transaction.Nonce, transaction.Raw, transaction.From, transaction.TxIndex, transaction.To, nullStringToZero(transaction.Value))
	if err != nil {
		return err
	}
	if hasReceipt(transaction) {
		receiptRepo := FullSyncReceiptRepository{}
		receiptID, err := receiptRepo.CreateFullSyncReceiptInTx(blockID, transaction.Receipt, tx)
		if err != nil {
			return err
		}
		if hasLogs(transaction) {
			err = blockRepository.createLogs(tx, transaction.Receipt.Logs, receiptID)
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

func (blockRepository BlockRepository) getBlockHash(block core.Block) (string, bool) {
	var retrievedBlockHash string
	// TODO: handle possible error
	blockRepository.database.Get(&retrievedBlockHash,
		`SELECT hash
               FROM eth_blocks
               WHERE number = $1 AND eth_node_id = $2`,
		block.Number, blockRepository.database.NodeID)
	return retrievedBlockHash, blockExists(retrievedBlockHash)
}

func (blockRepository BlockRepository) createLogs(tx *sqlx.Tx, logs []core.FullSyncLog, receiptID int64) error {
	for _, tlog := range logs {
		_, err := tx.Exec(
			`INSERT INTO full_sync_logs (block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data, receipt_id)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                `,
			tlog.BlockNumber, tlog.Address, tlog.TxHash, tlog.Index, tlog.Topics[0], tlog.Topics[1], tlog.Topics[2], tlog.Topics[3], tlog.Data, receiptID,
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

func (blockRepository BlockRepository) removeBlock(blockNumber int64) error {
	_, err := blockRepository.database.Exec(
		`DELETE FROM eth_blocks WHERE number=$1 AND eth_node_id=$2`,
		blockNumber, blockRepository.database.NodeID)
	if err != nil {
		return postgres.ErrDBDeleteFailed(err)
	}
	return nil
}

func (blockRepository BlockRepository) loadBlock(blockRows *sqlx.Row) (core.Block, error) {
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
	transactionRows, err := blockRepository.database.Queryx(`
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
	block.Transactions = blockRepository.LoadTransactions(transactionRows)
	return block.Block, nil
}

func (blockRepository BlockRepository) LoadTransactions(transactionRows *sqlx.Rows) []core.TransactionModel {
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
