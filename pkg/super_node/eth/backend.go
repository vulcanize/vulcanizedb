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
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/postgres"
)

var (
	errPendingBlockNumber = errors.New("pending block number not supported")
)

type Backend struct {
	Retriever *CIDRetriever
	Fetcher   *IPLDFetcher
	DB        *postgres.DB
}

func NewEthBackend(db *postgres.DB, ipfsPath string) (*Backend, error) {
	r := NewCIDRetriever(db)
	f, err := NewIPLDFetcher(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &Backend{
		Retriever: r,
		Fetcher:   f,
		DB:        db,
	}, nil
}

func (b *Backend) HeaderByNumber(ctx context.Context, blockNumber rpc.BlockNumber) (*types.Header, error) {
	number := blockNumber.Int64()
	var err error
	if blockNumber == rpc.LatestBlockNumber {
		number, err = b.Retriever.RetrieveLastBlockNumber()
		if err != nil {
			return nil, err
		}
	}
	if blockNumber == rpc.PendingBlockNumber {
		return nil, errPendingBlockNumber
	}
	// Retrieve the CIDs for headers at this height
	tx, err := b.DB.Beginx()
	if err != nil {
		return nil, err
	}
	headerCids, err := b.Retriever.RetrieveHeaderCIDs(tx, number)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logrus.Error(err)
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	// If there are none, throw an error
	if len(headerCids) < 1 {
		return nil, fmt.Errorf("header at block %d is not available", number)
	}
	// Fetch the header IPLDs for those CIDs
	headerIPLD, err := b.Fetcher.FetchHeader(headerCids[0])
	if err != nil {
		return nil, err
	}
	// Decode the first header at this block height and return it
	// We throw an error in FetchHeaders() if the number of headers does not match the number of CIDs and we already
	// confirmed the number of CIDs is greater than 0 so there is no need to bound check the slice before accessing
	var header types.Header
	if err := rlp.DecodeBytes(headerIPLD.Data, &header); err != nil {
		return nil, err
	}
	return &header, nil
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
	tx, err := b.DB.Beginx()
	if err != nil {
		return nil, err
	}
	receiptCIDs, err := b.Retriever.RetrieveRctCIDs(tx, ReceiptFilter{}, 0, &hash, nil)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logrus.Error(err)
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	if len(receiptCIDs) == 0 {
		return nil, nil
	}
	receiptIPLDs, err := b.Fetcher.FetchRcts(receiptCIDs)
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
	return logs, nil
}

// BlockByNumber returns the requested canonical block.
// Since the SuperNode can contain forked blocks, it is recommended to fetch BlockByHash as
// fetching by number can return non-deterministic results (returns the first block found at that height)
func (b *Backend) BlockByNumber(ctx context.Context, blockNumber rpc.BlockNumber) (*types.Block, error) {
	number := blockNumber.Int64()
	var err error
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

	// Fetch and decode the header IPLD
	headerIPLD, err := b.Fetcher.FetchHeader(headerCID)
	if err != nil {
		return nil, err
	}
	var header types.Header
	if err := rlp.DecodeBytes(headerIPLD.Data, &header); err != nil {
		return nil, err
	}
	// Fetch and decode the uncle IPLDs
	uncleIPLDs, err := b.Fetcher.FetchUncles(uncleCIDs)
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
	txIPLDs, err := b.Fetcher.FetchTrxs(txCIDs)
	if err != nil {
		return nil, err
	}
	var transactions []*types.Transaction
	for _, txIPLD := range txIPLDs {
		var tx types.Transaction
		if err := rlp.DecodeBytes(txIPLD.Data, &tx); err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}
	// Fetch and decode the receipt IPLDs
	rctIPLDs, err := b.Fetcher.FetchRcts(rctCIDs)
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
	return types.NewBlock(&header, transactions, uncles, receipts), nil
}

// BlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.
func (b *Backend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	// Retrieve all the CIDs for the block
	headerCID, uncleCIDs, txCIDs, rctCIDs, err := b.Retriever.RetrieveBlockByHash(hash)
	if err != nil {
		return nil, err
	}
	// Fetch and decode the header IPLD
	headerIPLD, err := b.Fetcher.FetchHeader(headerCID)
	if err != nil {
		return nil, err
	}
	var header types.Header
	if err := rlp.DecodeBytes(headerIPLD.Data, &header); err != nil {
		return nil, err
	}
	// Fetch and decode the uncle IPLDs
	uncleIPLDs, err := b.Fetcher.FetchUncles(uncleCIDs)
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
	txIPLDs, err := b.Fetcher.FetchTrxs(txCIDs)
	if err != nil {
		return nil, err
	}
	var transactions []*types.Transaction
	for _, txIPLD := range txIPLDs {
		var tx types.Transaction
		if err := rlp.DecodeBytes(txIPLD.Data, &tx); err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}
	// Fetch and decode the receipt IPLDs
	rctIPLDs, err := b.Fetcher.FetchRcts(rctCIDs)
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
	return types.NewBlock(&header, transactions, uncles, receipts), nil
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
	txIPLD, err := b.Fetcher.FetchTrxs([]TxModel{{CID: txCIDWithHeaderInfo.CID}})
	if err != nil {
		return nil, common.Hash{}, 0, 0, err
	}
	var transaction types.Transaction
	if err := rlp.DecodeBytes(txIPLD[0].Data, &transaction); err != nil {
		return nil, common.Hash{}, 0, 0, err
	}
	return &transaction, common.HexToHash(txCIDWithHeaderInfo.BlockHash), uint64(txCIDWithHeaderInfo.BlockNumber), uint64(txCIDWithHeaderInfo.Index), nil
}
