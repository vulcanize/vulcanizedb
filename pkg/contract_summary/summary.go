package contract_summary

import (
	"errors"
	"fmt"

	"math/big"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

type ContractSummary struct {
	ContractHash         string
	NumberOfTransactions int
	LastTransaction      *core.Transaction
	blockChain           core.Blockchain
	Attributes           core.ContractAttributes
	BlockNumber          *big.Int
	WatchedContract      core.WatchedContract
}

var NewContractNotWatchedErr = func(contractHash string) error {
	return errors.New(fmt.Sprintf("Contract %v not being watched", contractHash))
}

func NewSummary(blockchain core.Blockchain, repository repositories.Repository, contractHash string, blockNumber *big.Int) (ContractSummary, error) {
	watchedContract := repository.FindWatchedContract(contractHash)
	if watchedContract != nil {
		return newContractSummary(blockchain, *watchedContract, blockNumber), nil
	} else {
		return ContractSummary{}, NewContractNotWatchedErr(contractHash)
	}
}

func (contractSummary ContractSummary) GetStateAttribute(attributeName string) interface{} {
	var result interface{}
	result, _ = contractSummary.blockChain.GetAttribute(contractSummary.WatchedContract, attributeName, contractSummary.BlockNumber)
	return result
}

func newContractSummary(blockchain core.Blockchain, watchedContract core.WatchedContract, blockNumber *big.Int) ContractSummary {
	attributes, _ := blockchain.GetAttributes(watchedContract)
	return ContractSummary{
		blockChain:           blockchain,
		ContractHash:         watchedContract.Hash,
		NumberOfTransactions: len(watchedContract.Transactions),
		LastTransaction:      lastTransaction(watchedContract),
		Attributes:           attributes,
		BlockNumber:          blockNumber,
		WatchedContract:      watchedContract,
	}
}

func lastTransaction(watchedContract core.WatchedContract) *core.Transaction {
	if len(watchedContract.Transactions) > 0 {
		return &watchedContract.Transactions[0]
	} else {
		return nil
	}
}
