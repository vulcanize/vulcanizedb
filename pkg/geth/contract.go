// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package geth

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidStateAttribute = errors.New("invalid state attribute")
)

func (blockChain *BlockChain) FetchContractData(abiJSON string, address string, method string, methodArgs []interface{}, result interface{}, blockNumber int64) error {
	parsed, err := ParseAbi(abiJSON)
	if err != nil {
		return err
	}
	var input []byte
	if methodArgs != nil {
		input, err = parsed.Pack(method, methodArgs...)
	} else {
		input, err = parsed.Pack(method)
	}
	if err != nil {
		return err
	}
	output, err := blockChain.callContract(address, input, big.NewInt(blockNumber))
	if err != nil {
		return err
	}
	return parsed.Unpack(result, method, output)
}

func (blockChain *BlockChain) callContract(contractHash string, input []byte, blockNumber *big.Int) ([]byte, error) {
	to := common.HexToAddress(contractHash)
	msg := ethereum.CallMsg{To: &to, Data: input}
	return blockChain.ethClient.CallContract(context.Background(), msg, blockNumber)
}
