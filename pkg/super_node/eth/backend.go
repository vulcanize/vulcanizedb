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

package eth

import (
	"context"
	"errors"
	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
)

var (
	errPendingBlockNumber = errors.New("pending block number not supported")
)

type Backend struct {
	Retriever *CIDRetriever
	Fetcher   *IPLDPGFetcher
	DB        *postgres.DB
}

func NewEthBackend(db *postgres.DB) (*Backend, error) {
	r := NewCIDRetriever(db)
	return &Backend{
		Retriever: r,
		DB:        db,
	}, nil
}

func (b *Backend) HeaderByNumber(ctx context.Context, blockNumber rpc.BlockNumber) (*types.Header, error) {
	var err error
	number := blockNumber.Int64()
	if blockNumber == rpc.LatestBlockNumber {
		number, err = b.Retriever.RetrieveLastBlockNumber()
		if err != nil {
			return nil, err
		}
	}
	if blockNumber == rpc.PendingBlockNumber {
		return nil, errPendingBlockNumber
	}

	// Begin tx
	tx, err := b.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			shared.Rollback(tx)
			panic(p)
		} else if err != nil {
			shared.Rollback(tx)
		} else {
			err = tx.Commit()
		}
	}()

	// Retrieve the CIDs for headers at this height
	headerCids, err := b.Retriever.RetrieveHeaderCIDs(tx, number)
	if err != nil {
		return nil, err
	}
	// If there are none, throw an error
	if len(headerCids) < 1 {
		return nil, ethereum.NotFound
	}
	// Fetch the header IPLDs for those CIDs
	headerIPLD, err := b.Fetcher.FetchHeader(tx, headerCids[0])
	if err != nil {
		return nil, err
	}
	var header types.Header
	err = rlp.DecodeBytes(headerIPLD.Data, &header)
	return &header, err // return err explicitly so it can be assigned to in the defer
}

// GetTd retrieves and returns the total difficulty at the given block hash
func (b *Backend) GetTd(blockHash common.Hash) (*big.Int, error) {
	pgStr := `SELECT td FROM eth.header_cids
			WHERE header_cids.block_hash = $1`
	var tdStr string
	err := b.DB.Get(&tdStr, pgStr, blockHash.String())
	if err != nil {
		return nil, err
	}
	td, ok := new(big.Int).SetString(tdStr, 10)
	if !ok {
		return nil, errors.New("total difficulty retrieved from Postgres cannot be converted to an integer")
	}
	return td, nil
}

// GetLogs returns all the logs for the given block hash
func (b *Backend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	// Begin tx
	tx, err := b.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			shared.Rollback(tx)
			panic(p)
		} else if err != nil {
			shared.Rollback(tx)
		} else {
			err = tx.Commit()
		}
	}()
	receiptCIDs, err := b.Retriever.RetrieveRctCIDs(tx, ReceiptFilter{}, 0, &hash, nil)
	if err != nil {
		return nil, err
	}
	if len(receiptCIDs) == 0 {
		return nil, nil
	}
	receiptIPLDs, err := b.Fetcher.FetchRcts(tx, receiptCIDs)
	if err != nil {
		return nil, err
	}
	logs := make([][]*types.Log, len(receiptIPLDs))
	for i, rctIPLD := range receiptIPLDs {
		var rct types.Receipt
		if err := rlp.DecodeBytes(rctIPLD.Data, &rct); err != nil {
			return nil, err
		}
		logs[i] = rct.Logs
	}
	return logs, err
}

// BlockByNumber returns the requested canonical block.
// Since the SuperNode can contain forked blocks, it is recommended to fetch BlockByHash as
// fetching by number can return non-deterministic results (returns the first block found at that height)
func (b *Backend) BlockByNumber(ctx context.Context, blockNumber rpc.BlockNumber) (*types.Block, error) {
	var err error
	number := blockNumber.Int64()
	if blockNumber == rpc.LatestBlockNumber {
		number, err = b.Retriever.RetrieveLastBlockNumber()
		if err != nil {
			return nil, err
		}
	}
	if blockNumber == rpc.PendingBlockNumber {
		return nil, errPendingBlockNumber
	}
	// Retrieve all the CIDs for the block
	headerCID, uncleCIDs, txCIDs, rctCIDs, err := b.Retriever.RetrieveBlockByNumber(number)
	if err != nil {
		return nil, err
	}

	// Begin tx
	tx, err := b.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			shared.Rollback(tx)
			panic(p)
		} else if err != nil {
			shared.Rollback(tx)
		} else {
			err = tx.Commit()
		}
	}()

	// Fetch and decode the header IPLD
	headerIPLD, err := b.Fetcher.FetchHeader(tx, headerCID)
	if err != nil {
		return nil, err
	}
	var header types.Header
	if err := rlp.DecodeBytes(headerIPLD.Data, &header); err != nil {
		return nil, err
	}
	// Fetch and decode the uncle IPLDs
	uncleIPLDs, err := b.Fetcher.FetchUncles(tx, uncleCIDs)
	if err != nil {
		return nil, err
	}
	var uncles []*types.Header
	for _, uncleIPLD := range uncleIPLDs {
		var uncle types.Header
		if err := rlp.DecodeBytes(uncleIPLD.Data, &uncle); err != nil {
			return nil, err
		}
		uncles = append(uncles, &uncle)
	}
	// Fetch and decode the transaction IPLDs
	txIPLDs, err := b.Fetcher.FetchTrxs(tx, txCIDs)
	if err != nil {
		return nil, err
	}
	var transactions []*types.Transaction
	for _, txIPLD := range txIPLDs {
		var transaction types.Transaction
		if err := rlp.DecodeBytes(txIPLD.Data, &transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	// Fetch and decode the receipt IPLDs
	rctIPLDs, err := b.Fetcher.FetchRcts(tx, rctCIDs)
	if err != nil {
		return nil, err
	}
	var receipts []*types.Receipt
	for _, rctIPLD := range rctIPLDs {
		var receipt types.Receipt
		if err := rlp.DecodeBytes(rctIPLD.Data, &receipt); err != nil {
			return nil, err
		}
		receipts = append(receipts, &receipt)
	}
	// Compose everything together into a complete block
	return types.NewBlock(&header, transactions, uncles, receipts), err
}

// BlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.
func (b *Backend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	// Retrieve all the CIDs for the block
	headerCID, uncleCIDs, txCIDs, rctCIDs, err := b.Retriever.RetrieveBlockByHash(hash)
	if err != nil {
		return nil, err
	}

	// Begin tx
	tx, err := b.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			shared.Rollback(tx)
			panic(p)
		} else if err != nil {
			shared.Rollback(tx)
		} else {
			err = tx.Commit()
		}
	}()

	// Fetch and decode the header IPLD
	headerIPLD, err := b.Fetcher.FetchHeader(tx, headerCID)
	if err != nil {
		return nil, err
	}
	var header types.Header
	if err := rlp.DecodeBytes(headerIPLD.Data, &header); err != nil {
		return nil, err
	}
	// Fetch and decode the uncle IPLDs
	uncleIPLDs, err := b.Fetcher.FetchUncles(tx, uncleCIDs)
	if err != nil {
		return nil, err
	}
	var uncles []*types.Header
	for _, uncleIPLD := range uncleIPLDs {
		var uncle types.Header
		if err := rlp.DecodeBytes(uncleIPLD.Data, &uncle); err != nil {
			return nil, err
		}
		uncles = append(uncles, &uncle)
	}
	// Fetch and decode the transaction IPLDs
	txIPLDs, err := b.Fetcher.FetchTrxs(tx, txCIDs)
	if err != nil {
		return nil, err
	}
	var transactions []*types.Transaction
	for _, txIPLD := range txIPLDs {
		var transaction types.Transaction
		if err := rlp.DecodeBytes(txIPLD.Data, &transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	// Fetch and decode the receipt IPLDs
	rctIPLDs, err := b.Fetcher.FetchRcts(tx, rctCIDs)
	if err != nil {
		return nil, err
	}
	var receipts []*types.Receipt
	for _, rctIPLD := range rctIPLDs {
		var receipt types.Receipt
		if err := rlp.DecodeBytes(rctIPLD.Data, &receipt); err != nil {
			return nil, err
		}
		receipts = append(receipts, &receipt)
	}
	// Compose everything together into a complete block
	return types.NewBlock(&header, transactions, uncles, receipts), err
}

// GetTransaction retrieves a tx by hash
// It also returns the blockhash, blocknumber, and tx index associated with the transaction
func (b *Backend) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, error) {
	pgStr := `SELECT transaction_cids.cid, transaction_cids.index, header_cids.block_hash, header_cids.block_number
			FROM eth.transaction_cids, eth.header_cids
			WHERE transaction_cids.header_id = header_cids.id
			AND transaction_cids.tx_hash = $1`
	var txCIDWithHeaderInfo struct {
		CID         string `db:"cid"`
		Index       int64  `db:"index"`
		BlockHash   string `db:"block_hash"`
		BlockNumber int64  `db:"block_number"`
	}
	if err := b.DB.Get(&txCIDWithHeaderInfo, pgStr, txHash.String()); err != nil {
		return nil, common.Hash{}, 0, 0, err
	}

	// Begin tx
	tx, err := b.DB.Beginx()
	if err != nil {
		return nil, common.Hash{}, 0, 0, err
	}
	defer func() {
		if p := recover(); p != nil {
			shared.Rollback(tx)
			panic(p)
		} else if err != nil {
			shared.Rollback(tx)
		} else {
			err = tx.Commit()
		}
	}()

	txIPLD, err := b.Fetcher.FetchTrxs(tx, []TxModel{{CID: txCIDWithHeaderInfo.CID}})
	if err != nil {
		return nil, common.Hash{}, 0, 0, err
	}
	if err := tx.Commit(); err != nil {
		return nil, common.Hash{}, 0, 0, err
	}
	var transaction types.Transaction
	if err := rlp.DecodeBytes(txIPLD[0].Data, &transaction); err != nil {
		return nil, common.Hash{}, 0, 0, err
	}
	return &transaction, common.HexToHash(txCIDWithHeaderInfo.BlockHash), uint64(txCIDWithHeaderInfo.BlockNumber), uint64(txCIDWithHeaderInfo.Index), err
}

// extractLogsOfInterest returns logs from the receipt IPLD
func extractLogsOfInterest(rctIPLDs []ipfs.BlockModel, wantedTopics [][]string) ([]*types.Log, error) {
	var logs []*types.Log
	for _, rctIPLD := range rctIPLDs {
		rctRLP := rctIPLD
		var rct types.Receipt
		if err := rlp.DecodeBytes(rctRLP.Data, &rct); err != nil {
			return nil, err
		}
		for _, log := range rct.Logs {
			if wanted := wantedLog(wantedTopics, log.Topics); wanted == true {
				logs = append(logs, log)
			}
		}
	}
	return logs, nil
}

// returns true if the log matches on the filter
func wantedLog(wantedTopics [][]string, actualTopics []common.Hash) bool {
	// actualTopics will always have length <= 4
	// wantedTopics will always have length 4
	matches := 0
	for i, actualTopic := range actualTopics {
		// If we have topics in this filter slot, count as a match if the actualTopic matches one of the ones in this filter slot
		if len(wantedTopics[i]) > 0 {
			matches += sliceContainsHash(wantedTopics[i], actualTopic)
		} else {
			// Filter slot is empty, not matching any topics at this slot => counts as a match
			matches++
		}
	}
	if matches == len(actualTopics) {
		return true
	}
	return false
}

// returns 1 if the slice contains the hash, 0 if it does not
func sliceContainsHash(slice []string, hash common.Hash) int {
	for _, str := range slice {
		if str == hash.String() {
			return 1
		}
	}
	return 0
}

// rpcMarshalHeader uses the generalized output filler, then adds the total difficulty field, which requires
// a `PublicEthAPI`.
func (pea *PublicEthAPI) rpcMarshalHeader(header *types.Header) (map[string]interface{}, error) {
	fields := RPCMarshalHeader(header)
	td, err := pea.B.GetTd(header.Hash())
	if err != nil {
		return nil, err
	}
	fields["totalDifficulty"] = (*hexutil.Big)(td)
	return fields, nil
}

// RPCMarshalHeader converts the given header to the RPC output.
// This function is eth/internal so we have to make our own version here...
func RPCMarshalHeader(head *types.Header) map[string]interface{} {
	return map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             head.Hash(),
		"parentHash":       head.ParentHash,
		"nonce":            head.Nonce,
		"mixHash":          head.MixDigest,
		"sha3Uncles":       head.UncleHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"miner":            head.Coinbase,
		"difficulty":       (*hexutil.Big)(head.Difficulty),
		"extraData":        hexutil.Bytes(head.Extra),
		"size":             hexutil.Uint64(head.Size()),
		"gasLimit":         hexutil.Uint64(head.GasLimit),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        hexutil.Uint64(head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}
}

// rpcMarshalBlock uses the generalized output filler, then adds the total difficulty field, which requires
// a `PublicBlockchainAPI`.
func (pea *PublicEthAPI) rpcMarshalBlock(b *types.Block, inclTx bool, fullTx bool) (map[string]interface{}, error) {
	fields, err := RPCMarshalBlock(b, inclTx, fullTx)
	if err != nil {
		return nil, err
	}
	td, err := pea.B.GetTd(b.Hash())
	if err != nil {
		return nil, err
	}
	fields["totalDifficulty"] = (*hexutil.Big)(td)
	return fields, err
}

// RPCMarshalBlock converts the given block to the RPC output which depends on fullTx. If inclTx is true transactions are
// returned. When fullTx is true the returned block contains full transaction details, otherwise it will only contain
// transaction hashes.
func RPCMarshalBlock(block *types.Block, inclTx bool, fullTx bool) (map[string]interface{}, error) {
	fields := RPCMarshalHeader(block.Header())
	fields["size"] = hexutil.Uint64(block.Size())

	if inclTx {
		formatTx := func(tx *types.Transaction) (interface{}, error) {
			return tx.Hash(), nil
		}
		if fullTx {
			formatTx = func(tx *types.Transaction) (interface{}, error) {
				return NewRPCTransactionFromBlockHash(block, tx.Hash()), nil
			}
		}
		txs := block.Transactions()
		transactions := make([]interface{}, len(txs))
		var err error
		for i, tx := range txs {
			if transactions[i], err = formatTx(tx); err != nil {
				return nil, err
			}
		}
		fields["transactions"] = transactions
	}
	uncles := block.Uncles()
	uncleHashes := make([]common.Hash, len(uncles))
	for i, uncle := range uncles {
		uncleHashes[i] = uncle.Hash()
	}
	fields["uncles"] = uncleHashes

	return fields, nil
}

// NewRPCTransactionFromBlockHash returns a transaction that will serialize to the RPC representation.
func NewRPCTransactionFromBlockHash(b *types.Block, hash common.Hash) *RPCTransaction {
	for idx, tx := range b.Transactions() {
		if tx.Hash() == hash {
			return newRPCTransactionFromBlockIndex(b, uint64(idx))
		}
	}
	return nil
}

// newRPCTransactionFromBlockIndex returns a transaction that will serialize to the RPC representation.
func newRPCTransactionFromBlockIndex(b *types.Block, index uint64) *RPCTransaction {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	return NewRPCTransaction(txs[index], b.Hash(), b.NumberU64(), index)
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash        *common.Hash    `json:"blockHash"`
	BlockNumber      *hexutil.Big    `json:"blockNumber"`
	From             common.Address  `json:"from"`
	Gas              hexutil.Uint64  `json:"gas"`
	GasPrice         *hexutil.Big    `json:"gasPrice"`
	Hash             common.Hash     `json:"hash"`
	Input            hexutil.Bytes   `json:"input"`
	Nonce            hexutil.Uint64  `json:"nonce"`
	To               *common.Address `json:"to"`
	TransactionIndex *hexutil.Uint64 `json:"transactionIndex"`
	Value            *hexutil.Big    `json:"value"`
	V                *hexutil.Big    `json:"v"`
	R                *hexutil.Big    `json:"r"`
	S                *hexutil.Big    `json:"s"`
}

// NewRPCTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewRPCTransaction(tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index uint64) *RPCTransaction {
	var signer types.Signer = types.FrontierSigner{}
	if tx.Protected() {
		signer = types.NewEIP155Signer(tx.ChainId())
	}
	from, _ := types.Sender(signer, tx)
	v, r, s := tx.RawSignatureValues()

	result := &RPCTransaction{
		From:     from,
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    hexutil.Bytes(tx.Data()), // somehow this is ending up `nil`
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(v),
		R:        (*hexutil.Big)(r),
		S:        (*hexutil.Big)(s),
	}
	if blockHash != (common.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}
	return result
}
