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

package fakes

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockBlockRepository struct {
	CreateOrUpdateBlockCallCount                 int
	CreateOrUpdateBlockPassedBlock               core.Block
	CreateOrUpdateBlockPassedBlockNumbers        []int64
	createOrUpdateBlockReturnErr                 error
	createOrUpdateBlockReturnInt                 int64
	MissingBlockNumbersPassedEndingBlockNumber   int64
	MissingBlockNumbersPassedNodeId              string
	MissingBlockNumbersPassedStartingBlockNumber int64
	missingBlockNumbersReturnArray               []int64
	SetBlockStatusPassedChainHead                int64
}

func NewMockBlockRepository() *MockBlockRepository {
	return &MockBlockRepository{
		CreateOrUpdateBlockCallCount:                 0,
		CreateOrUpdateBlockPassedBlock:               core.Block{},
		CreateOrUpdateBlockPassedBlockNumbers:        nil,
		createOrUpdateBlockReturnErr:                 nil,
		createOrUpdateBlockReturnInt:                 0,
		MissingBlockNumbersPassedEndingBlockNumber:   0,
		MissingBlockNumbersPassedNodeId:              "",
		MissingBlockNumbersPassedStartingBlockNumber: 0,
		missingBlockNumbersReturnArray:               nil,
		SetBlockStatusPassedChainHead:                0,
	}
}

func (repository *MockBlockRepository) SetCreateOrUpdateBlockReturnVals(i int64, err error) {
	repository.createOrUpdateBlockReturnInt = i
	repository.createOrUpdateBlockReturnErr = err
}

func (repository *MockBlockRepository) SetMissingBlockNumbersReturnArray(returnArray []int64) {
	repository.missingBlockNumbersReturnArray = returnArray
}

func (repository *MockBlockRepository) CreateOrUpdateBlock(block core.Block) (int64, error) {
	repository.CreateOrUpdateBlockCallCount++
	repository.CreateOrUpdateBlockPassedBlock = block
	repository.CreateOrUpdateBlockPassedBlockNumbers = append(repository.CreateOrUpdateBlockPassedBlockNumbers, block.Number)
	return repository.createOrUpdateBlockReturnInt, repository.createOrUpdateBlockReturnErr
}

func (repository *MockBlockRepository) GetBlock(blockNumber int64) (core.Block, error) {
	return core.Block{Number: blockNumber}, nil
}

func (repository *MockBlockRepository) MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64, nodeId string) []int64 {
	repository.MissingBlockNumbersPassedStartingBlockNumber = startingBlockNumber
	repository.MissingBlockNumbersPassedEndingBlockNumber = endingBlockNumber
	repository.MissingBlockNumbersPassedNodeId = nodeId
	return repository.missingBlockNumbersReturnArray
}

func (repository *MockBlockRepository) SetBlocksStatus(chainHead int64) error {
	repository.SetBlockStatusPassedChainHead = chainHead
	return nil
}
