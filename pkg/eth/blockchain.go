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
	"errors"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/eth/converters"
	"golang.org/x/net/context"
)

var ErrEmptyHeader = errors.New("empty header returned over RPC")

const MAX_BATCH_SIZE = 100

type BlockChain struct {
	ethClient            core.EthClient
	headerConverter      converters.HeaderConverter
	node                 core.Node
	rpcClient            core.RpcClient
	transactionConverter converters.TransactionConverter
}

func NewBlockChain(ethClient core.EthClient, rpcClient core.RpcClient, node core.Node, converter converters.TransactionConverter) *BlockChain {
	return &BlockChain{
		ethClient:            ethClient,
		headerConverter:      converters.HeaderConverter{},
		node:                 node,
		rpcClient:            rpcClient,
		transactionConverter: converter,
	}
}

func (blockChain *BlockChain) GetEthLogsWithCustomQuery(query ethereum.FilterQuery) ([]types.Log, error) {
	gethLogs, err := blockChain.ethClient.FilterLogs(context.Background(), query)
	if err != nil {
		return []types.Log{}, err
	}
	return gethLogs, nil
}

func (blockChain *BlockChain) GetHeaderByNumber(blockNumber int64) (header core.Header, err error) {
	if blockChain.node.NetworkID == core.KOVAN_NETWORK_ID {
		return blockChain.getPOAHeader(blockNumber)
	}
	return blockChain.getPOWHeader(blockNumber)
}

func (blockChain *BlockChain) GetHeadersByNumbers(blockNumbers []int64) (header []core.Header, err error) {
	if blockChain.node.NetworkID == core.KOVAN_NETWORK_ID {
		return blockChain.getPOAHeaders(blockNumbers)
	}
	return blockChain.getPOWHeaders(blockNumbers)
}

func (blockChain *BlockChain) GetTransactions(transactionHashes []common.Hash) ([]core.TransactionModel, error) {
	numTransactions := len(transactionHashes)
	var batch []core.BatchElem
	transactions := make([]core.RpcTransaction, numTransactions)

	for index, transactionHash := range transactionHashes {
		batchElem := core.BatchElem{
			Method: "eth_getTransactionByHash",
			Result: &transactions[index],
			Args:   []interface{}{transactionHash},
		}
		batch = append(batch, batchElem)
	}

	rpcErr := blockChain.rpcClient.BatchCall(batch)
	if rpcErr != nil {
		return []core.TransactionModel{}, rpcErr
	}

	return blockChain.transactionConverter.ConvertRpcTransactionsToModels(transactions)
}

func (blockChain *BlockChain) LastBlock() (*big.Int, error) {
	block, err := blockChain.ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return big.NewInt(0), err
	}
	return block.Number, err
}

func (blockChain *BlockChain) BatchGetStorageAt(account common.Address, keys []common.Hash, blockNumber *big.Int) (map[common.Hash][]byte, error) {
	numStorageValues := len(keys)
	var batch []core.BatchElem
	storageValues := make([]hexutil.Bytes, numStorageValues)

	for index, key := range keys {
		batchElem := core.BatchElem{
			Method: "eth_getStorageAt",
			Args:   []interface{}{account.Hex(), key.Hex(), hexutil.EncodeBig(blockNumber)},
			Result: &storageValues[index],
			Error:  nil,
		}
		batch = append(batch, batchElem)
	}

	rpcErr := blockChain.rpcClient.BatchCall(batch)
	if rpcErr != nil {
		return nil, rpcErr
	}

	result := make(map[common.Hash][]byte)
	for index, key := range keys {
		result[key] = storageValues[index]
	}
	return result, nil
}

func (blockChain *BlockChain) Node() core.Node {
	return blockChain.node
}

func (blockChain *BlockChain) getPOAHeader(blockNumber int64) (header core.Header, err error) {
	var POAHeader core.POAHeader
	blockNumberArg := hexutil.EncodeBig(big.NewInt(blockNumber))
	includeTransactions := false
	err = blockChain.rpcClient.CallContext(context.Background(), &POAHeader, "eth_getBlockByNumber", blockNumberArg, includeTransactions)
	if err != nil {
		return header, err
	}
	if POAHeader.Number == nil {
		return header, ErrEmptyHeader
	}
	return blockChain.headerConverter.Convert(&types.Header{
		ParentHash:  POAHeader.ParentHash,
		UncleHash:   POAHeader.UncleHash,
		Coinbase:    POAHeader.Coinbase,
		Root:        POAHeader.Root,
		TxHash:      POAHeader.TxHash,
		ReceiptHash: POAHeader.ReceiptHash,
		Bloom:       POAHeader.Bloom,
		Difficulty:  POAHeader.Difficulty.ToInt(),
		Number:      POAHeader.Number.ToInt(),
		GasLimit:    uint64(POAHeader.GasLimit),
		GasUsed:     uint64(POAHeader.GasUsed),
		Time:        uint64(POAHeader.Time),
		Extra:       POAHeader.Extra,
	}, POAHeader.Hash.String()), nil
}

func (blockChain *BlockChain) getPOAHeaders(blockNumbers []int64) (headers []core.Header, err error) {

	var batch []core.BatchElem
	var POAHeaders [MAX_BATCH_SIZE]core.POAHeader
	includeTransactions := false

	for index, blockNumber := range blockNumbers {

		if index >= MAX_BATCH_SIZE {
			break
		}

		blockNumberArg := hexutil.EncodeBig(big.NewInt(blockNumber))

		batchElem := core.BatchElem{
			Method: "eth_getBlockByNumber",
			Result: &POAHeaders[index],
			Args:   []interface{}{blockNumberArg, includeTransactions},
		}

		batch = append(batch, batchElem)
	}

	err = blockChain.rpcClient.BatchCall(batch)
	if err != nil {
		return headers, err
	}

	for _, POAHeader := range POAHeaders {
		var header core.Header
		//Header.Number of the newest block will return nil.
		if _, err := strconv.ParseUint(POAHeader.Number.ToInt().String(), 16, 64); err == nil {
			header = blockChain.headerConverter.Convert(&types.Header{
				ParentHash:  POAHeader.ParentHash,
				UncleHash:   POAHeader.UncleHash,
				Coinbase:    POAHeader.Coinbase,
				Root:        POAHeader.Root,
				TxHash:      POAHeader.TxHash,
				ReceiptHash: POAHeader.ReceiptHash,
				Bloom:       POAHeader.Bloom,
				Difficulty:  POAHeader.Difficulty.ToInt(),
				Number:      POAHeader.Number.ToInt(),
				GasLimit:    uint64(POAHeader.GasLimit),
				GasUsed:     uint64(POAHeader.GasUsed),
				Time:        uint64(POAHeader.Time),
				Extra:       POAHeader.Extra,
			}, POAHeader.Hash.String())

			headers = append(headers, header)
		}
	}

	return headers, err
}

func (blockChain *BlockChain) getPOWHeader(blockNumber int64) (header core.Header, err error) {
	gethHeader, err := blockChain.ethClient.HeaderByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		return header, err
	}
	return blockChain.headerConverter.Convert(gethHeader, gethHeader.Hash().String()), nil
}

func (blockChain *BlockChain) getPOWHeaders(blockNumbers []int64) (headers []core.Header, err error) {
	var batch []core.BatchElem
	var POWHeaders [MAX_BATCH_SIZE]types.Header
	includeTransactions := false

	for index, blockNumber := range blockNumbers {

		if index >= MAX_BATCH_SIZE {
			break
		}

		blockNumberArg := hexutil.EncodeBig(big.NewInt(blockNumber))

		batchElem := core.BatchElem{
			Method: "eth_getBlockByNumber",
			Result: &POWHeaders[index],
			Args:   []interface{}{blockNumberArg, includeTransactions},
		}

		batch = append(batch, batchElem)
	}

	err = blockChain.rpcClient.BatchCall(batch)
	if err != nil {
		return headers, err
	}

	for _, POWHeader := range POWHeaders {
		if POWHeader.Number != nil {
			header := blockChain.headerConverter.Convert(&POWHeader, POWHeader.Hash().String())
			headers = append(headers, header)
		}
	}

	return headers, err
}
