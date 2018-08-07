// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test_data

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
)

type MockLogFetcher struct {
	FetchedContractAddress string
	FetchedTopics          [][]common.Hash
	FetchedBlocks          []int64
	FetcherError           error
	FetchedLogs            []types.Log
}

func (mlf *MockLogFetcher) FetchLogs(contractAddress string, topics [][]common.Hash, blockNumber int64) ([]types.Log, error) {
	mlf.FetchedContractAddress = contractAddress
	mlf.FetchedTopics = topics
	mlf.FetchedBlocks = append(mlf.FetchedBlocks, blockNumber)

	return mlf.FetchedLogs, mlf.FetcherError
}

func (mlf *MockLogFetcher) SetFetcherError(err error) {
	mlf.FetcherError = err
}

func (mlf *MockLogFetcher) SetFetchedLogs(logs []types.Log) {
	mlf.FetchedLogs = logs
}

type MockFlipKickConverter struct {
	ConverterContract string
	ConverterAbi      string
	LogsToConvert     []types.Log
	EntitiesToConvert []flip_kick.FlipKickEntity
	ConverterError    error
}

func (mfkc *MockFlipKickConverter) ToEntity(contractAddress string, contractAbi string, ethLog types.Log) (*flip_kick.FlipKickEntity, error) {
	mfkc.ConverterContract = contractAddress
	mfkc.ConverterAbi = contractAbi
	mfkc.LogsToConvert = append(mfkc.LogsToConvert, ethLog)
	return &FlipKickEntity, mfkc.ConverterError
}

func (mfkc *MockFlipKickConverter) ToModel(flipKick flip_kick.FlipKickEntity) (flip_kick.FlipKickModel, error) {
	mfkc.EntitiesToConvert = append(mfkc.EntitiesToConvert, flipKick)
	return FlipKickModel, nil
}
func (mfkc *MockFlipKickConverter) SetConverterError(err error) {
	mfkc.ConverterError = err
}

type MockFlipKickRepository struct {
	HeaderIds           []int64
	HeadersToReturn     []core.Header
	StartingBlockNumber int64
	EndingBlockNumber   int64
	FlipKicksCreated    []flip_kick.FlipKickModel
	CreateRecordError   error
	MissingHeadersError error
}

func (mfkr *MockFlipKickRepository) Create(headerId int64, flipKick flip_kick.FlipKickModel) error {
	mfkr.HeaderIds = append(mfkr.HeaderIds, headerId)
	mfkr.FlipKicksCreated = append(mfkr.FlipKicksCreated, flipKick)

	return mfkr.CreateRecordError
}

func (mfkr *MockFlipKickRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	mfkr.StartingBlockNumber = startingBlockNumber
	mfkr.EndingBlockNumber = endingBlockNumber

	return mfkr.HeadersToReturn, mfkr.MissingHeadersError
}

func (mfkr *MockFlipKickRepository) SetHeadersToReturn(headers []core.Header) {
	mfkr.HeadersToReturn = headers
}

func (mfkr *MockFlipKickRepository) SetCreateRecordError(err error) {
	mfkr.CreateRecordError = err
}
func (mfkr *MockFlipKickRepository) SetMissingHeadersError(err error) {
	mfkr.MissingHeadersError = err
}
