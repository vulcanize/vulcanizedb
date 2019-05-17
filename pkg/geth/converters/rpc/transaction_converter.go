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

package rpc

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/sync/errgroup"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	vulcCommon "github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
)

type RpcTransactionConverter struct {
	client core.EthClient
}

// raw transaction data, required for generating RLP
type transactionData struct {
	AccountNonce uint64
	Price        *big.Int
	GasLimit     uint64
	Recipient    *common.Address `rlp:"nil"` // nil means contract creation
	Amount       *big.Int
	Payload      []byte
	V            *big.Int
	R            *big.Int
	S            *big.Int
}

func NewRpcTransactionConverter(client core.EthClient) *RpcTransactionConverter {
	return &RpcTransactionConverter{client: client}
}

func (converter *RpcTransactionConverter) ConvertRpcTransactionsToModels(transactions []core.RpcTransaction) ([]core.TransactionModel, error) {
	var results []core.TransactionModel
	for _, transaction := range transactions {
		txData, convertErr := getTransactionData(transaction)
		if convertErr != nil {
			return nil, convertErr
		}
		txRLP, rlpErr := getTransactionRLP(txData)
		if rlpErr != nil {
			return nil, rlpErr
		}
		txIndex, txIndexErr := hexToBigInt(transaction.TransactionIndex)
		if txIndexErr != nil {
			return nil, txIndexErr
		}
		transactionModel := core.TransactionModel{
			Data:     txData.Payload,
			From:     transaction.From,
			GasLimit: txData.GasLimit,
			GasPrice: txData.Price.Int64(),
			Hash:     transaction.Hash,
			Nonce:    txData.AccountNonce,
			Raw:      txRLP,
			// NOTE: Header Sync transactions don't include receipt; would require separate RPC call
			To:      transaction.Recipient,
			TxIndex: txIndex.Int64(),
			Value:   txData.Amount.String(),
		}
		results = append(results, transactionModel)
	}
	return results, nil
}

func (converter *RpcTransactionConverter) ConvertBlockTransactionsToCore(gethBlock *types.Block) ([]core.TransactionModel, error) {
	var g errgroup.Group
	coreTransactions := make([]core.TransactionModel, len(gethBlock.Transactions()))

	for gethTransactionIndex, gethTransaction := range gethBlock.Transactions() {
		//https://golang.org/doc/faq#closures_and_goroutines
		transaction := gethTransaction
		transactionIndex := uint(gethTransactionIndex)
		g.Go(func() error {
			from, err := converter.client.TransactionSender(context.Background(), transaction, gethBlock.Hash(), transactionIndex)
			if err != nil {
				log.Println("transaction sender: ", err)
				return err
			}
			coreTransaction, convertErr := convertGethTransactionToModel(transaction, &from, int64(gethTransactionIndex))
			if convertErr != nil {
				return convertErr
			}
			coreTransaction, err = converter.appendReceiptToTransaction(coreTransaction)
			if err != nil {
				log.Println("receipt: ", err)
				return err
			}
			coreTransactions[transactionIndex] = coreTransaction
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Println("transactions: ", err)
		return coreTransactions, err
	}
	return coreTransactions, nil
}

func (rtc *RpcTransactionConverter) appendReceiptToTransaction(transaction core.TransactionModel) (core.TransactionModel, error) {
	gethReceipt, err := rtc.client.TransactionReceipt(context.Background(), common.HexToHash(transaction.Hash))
	if err != nil {
		return transaction, err
	}
	receipt, err := vulcCommon.ToCoreReceipt(gethReceipt)
	if err != nil {
		return transaction, err
	}
	transaction.Receipt = receipt
	return transaction, nil
}

func convertGethTransactionToModel(transaction *types.Transaction, from *common.Address, transactionIndex int64) (core.TransactionModel, error) {
	raw := bytes.Buffer{}
	encodeErr := transaction.EncodeRLP(&raw)
	if encodeErr != nil {
		return core.TransactionModel{}, encodeErr
	}
	return core.TransactionModel{
		Data:     transaction.Data(),
		From:     strings.ToLower(addressToHex(from)),
		GasLimit: transaction.Gas(),
		GasPrice: transaction.GasPrice().Int64(),
		Hash:     transaction.Hash().Hex(),
		Nonce:    transaction.Nonce(),
		Raw:      raw.Bytes(),
		To:       strings.ToLower(addressToHex(transaction.To())),
		TxIndex:  transactionIndex,
		Value:    transaction.Value().String(),
	}, nil
}

func getTransactionData(transaction core.RpcTransaction) (transactionData, error) {
	nonce, nonceErr := hexToBigInt(transaction.Nonce)
	if nonceErr != nil {
		return transactionData{}, nonceErr
	}
	gasPrice, gasPriceErr := hexToBigInt(transaction.GasPrice)
	if gasPriceErr != nil {
		return transactionData{}, gasPriceErr
	}
	gasLimit, gasLimitErr := hexToBigInt(transaction.GasLimit)
	if gasLimitErr != nil {
		return transactionData{}, gasLimitErr
	}
	recipient := common.HexToAddress(transaction.Recipient)
	amount, amountErr := hexToBigInt(transaction.Amount)
	if amountErr != nil {
		return transactionData{}, amountErr
	}
	v, vErr := hexToBigInt(transaction.V)
	if vErr != nil {
		return transactionData{}, vErr
	}
	r, rErr := hexToBigInt(transaction.R)
	if rErr != nil {
		return transactionData{}, rErr
	}
	s, sErr := hexToBigInt(transaction.S)
	if sErr != nil {
		return transactionData{}, sErr
	}
	return transactionData{
		AccountNonce: nonce.Uint64(),
		Price:        gasPrice,
		GasLimit:     gasLimit.Uint64(),
		Recipient:    &recipient,
		Amount:       amount,
		Payload:      hexutil.MustDecode(transaction.Payload),
		V:            v,
		R:            r,
		S:            s,
	}, nil
}

func getTransactionRLP(txData transactionData) ([]byte, error) {
	transactionRlp := bytes.Buffer{}
	encodeErr := rlp.Encode(&transactionRlp, txData)
	if encodeErr != nil {
		return nil, encodeErr
	}
	return transactionRlp.Bytes(), nil
}

func addressToHex(to *common.Address) string {
	if to == nil {
		return ""
	}
	return to.Hex()
}

func hexToBigInt(hex string) (*big.Int, error) {
	result := big.NewInt(0)
	_, scanErr := fmt.Sscan(hex, result)
	if scanErr != nil {
		return nil, scanErr
	}
	return result, nil
}
