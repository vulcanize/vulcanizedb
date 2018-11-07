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

package every_block

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type ERC20Transformer struct {
	Getter     ERC20GetterInterface
	Repository ERC20TokenDatastore
	Retriever  erc20_watcher.TokenHolderRetriever
	Config     shared.ContractConfig
}

func (t *ERC20Transformer) SetConfiguration(config shared.ContractConfig) {
	t.Config = config
}

func NewERC20TokenTransformer(db *postgres.DB, blockchain core.BlockChain, con shared.ContractConfig) (shared.Transformer, error) {
	getter := NewGetter(blockchain)
	repository := ERC20TokenRepository{DB: db}
	retriever := erc20_watcher.NewTokenHolderRetriever(db, con.Address)

	transformer := ERC20Transformer{
		Getter:     &getter,
		Repository: &repository,
		Retriever:  retriever,
		Config:     con,
	}

	return transformer, nil
}

const (
	FetchingBlocksError         = "Error fetching missing blocks starting at block number %d: %s"
	FetchingSupplyError         = "Error fetching supply for block %d: %s"
	CreateSupplyError           = "Error inserting token_supply for block %d: %s"
	FetchingTokenAddressesError = "Error fetching token holder addresses at block %d: %s"
	FetchingBalanceError        = "Error fetching balance at block %d: %s"
	CreateBalanceError          = "Error inserting token_balance at block %d: %s"
	FetchingAllowanceError      = "Error fetching allowance at block %d: %s"
	CreateAllowanceError        = "Error inserting allowance at block %d: %s"
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

func (t ERC20Transformer) Execute() error {
	var upperBoundBlock int64
	blockchain := t.Getter.GetBlockChain()
	lastBlock := blockchain.LastBlock().Int64()

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
	log.Printf("Fetching totalSupply for %d blocks", len(blocks))

	// For each block missing total supply, create supply model and feed the missing data into the repository
	for _, blockNumber := range blocks {
		totalSupply, err := t.Getter.GetTotalSupply(t.Config.Abi, t.Config.Address, blockNumber)

		if err != nil {
			return newTransformerError(err, blockNumber, FetchingSupplyError)
		}
		// Create the supply model
		model := createTokenSupplyModel(totalSupply, t.Config.Address, blockNumber)
		// Feed it into the repository
		err = t.Repository.CreateSupply(model)

		if err != nil {
			return newTransformerError(err, blockNumber, CreateSupplyError)
		}
	}

	// Balance and allowance transformations:

	// Retrieve all token holder addresses for the given contract configuration

	tokenHolderAddresses, err := t.Retriever.RetrieveTokenHolderAddresses()
	if err != nil {
		return newTransformerError(err, t.Config.FirstBlock, FetchingTokenAddressesError)
	}

	// Iterate over the addresses and add their balances and allowances at each block height to the repository
	for holderAddr := range tokenHolderAddresses {

		// Balance transformations:

		blocks, err := t.Repository.MissingBalanceBlocks(t.Config.FirstBlock, upperBoundBlock, t.Config.Address, holderAddr.String())

		if err != nil {
			return newTransformerError(err, t.Config.FirstBlock, FetchingBlocksError)
		}

		log.Printf("Fetching balances for %d blocks", len(blocks))

		// For each block missing balances for the given address, create a balance model and feed the missing data into the repository
		for _, blockNumber := range blocks {

			hashArgs := []common.Address{holderAddr}
			balanceOfArgs := make([]interface{}, len(hashArgs))
			for i, s := range hashArgs {
				balanceOfArgs[i] = s
			}

			totalSupply, err := t.Getter.GetBalance(t.Config.Abi, t.Config.Address, blockNumber, balanceOfArgs)

			if err != nil {
				return newTransformerError(err, blockNumber, FetchingBalanceError)
			}

			model := createTokenBalanceModel(totalSupply, t.Config.Address, blockNumber, holderAddr.String())

			err = t.Repository.CreateBalance(model)

			if err != nil {
				return newTransformerError(err, blockNumber, CreateBalanceError)
			}
		}

		// Allowance transformations:

		for spenderAddr := range tokenHolderAddresses {

			blocks, err := t.Repository.MissingAllowanceBlocks(t.Config.FirstBlock, upperBoundBlock, t.Config.Address, holderAddr.String(), spenderAddr.String())

			if err != nil {
				return newTransformerError(err, t.Config.FirstBlock, FetchingBlocksError)
			}

			log.Printf("Fetching allowances for %d blocks", len(blocks))

			// For each block missing allowances for the given holder and spender addresses, create a allowance model and feed the missing data into the repository
			for _, blockNumber := range blocks {

				hashArgs := []common.Address{holderAddr, spenderAddr}
				allowanceArgs := make([]interface{}, len(hashArgs))
				for i, s := range hashArgs {
					allowanceArgs[i] = s
				}

				totalSupply, err := t.Getter.GetAllowance(t.Config.Abi, t.Config.Address, blockNumber, allowanceArgs)

				if err != nil {
					return newTransformerError(err, blockNumber, FetchingAllowanceError)
				}

				model := createTokenAllowanceModel(totalSupply, t.Config.Address, blockNumber, holderAddr.String(), spenderAddr.String())

				err = t.Repository.CreateAllowance(model)

				if err != nil {
					return newTransformerError(err, blockNumber, CreateAllowanceError)
				}

			}

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

func createTokenBalanceModel(tokenBalance big.Int, tokenAddress string, blockNumber int64, tokenHolderAddress string) TokenBalance {
	return TokenBalance{
		Value:              tokenBalance.String(),
		TokenAddress:       tokenAddress,
		BlockNumber:        blockNumber,
		TokenHolderAddress: tokenHolderAddress,
	}
}

func createTokenAllowanceModel(tokenBalance big.Int, tokenAddress string, blockNumber int64, tokenHolderAddress, tokenSpenderAddress string) TokenAllowance {
	return TokenAllowance{
		Value:               tokenBalance.String(),
		TokenAddress:        tokenAddress,
		BlockNumber:         blockNumber,
		TokenHolderAddress:  tokenHolderAddress,
		TokenSpenderAddress: tokenSpenderAddress,
	}
}
