package repositories

import (
	"database/sql"

	"context"

	"errors"

	"github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Db *sqlx.DB
}

var (
	ErrDBInsertFailed     = errors.New("postgres: insert failed")
	ErrDBConnectionFailed = errors.New("postgres: db connection failed")
)

func NewPostgres(databaseConfig config.Database) (Postgres, error) {
	connectString := config.DbConnectionString(databaseConfig)
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		return Postgres{}, ErrDBConnectionFailed
	}
	return Postgres{Db: db}, nil
}

func (repository Postgres) CreateContract(contract core.Contract) error {
	abi := contract.Abi
	var abiToInsert *string
	if abi != "" {
		abiToInsert = &abi
	}
	_, err := repository.Db.Exec(
		`INSERT INTO watched_contracts (contract_hash, contract_abi) VALUES ($1, $2)`, contract.Hash, abiToInsert)
	if err != nil {
		return ErrDBInsertFailed
	}
	return nil
}

func (repository Postgres) ContractExists(contractHash string) bool {
	var exists bool
	repository.Db.QueryRow(
		`SELECT exists(SELECT 1 FROM watched_contracts WHERE contract_hash=$1) FROM watched_contracts`, contractHash).Scan(&exists)
	return exists
}

func (repository Postgres) FindContract(contractHash string) *core.Contract {
	var savedContracts []core.Contract
	contractRows, _ := repository.Db.Query(
		`SELECT contract_hash, contract_abi FROM watched_contracts WHERE contract_hash=$1`, contractHash)
	savedContracts = repository.loadContract(contractRows)
	if len(savedContracts) > 0 {
		return &savedContracts[0]
	} else {
		return nil
	}
}

func (repository Postgres) MaxBlockNumber() int64 {
	var highestBlockNumber int64
	repository.Db.Get(&highestBlockNumber, `SELECT MAX(block_number) FROM blocks`)
	return highestBlockNumber
}

func (repository Postgres) MissingBlockNumbers(startingBlockNumber int64, highestBlockNumber int64) []int64 {
	numbers := []int64{}
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

func (repository Postgres) FindBlockByNumber(blockNumber int64) *core.Block {
	blockRows, _ := repository.Db.Query(
		`SELECT id, block_number, block_gaslimit, block_gasused, block_time, block_difficulty, block_hash, block_nonce, block_parenthash, block_size, uncle_hash FROM blocks`)
	var savedBlocks []core.Block
	for blockRows.Next() {
		savedBlock := repository.loadBlock(blockRows)
		savedBlocks = append(savedBlocks, savedBlock)
	}
	if len(savedBlocks) > 0 {
		return &savedBlocks[0]
	} else {
		return nil
	}
}

func (repository Postgres) BlockCount() int {
	var count int
	repository.Db.Get(&count, `SELECT COUNT(*) FROM blocks`)
	return count
}

func (repository Postgres) CreateBlock(block core.Block) error {
	tx, _ := repository.Db.BeginTx(context.Background(), nil)
	var blockId int64
	err := tx.QueryRow(
		`INSERT INTO blocks
			(block_number, block_gaslimit, block_gasused, block_time, block_difficulty, block_hash, block_nonce, block_parenthash, block_size, uncle_hash)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id `,
		block.Number, block.GasLimit, block.GasUsed, block.Time, block.Difficulty, block.Hash, block.Nonce, block.ParentHash, block.Size, block.UncleHash).
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
		_, err := tx.Exec(
			`INSERT INTO transactions
			(block_id, tx_hash, tx_nonce, tx_to, tx_from, tx_gaslimit, tx_gasprice, tx_value)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			blockId, transaction.Hash, transaction.Nonce, transaction.To, transaction.From, transaction.GasLimit, transaction.GasPrice, transaction.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repository Postgres) loadBlock(blockRows *sql.Rows) core.Block {
	var blockId int64
	var blockHash string
	var blockNonce string
	var blockNumber int64
	var blockParentHash string
	var blockSize int64
	var blockTime float64
	var difficulty int64
	var gasLimit float64
	var gasUsed float64
	var uncleHash string
	blockRows.Scan(&blockId, &blockNumber, &gasLimit, &gasUsed, &blockTime, &difficulty, &blockHash, &blockNonce, &blockParentHash, &blockSize, &uncleHash)
	transactionRows, _ := repository.Db.Query(`SELECT tx_hash, tx_nonce, tx_to, tx_from, tx_gaslimit, tx_gasprice, tx_value FROM transactions WHERE block_id = $1`, blockId)
	transactions := repository.loadTransactions(transactionRows)
	return core.Block{
		Difficulty:   difficulty,
		GasLimit:     int64(gasLimit),
		GasUsed:      int64(gasUsed),
		Hash:         blockHash,
		Nonce:        blockNonce,
		Number:       blockNumber,
		ParentHash:   blockParentHash,
		Size:         blockSize,
		Time:         int64(blockTime),
		Transactions: transactions,
		UncleHash:    uncleHash,
	}
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
		var value int64
		transactionRows.Scan(&hash, &nonce, &to, &from, &gasLimit, &gasPrice, &value)
		transaction := core.Transaction{
			Hash:     hash,
			Nonce:    nonce,
			To:       to,
			From:     from,
			GasLimit: gasLimit,
			GasPrice: gasPrice,
			Value:    value,
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}

func (repository Postgres) loadContract(contractRows *sql.Rows) []core.Contract {
	var savedContracts []core.Contract
	for contractRows.Next() {
		var savedContractHash string
		var savedContractAbi string
		contractRows.Scan(&savedContractHash, &savedContractAbi)
		transactionRows, _ := repository.Db.Query(`SELECT tx_hash, tx_nonce, tx_to, tx_from, tx_gaslimit, tx_gasprice, tx_value FROM transactions WHERE tx_to = $1 ORDER BY block_id desc`, savedContractHash)
		transactions := repository.loadTransactions(transactionRows)
		savedContract := core.Contract{Hash: savedContractHash, Transactions: transactions, Abi: savedContractAbi}
		savedContracts = append(savedContracts, savedContract)
	}
	return savedContracts
}
