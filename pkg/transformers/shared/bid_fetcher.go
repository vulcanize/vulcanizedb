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

package shared

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type IBidFetcher interface {
	FetchBid(contractAbi, contractAddress string, blockNumber int64, methodArgs interface{}) (Bid, error)
}

type BidFetcher struct {
	blockChain core.BlockChain
}

func NewBidFetcher(blockchain core.BlockChain) IBidFetcher {
	return BidFetcher{
		blockChain: blockchain,
	}
}

func (fetcher BidFetcher) FetchBid(contractAbi, contractAddress string, blockNumber int64, methodArgs interface{}) (Bid, error) {
	method := "bids"
	result := Bid{}
	err := fetcher.blockChain.FetchContractData(contractAbi, contractAddress, method, methodArgs, &result, blockNumber)

	if err != nil {
		return result, err
	}

	return result, nil
}
