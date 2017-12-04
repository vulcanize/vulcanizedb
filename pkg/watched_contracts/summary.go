package watched_contracts

import (
	"errors"
	"fmt"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

type ContractSummary struct {
	Contract             core.Contract
	ContractHash         string
	NumberOfTransactions int
	LastTransaction      *core.Transaction
	blockChain           core.Blockchain
	Attributes           core.ContractAttributes
}

var NewContractNotWatchedErr = func(contractHash string) error {
	return errors.New(fmt.Sprintf("Contract %v not being watched", contractHash))
}

func NewSummary(blockchain core.Blockchain, repository repositories.Repository, contractHash string) (ContractSummary, error) {
	watchedContract := repository.FindWatchedContract(contractHash)
	if watchedContract != nil {
		return newContractSummary(blockchain, *watchedContract), nil
	} else {
		return ContractSummary{}, NewContractNotWatchedErr(contractHash)
	}
}

func (contractSummary ContractSummary) GetStateAttribute(attributeName string) interface{} {
	result, _ := contractSummary.blockChain.GetAttribute(contractSummary.Contract, attributeName)
	return result
}

func newContractSummary(blockchain core.Blockchain, watchedContract core.WatchedContract) ContractSummary {
	contract, _ := blockchain.GetContract(watchedContract.Hash)
	return ContractSummary{
		blockChain:           blockchain,
		Contract:             contract,
		ContractHash:         watchedContract.Hash,
		NumberOfTransactions: len(watchedContract.Transactions),
		LastTransaction:      lastTransaction(watchedContract),
		Attributes:           contract.Attributes,
	}
}

func lastTransaction(contract core.WatchedContract) *core.Transaction {
	if len(contract.Transactions) > 0 {
		return &contract.Transactions[0]
	} else {
		return nil
	}
}
