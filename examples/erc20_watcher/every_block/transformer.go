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

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type Transformer struct {
	Fetcher    ERC20FetcherInterface
	Repository ERC20RepositoryInterface
	Config     erc20_watcher.ContractConfig
}

func (t *Transformer) SetConfiguration(config erc20_watcher.ContractConfig) {
	t.Config = config
}

type TokenSupplyTransformerInitializer struct {
	Config erc20_watcher.ContractConfig
}

func (i TokenSupplyTransformerInitializer) NewTokenSupplyTransformer(db *postgres.DB, blockchain core.BlockChain) shared.Transformer {
	fetcher := NewFetcher(blockchain)
	repository := TokenSupplyRepository{DB: db}
	transformer := Transformer{
		Fetcher:    &fetcher,
		Repository: &repository,
		Config:     i.Config,
	}

	return transformer
}

const (
	FetchingBlocksError = "Error fetching missing blocks starting at block number %d: %s"
	FetchingSupplyError = "Error fetching supply for block %d: %s"
	CreateSupplyError   = "Error inserting token_supply for block %d: %s"
)

type transformerError struct {
	err         string
	blockNumber int64
	msg         string
}

func (te *transformerError) Error() string {
	return fmt.Sprintf(te.msg, te.blockNumber, te.err)
}

func newTransformerError(err error, blockNumber int64, msg string) error {
	e := transformerError{err.Error(), blockNumber, msg}
	log.Println(e.Error())
	return &e
}

func (t Transformer) Execute() error {
	var upperBoundBlock int64
	blockchain := t.Fetcher.GetBlockChain()
	lastBlock := blockchain.LastBlock().Int64()

	if t.Config.LastBlock == -1 {
		upperBoundBlock = lastBlock
	} else {
		upperBoundBlock = t.Config.LastBlock
	}

	blocks, err := t.Repository.MissingBlocks(t.Config.FirstBlock, upperBoundBlock)

	if err != nil {
		return newTransformerError(err, t.Config.FirstBlock, FetchingBlocksError)
	}

	log.Printf("Fetching totalSupply for %d blocks", len(blocks))
	for _, blockNumber := range blocks {
		totalSupply, err := t.Fetcher.FetchSupplyOf(t.Config.Abi, t.Config.Address, blockNumber)

		if err != nil {
			return newTransformerError(err, blockNumber, FetchingSupplyError)
		}
		model := createTokenSupplyModel(totalSupply, t.Config.Address, blockNumber)
		err = t.Repository.Create(model)

		if err != nil {
			return newTransformerError(err, blockNumber, CreateSupplyError)
		}
	}

	return nil
}

func createTokenSupplyModel(totalSupply big.Int, address string, blockNumber int64) TokenSupply {
	return TokenSupply{
		Value:        totalSupply.String(),
		TokenAddress: address,
		BlockNumber:  blockNumber,
	}
}
