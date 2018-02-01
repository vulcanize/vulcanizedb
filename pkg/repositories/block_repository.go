package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"errors"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type BlockRepository interface {
	CreateOrUpdateBlock(block core.Block) error
	BlockCount() int
	FindBlockByNumber(blockNumber int64) (core.Block, error)
	MaxBlockNumber() int64
	MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64) []int64
	SetBlocksStatus(chainHead int64)
}

var ErrBlockDoesNotExist = func(blockNumber int64) error {
	return errors.New(fmt.Sprintf("Block number %d does not exist", blockNumber))
}

func (repository Postgres) SetBlocksStatus(chainHead int64) {
	cutoff := chainHead - blocksFromHeadBeforeFinal
	repository.Db.Exec(`
                  UPDATE blocks SET is_final = TRUE
                  WHERE is_final = FALSE AND block_number < $1`,
		cutoff)
}

func (repository Postgres) CreateOrUpdateBlock(block core.Block) error {
	var err error
	retrievedBlockHash, ok := repository.getBlockHash(block)
	if !ok {
		err = repository.insertBlock(block)
		return err
	}
	if ok && retrievedBlockHash != block.Hash {
		err = repository.removeBlock(block.Number)
		if err != nil {
			return err
		}
		err = repository.insertBlock(block)
		return err
	}
	return nil
}

func (repository Postgres) BlockCount() int {
	var count int
	repository.Db.Get(&count, `SELECT COUNT(*) FROM blocks`)
	return count
}

func (repository Postgres) MaxBlockNumber() int64 {
	var highestBlockNumber int64
	repository.Db.Get(&highestBlockNumber, `SELECT MAX(block_number) FROM blocks`)
	return highestBlockNumber
}

func (repository Postgres) MissingBlockNumbers(startingBlockNumber int64, highestBlockNumber int64) []int64 {
	numbers := make([]int64, 0)
	repository.Db.Select(&numbers,
		`SELECT all_block_numbers
            FROM (
                SELECT generate_series($1::INT, $2::INT) AS all_block_numbers) series
                LEFT JOIN blocks
                    ON block_number = all_block_numbers
            WHERE block_number ISNULL`,
		startingBlockNumber,
		highestBlockNumber)
	return numbers
}

func (repository Postgres) FindBlockByNumber(blockNumber int64) (core.Block, error) {
	blockRows := repository.Db.QueryRow(
		`SELECT id,
                       block_number,
                       block_gaslimit,
                       block_gasused,
                       block_time,
                       block_difficulty,
                       block_hash,
                       block_nonce,
                       block_parenthash,
                       block_size,
                       uncle_hash,
                       is_final,
                       block_miner,
                       block_extra_data,
                       block_reward,
                       block_uncles_reward
               FROM blocks
               WHERE node_id = $1 AND block_number = $2`, repository.nodeId, blockNumber)
	savedBlock, err := repository.loadBlock(blockRows)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return core.Block{}, ErrBlockDoesNotExist(blockNumber)
		default:
			return savedBlock, err
		}
	}
	return savedBlock, nil
}

func (repository Postgres) insertBlock(block core.Block) error {
	var blockId int64
	tx, _ := repository.Db.BeginTx(context.Background(), nil)
	err := tx.QueryRow(
		`INSERT INTO blocks
                (node_id, block_number, block_gaslimit, block_gasused, block_time, block_difficulty, block_hash, block_nonce, block_parenthash, block_size, uncle_hash, is_final, block_miner, block_extra_data, block_reward, block_uncles_reward)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
                RETURNING id `,
		repository.nodeId, block.Number, block.GasLimit, block.GasUsed, block.Time, block.Difficulty, block.Hash, block.Nonce, block.ParentHash, block.Size, block.UncleHash, block.IsFinal, block.Miner, block.ExtraData, block.Reward, block.UnclesReward).
		Scan(&blockId)
	if err != nil {
		tx.Rollback()
		return ErrDBInsertFailed
	}
	err = repository.createTransactions(tx, blockId, block.Transactions)
	if err != nil {
		tx.Rollback()
		return ErrDBInsertFailed
	}
	tx.Commit()
	return nil
}

func (repository Postgres) createTransactions(tx *sql.Tx, blockId int64, transactions []core.Transaction) error {
	for _, transaction := range transactions {
		err := repository.createTransaction(tx, blockId, transaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repository Postgres) createTransaction(tx *sql.Tx, blockId int64, transaction core.Transaction) error {
	var transactionId int
	err := tx.QueryRow(
		`INSERT INTO transactions
       (block_id, tx_hash, tx_nonce, tx_to, tx_from, tx_gaslimit, tx_gasprice, tx_value, tx_input_data)
       VALUES ($1, $2, $3, $4, $5, $6, $7,  cast(NULLIF($8, '') AS NUMERIC), $9)
       RETURNING id`,
		blockId, transaction.Hash, transaction.Nonce, transaction.To, transaction.From, transaction.GasLimit, transaction.GasPrice, transaction.Value, transaction.Data).
		Scan(&transactionId)
	if err != nil {
		return err
	}
	if hasReceipt(transaction) {
		receiptId, err := repository.createReceipt(tx, transactionId, transaction.Receipt)
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

func hasLogs(transaction core.Transaction) bool {
	return len(transaction.Receipt.Logs) > 0
}

func hasReceipt(transaction core.Transaction) bool {
	return transaction.Receipt.TxHash != ""
}

func (repository Postgres) createReceipt(tx *sql.Tx, transactionId int, receipt core.Receipt) (int, error) {
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

func (repository Postgres) getBlockHash(block core.Block) (string, bool) {
	var retrievedBlockHash string
	repository.Db.Get(&retrievedBlockHash,
		`SELECT block_hash
               FROM blocks
               WHERE block_number = $1 AND node_id = $2`,
		block.Number, repository.nodeId)
	return retrievedBlockHash, blockExists(retrievedBlockHash)
}

func (repository Postgres) createLogs(tx *sql.Tx, logs []core.Log, receiptId int) error {
	for _, tlog := range logs {
		_, err := tx.Exec(
			`INSERT INTO logs (block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data, receipt_id)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                `,
			tlog.BlockNumber, tlog.Address, tlog.TxHash, tlog.Index, tlog.Topics[0], tlog.Topics[1], tlog.Topics[2], tlog.Topics[3], tlog.Data, receiptId,
		)
		if err != nil {
			return ErrDBInsertFailed
		}
	}
	return nil
}

func blockExists(retrievedBlockHash string) bool {
	return retrievedBlockHash != ""
}

func (repository Postgres) removeBlock(blockNumber int64) error {
	_, err := repository.Db.Exec(
		`DELETE FROM
                blocks
                WHERE block_number=$1 AND node_id=$2`,
		blockNumber, repository.nodeId)
	if err != nil {
		return ErrDBDeleteFailed
	}
	return nil
}

func (repository Postgres) loadBlock(blockRows *sql.Row) (core.Block, error) {
	var blockId int64
	var blockHash string
	var blockNonce string
	var blockNumber int64
	var blockMiner string
	var blockExtraData string
	var blockParentHash string
	var blockSize int64
	var blockTime float64
	var blockReward float64
	var difficulty int64
	var gasLimit float64
	var gasUsed float64
	var uncleHash string
	var unclesReward float64
	var isFinal bool
	err := blockRows.Scan(&blockId, &blockNumber, &gasLimit, &gasUsed, &blockTime, &difficulty, &blockHash, &blockNonce, &blockParentHash, &blockSize, &uncleHash, &isFinal, &blockMiner, &blockExtraData, &blockReward, &unclesReward)
	if err != nil {
		return core.Block{}, err
	}
	transactionRows, _ := repository.Db.Query(`
            SELECT tx_hash,
				   tx_nonce,
				   tx_to,
				   tx_from,
				   tx_gaslimit,
				   tx_gasprice,
				   tx_value,
				   tx_input_data
            FROM transactions
            WHERE block_id = $1
            ORDER BY tx_hash`, blockId)
	transactions := repository.loadTransactions(transactionRows)
	return core.Block{
		Reward:       blockReward,
		Difficulty:   difficulty,
		ExtraData:    blockExtraData,
		GasLimit:     int64(gasLimit),
		GasUsed:      int64(gasUsed),
		Hash:         blockHash,
		IsFinal:      isFinal,
		Miner:        blockMiner,
		Nonce:        blockNonce,
		Number:       blockNumber,
		ParentHash:   blockParentHash,
		Size:         blockSize,
		Time:         int64(blockTime),
		Transactions: transactions,
		UncleHash:    uncleHash,
		UnclesReward: unclesReward,
	}, nil
}

func (repository Postgres) loadTransactions(transactionRows *sql.Rows) []core.Transaction {
	var transactions []core.Transaction
	for transactionRows.Next() {
		var hash string
		var nonce uint64
		var to string
		var from string
		var gasLimit int64
		var gasPrice int64
		var inputData string
		var value string
		transactionRows.Scan(&hash, &nonce, &to, &from, &gasLimit, &gasPrice, &value, &inputData)
		transaction := core.Transaction{
			Hash:     hash,
			Nonce:    nonce,
			To:       to,
			From:     from,
			GasLimit: gasLimit,
			GasPrice: gasPrice,
			Value:    value,
			Data:     inputData,
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}
