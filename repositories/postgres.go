package repositories

import (
	"database/sql"

	"github.com/8thlight/vulcanizedb/core"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) Postgres {
	return Postgres{Db: db}
}

func (repository Postgres) FindBlockByNumber(blockNumber int64) *core.Block {
	blockRows, _ := repository.Db.Query("SELECT id, block_number, block_gaslimit, block_gasused, block_time, block_difficulty, block_hash, block_nonce, block_parenthash, block_size, uncle_hash FROM blocks")
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
	repository.Db.Get(&count, "SELECT COUNT(*) FROM blocks")
	return count
}

func (repository Postgres) CreateBlock(block core.Block) {
	insertedBlock := repository.Db.QueryRow(
		"Insert INTO blocks "+
			"(block_number, block_gaslimit, block_gasused, block_time, block_difficulty, block_hash, block_nonce, block_parenthash, block_size, uncle_hash) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		block.Number, block.GasLimit, block.GasUsed, block.Time, block.Difficulty, block.Hash, block.Nonce, block.ParentHash, block.Size, block.UncleHash)
	var blockId int64
	insertedBlock.Scan(&blockId)
	repository.createTransactions(blockId, block.Transactions)
}

func (repository Postgres) createTransactions(blockId int64, transactions []core.Transaction) {
	for _, transaction := range transactions {
		repository.Db.MustExec("Insert INTO transactions "+
			"(block_id, tx_hash, tx_nonce, tx_to, tx_gaslimit, tx_gasprice, tx_value) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			blockId, transaction.Hash, transaction.Nonce, transaction.To, transaction.GasLimit, transaction.GasPrice, transaction.Value)
	}
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
	transactions := repository.loadTransactions(blockId)
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
func (repository Postgres) loadTransactions(blockId int64) []core.Transaction {
	transactionRows, _ := repository.Db.Query("SELECT tx_hash, tx_nonce, tx_to, tx_gaslimit, tx_gasprice, tx_value FROM transactions")
	var transactions []core.Transaction
	for transactionRows.Next() {
		var hash string
		var nonce uint64
		var to string
		var gasLimit int64
		var gasPrice int64
		var value int64
		transactionRows.Scan(&hash, &nonce, &to, &gasLimit, &gasPrice, &value)
		transaction := core.Transaction{
			Hash:     hash,
			Nonce:    nonce,
			To:       to,
			GasLimit: gasLimit,
			GasPrice: gasPrice,
			Value:    value,
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}
