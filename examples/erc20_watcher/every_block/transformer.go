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
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"log"
	"math/big"
)

type Transformer struct {
	Getter     ERC20GetterInterface
	Repository ERC20RepositoryInterface
	Config     erc20_watcher.ContractConfig
}

func (t *Transformer) SetConfiguration(config erc20_watcher.ContractConfig) {
	t.Config = config
}

type TokenSupplyTransformerInitializer struct {
	Config erc20_watcher.ContractConfig
}

func (i TokenSupplyTransformerInitializer) NewTokenSupplyTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	getter := NewGetter(blockChain)
	repository := ERC20TokenRepository{DB: db}
	transformer := Transformer{
		Getter:     &getter,
		Repository: &repository,
		Config:     i.Config,
	}

	return transformer
}

const (
	FetchingBlocksError = "Error getting missing blocks starting at block number %d: %s"
	GetSupplyError      = "Error getting supply for block %d: %s"
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
	blockChain := t.Getter.GetBlockChain()
	lastBlock := blockChain.LastBlock().Int64()

	if t.Config.LastBlock == -1 {
		upperBoundBlock = lastBlock
	} else {
		upperBoundBlock = t.Config.LastBlock
	}

	// Supply transformations:

	// Fetch missing supply blocks
	blocks, err := t.Repository.MissingSupplyBlocks(t.Config.FirstBlock, upperBoundBlock, t.Config.Address)

	if err != nil {
		return newTransformerError(err, t.Config.FirstBlock, FetchingBlocksError)
	}

	// Fetch supply for missing blocks
	log.Printf("Gets totalSupply for %d blocks", len(blocks))

	// For each block missing total supply, create supply model and feed the missing data into the repository
	for _, blockNumber := range blocks {
		totalSupply, err := t.Getter.GetTotalSupply(t.Config.Abi, t.Config.Address, blockNumber)

		if err != nil {
			return newTransformerError(err, blockNumber, GetSupplyError)
		}
		// Create the supply model
		model := createTokenSupplyModel(totalSupply, t.Config.Address, blockNumber)
		// Feed it into the repository
		err = t.Repository.CreateSupply(model)

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
