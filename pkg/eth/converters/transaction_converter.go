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

package converters

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

type TransactionConverter interface {
	ConvertRpcTransactionsToModels(transactions []core.RpcTransaction) ([]core.TransactionModel, error)
}

type transactionConverter struct {
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

func NewTransactionConverter(client core.EthClient) TransactionConverter {
	return &transactionConverter{client: client}
}

func (converter *transactionConverter) ConvertRpcTransactionsToModels(transactions []core.RpcTransaction) ([]core.TransactionModel, error) {
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

func hexToBigInt(hex string) (*big.Int, error) {
	result := big.NewInt(0)
	_, scanErr := fmt.Sscan(hex, result)
	if scanErr != nil {
		return nil, scanErr
	}
	return result, nil
}
