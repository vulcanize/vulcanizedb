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

package every_block

import (
	"fmt"
	"log"
	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type ERC20FetcherInterface interface {
	FetchSupplyOf(contractAbi string, contractAddress string, blockNumber int64) (big.Int, error)
	GetBlockChain() core.BlockChain
}

func NewFetcher(blockchain core.BlockChain) Fetcher {
	return Fetcher{
		Blockchain: blockchain,
	}
}

type Fetcher struct {
	Blockchain      core.BlockChain
	ContractAbi     string
	ContractAddress string
}

type fetcherError struct {
	err         string
	fetchMethod string
}

func (fe *fetcherError) Error() string {
	return fmt.Sprintf("Error fetching %s: %s", fe.fetchMethod, fe.err)
}

func newFetcherError(err error, fetchMethod string) *fetcherError {
	e := fetcherError{err.Error(), fetchMethod}
	log.Println(e.Error())
	return &e
}

func (f Fetcher) FetchSupplyOf(contractAbi string, contractAddress string, blockNumber int64) (big.Int, error) {
	method := "totalSupply"
	var result = new(big.Int)
	err := f.Blockchain.FetchContractData(contractAbi, contractAddress, method, nil, &result, blockNumber)

	if err != nil {
		return *result, newFetcherError(err, method)
	}

	return *result, nil
}

func (f Fetcher) GetBlockChain() core.BlockChain {
	return f.Blockchain
}
