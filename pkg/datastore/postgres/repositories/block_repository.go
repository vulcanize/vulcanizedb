package repositories

import (
	"context"
	"database/sql"

	"log"

	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	blocksFromHeadBeforeFinal = 20
)

type BlockRepository struct {
	*postgres.DB
}

func (blockRepository BlockRepository) SetBlocksStatus(chainHead int64) {
	cutoff := chainHead - blocksFromHeadBeforeFinal
	blockRepository.DB.Exec(`
                  UPDATE blocks SET is_final = TRUE
                  WHERE is_final = FALSE AND number < $1`,
		cutoff)
}

func (blockRepository BlockRepository) CreateOrUpdateBlock(block core.Block) error {
	var err error
	retrievedBlockHash, ok := blockRepository.getBlockHash(block)
	if !ok {
		err = blockRepository.insertBlock(block)
		return err
	}
	if ok && retrievedBlockHash != block.Hash {
		err = blockRepository.removeBlock(block.Number)
		if err != nil {
			return err
		}
		err = blockRepository.insertBlock(block)
		return err
	}
	return nil
}

func (blockRepository BlockRepository) MissingBlockNumbers(startingBlockNumber int64, highestBlockNumber int64) []int64 {
	numbers := make([]int64, 0)
	blockRepository.DB.Select(&numbers,
		`SELECT all_block_numbers
            FROM (
                SELECT generate_series($1::INT, $2::INT) AS all_block_numbers) series
                LEFT JOIN blocks
                    ON number = all_block_numbers
            WHERE number ISNULL`,
		startingBlockNumber,
		highestBlockNumber)
	return numbers
}

func (blockRepository BlockRepository) GetBlock(blockNumber int64) (core.Block, error) {
	blockRows := blockRepository.DB.QueryRowx(
		`SELECT id,
                       number,
                       gaslimit,
                       gasused,
                       time,
                       difficulty,
                       hash,
                       nonce,
                       parenthash,
                       size,
                       uncle_hash,
                       is_final,
                       miner,
                       extra_data,
                       reward,
                       uncles_reward
               FROM blocks
               WHERE eth_node_id = $1 AND number = $2`, blockRepository.NodeID, blockNumber)
	savedBlock, err := blockRepository.loadBlock(blockRows)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return core.Block{}, datastore.ErrBlockDoesNotExist(blockNumber)
		default:
			return savedBlock, err
		}
	}
	return savedBlock, nil
}

func (blockRepository BlockRepository) insertBlock(block core.Block) error {
	var blockId int64
	tx, _ := blockRepository.DB.BeginTx(context.Background(), nil)
	err := tx.QueryRow(
		`INSERT INTO blocks
                (eth_node_id, number, gaslimit, gasused, time, difficulty, hash, nonce, parenthash, size, uncle_hash, is_final, miner, extra_data, reward, uncles_reward)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
                RETURNING id `,
		blockRepository.NodeID, block.Number, block.GasLimit, block.GasUsed, block.Time, block.Difficulty, block.Hash, block.Nonce, block.ParentHash, block.Size, block.UncleHash, block.IsFinal, block.Miner, block.ExtraData, block.Reward, block.UnclesReward).
		Scan(&blockId)
	if err != nil {
		tx.Rollback()
		return postgres.ErrDBInsertFailed
	}
	err = blockRepository.createTransactions(tx, blockId, block.Transactions)
	if err != nil {
		tx.Rollback()
		return postgres.ErrDBInsertFailed
	}
	tx.Commit()
	return nil
}

func (blockRepository BlockRepository) createTransactions(tx *sql.Tx, blockId int64, transactions []core.Transaction) error {
	for _, transaction := range transactions {
		err := blockRepository.createTransaction(tx, blockId, transaction)
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

func (blockRepository BlockRepository) createTransaction(tx *sql.Tx, blockId int64, transaction core.Transaction) error {
	var transactionId int
	err := tx.QueryRow(
		`INSERT INTO transactions
       (block_id, hash, nonce, tx_to, tx_from, gaslimit, gasprice, value, input_data)
       VALUES ($1, $2, $3, $4, $5, $6, $7,  $8::NUMERIC, $9)
       RETURNING id`,
		blockId, transaction.Hash, transaction.Nonce, transaction.To, transaction.From, transaction.GasLimit, transaction.GasPrice, nullStringToZero(transaction.Value), transaction.Data).
		Scan(&transactionId)
	if err != nil {
		return err
	}
	if hasReceipt(transaction) {
		receiptId, err := blockRepository.createReceipt(tx, transactionId, transaction.Receipt)
		if err != nil {
			return err
		}
		if hasLogs(transaction) {
			err = blockRepository.createLogs(tx, transaction.Receipt.Logs, receiptId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func hasLogs(transaction core.Transaction) bool {
	return len(transaction.Receipt.Logs) > 0
}

func hasReceipt(transaction core.Transaction) bool {
	return transaction.Receipt.TxHash != ""
}

func (blockRepository BlockRepository) createReceipt(tx *sql.Tx, transactionId int, receipt core.Receipt) (int, error) {
	//Not currently persisting log bloom filters
	var receiptId int
	err := tx.QueryRow(
		`INSERT INTO receipts
               (contract_address, tx_hash, cumulative_gas_used, gas_used, state_root, status, transaction_id)
               VALUES ($1, $2, $3, $4, $5, $6, $7) 
               RETURNING id`,
		receipt.ContractAddress, receipt.TxHash, receipt.CumulativeGasUsed, receipt.GasUsed, receipt.StateRoot, receipt.Status, transactionId).Scan(&receiptId)
	if err != nil {
		return receiptId, err
	}
	return receiptId, nil
}

func (blockRepository BlockRepository) getBlockHash(block core.Block) (string, bool) {
	var retrievedBlockHash string
	blockRepository.DB.Get(&retrievedBlockHash,
		`SELECT hash
               FROM blocks
               WHERE number = $1 AND eth_node_id = $2`,
		block.Number, blockRepository.NodeID)
	return retrievedBlockHash, blockExists(retrievedBlockHash)
}

func (blockRepository BlockRepository) createLogs(tx *sql.Tx, logs []core.Log, receiptId int) error {
	for _, tlog := range logs {
		_, err := tx.Exec(
			`INSERT INTO logs (block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data, receipt_id)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                `,
			tlog.BlockNumber, tlog.Address, tlog.TxHash, tlog.Index, tlog.Topics[0], tlog.Topics[1], tlog.Topics[2], tlog.Topics[3], tlog.Data, receiptId,
		)
		if err != nil {
			return postgres.ErrDBInsertFailed
		}
	}
	return nil
}

func blockExists(retrievedBlockHash string) bool {
	return retrievedBlockHash != ""
}

func (blockRepository BlockRepository) removeBlock(blockNumber int64) error {
	_, err := blockRepository.DB.Exec(
		`DELETE FROM
                blocks
                WHERE number=$1 AND eth_node_id=$2`,
		blockNumber, blockRepository.NodeID)
	if err != nil {
		return postgres.ErrDBDeleteFailed
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
		return core.Block{}, err
	}
	transactionRows, err := blockRepository.DB.Queryx(`
            SELECT hash,
				   nonce,
				   tx_to,
				   tx_from,
				   gaslimit,
				   gasprice,
				   value,
				   input_data
            FROM transactions
            WHERE block_id = $1
            ORDER BY hash`, block.ID)
	if err != nil {
		return core.Block{}, err
	}
	block.Transactions = blockRepository.LoadTransactions(transactionRows)
	return block.Block, nil
}

func (blockRepository BlockRepository) LoadTransactions(transactionRows *sqlx.Rows) []core.Transaction {
	var transactions []core.Transaction
	for transactionRows.Next() {
		var transaction core.Transaction
		err := transactionRows.StructScan(&transaction)
		if err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}
