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

package fetcher

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

// Fetcher serves as the lower level data fetcher that calls the underlying
// blockchain's FetchConctractData method for a given return type

// Interface definition for a Fetcher
type FetcherInterface interface {
	FetchBigInt(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (big.Int, error)
	FetchBool(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (bool, error)
	FetchAddress(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (common.Address, error)
	FetchString(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (string, error)
	FetchHash(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (common.Hash, error)
}

// Used to create a new Fetcher error for a given error and fetch method
func newFetcherError(err error, fetchMethod string) *fetcherError {
	e := fetcherError{err.Error(), fetchMethod}
	log.Println(e.Error())
	return &e
}

// Fetcher struct
type Fetcher struct {
	BlockChain core.BlockChain // Underyling Blockchain
}

// Fetcher error
type fetcherError struct {
	err         string
	fetchMethod string
}

// Fetcher error method
func (fe *fetcherError) Error() string {
	return fmt.Sprintf("Error fetching %s: %s", fe.fetchMethod, fe.err)
}

//  Generic Fetcher methods used by Getters to call contract methods

// Method used to fetch big.Int value from contract
func (f Fetcher) FetchBigInt(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (big.Int, error) {
	var result = new(big.Int)
	err := f.BlockChain.FetchContractData(contractAbi, contractAddress, method, methodArgs, &result, blockNumber)

	if err != nil {
		return *result, newFetcherError(err, method)
	}

	return *result, nil
}

// Method used to fetch bool value from contract
func (f Fetcher) FetchBool(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (bool, error) {
	var result = new(bool)
	err := f.BlockChain.FetchContractData(contractAbi, contractAddress, method, methodArgs, &result, blockNumber)

	if err != nil {
		return *result, newFetcherError(err, method)
	}

	return *result, nil
}

// Method used to fetch address value from contract
func (f Fetcher) FetchAddress(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (common.Address, error) {
	var result = new(common.Address)
	err := f.BlockChain.FetchContractData(contractAbi, contractAddress, method, methodArgs, &result, blockNumber)

	if err != nil {
		return *result, newFetcherError(err, method)
	}

	return *result, nil
}

// Method used to fetch string value from contract
func (f Fetcher) FetchString(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (string, error) {
	var result = new(string)
	err := f.BlockChain.FetchContractData(contractAbi, contractAddress, method, methodArgs, &result, blockNumber)

	if err != nil {
		return *result, newFetcherError(err, method)
	}

	return *result, nil
}

// Method used to fetch hash value from contract
func (f Fetcher) FetchHash(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (common.Hash, error) {
	var result = new(common.Hash)
	err := f.BlockChain.FetchContractData(contractAbi, contractAddress, method, methodArgs, &result, blockNumber)

	if err != nil {
		return *result, newFetcherError(err, method)
	}

	return *result, nil
}
