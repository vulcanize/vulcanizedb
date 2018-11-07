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

package fakes

import (
	"math/big"

	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockBlockChain struct {
	fetchContractDataErr               error
	fetchContractDataPassedAbi         string
	fetchContractDataPassedAddress     string
	fetchContractDataPassedMethod      string
	fetchContractDataPassedMethodArgs  []interface{}
	fetchContractDataPassedResult      interface{}
	fetchContractDataPassedBlockNumber int64
	getBlockByNumberErr                error
	lastBlock                          *big.Int
	node                               core.Node
}

func NewMockBlockChain() *MockBlockChain {
	return &MockBlockChain{
		node: core.Node{GenesisBlock: "GENESIS", NetworkID: 1, ID: "x123", ClientName: "Geth"},
	}
}

func (blockChain *MockBlockChain) SetFetchContractDataErr(err error) {
	blockChain.fetchContractDataErr = err
}

func (blockChain *MockBlockChain) SetLastBlock(blockNumber *big.Int) {
	blockChain.lastBlock = blockNumber
}

func (blockChain *MockBlockChain) SetGetBlockByNumberErr(err error) {
	blockChain.getBlockByNumberErr = err
}

func (blockChain *MockBlockChain) GetHeaderByNumber(blockNumber int64) (core.Header, error) {
	return core.Header{BlockNumber: blockNumber}, nil
}

func (blockChain *MockBlockChain) FetchContractData(abiJSON string, address string, method string, methodArgs []interface{}, result interface{}, blockNumber int64) error {
	blockChain.fetchContractDataPassedAbi = abiJSON
	blockChain.fetchContractDataPassedAddress = address
	blockChain.fetchContractDataPassedMethod = method
	blockChain.fetchContractDataPassedMethodArgs = methodArgs
	blockChain.fetchContractDataPassedResult = result
	blockChain.fetchContractDataPassedBlockNumber = blockNumber
	return blockChain.fetchContractDataErr
}

func (blockChain *MockBlockChain) CallContract(contractHash string, input []byte, blockNumber *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (blockChain *MockBlockChain) LastBlock() *big.Int {
	return blockChain.lastBlock
}

func (blockChain *MockBlockChain) GetLogs(contract core.Contract, startingBlock *big.Int, endingBlock *big.Int) ([]core.Log, error) {
	return []core.Log{}, nil
}

func (blockChain *MockBlockChain) Node() core.Node {
	return blockChain.node
}

func (blockChain *MockBlockChain) GetBlockByNumber(blockNumber int64) (core.Block, error) {
	return core.Block{Number: blockNumber}, blockChain.getBlockByNumberErr
}

// TODO: handle methodArg being nil (can't match nil to nil in Gomega)
func (blockChain *MockBlockChain) AssertFetchContractDataCalledWith(abiJSON string, address string, method string, methodArgs []interface{}, result interface{}, blockNumber int64) {
	Expect(blockChain.fetchContractDataPassedAbi).To(Equal(abiJSON))
	Expect(blockChain.fetchContractDataPassedAddress).To(Equal(address))
	Expect(blockChain.fetchContractDataPassedMethod).To(Equal(method))
	Expect(blockChain.fetchContractDataPassedResult).To(Equal(result))
	Expect(blockChain.fetchContractDataPassedBlockNumber).To(Equal(blockNumber))
}
