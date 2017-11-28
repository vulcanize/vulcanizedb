package watched_contracts

import (
	"errors"
	"fmt"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

type ContractSummary struct {
	ContractHash         string
	NumberOfTransactions int
	LastTransaction      *core.Transaction
	blockChain           core.Blockchain
	Attributes           core.ContractAttributes
}

var NewContractNotWatchedErr = func(contractHash string) error {
	return errors.New(fmt.Sprintf("Contract %v not being watched", contractHash))
}

func NewSummary(blockchain core.Blockchain, repository repositories.Repository, contractHash string) (*ContractSummary, error) {
	watchedContract := repository.FindWatchedContract(contractHash)
	if watchedContract != nil {
		return newContractSummary(blockchain, *watchedContract), nil
	} else {
		return nil, NewContractNotWatchedErr(contractHash)
	}
}

func (contractSummary ContractSummary) GetStateAttribute(attributeName string) string {
	result, _ := contractSummary.blockChain.GetContractStateAttribute(contractSummary.ContractHash, attributeName)
	return *result
}

func newContractSummary(blockchain core.Blockchain, watchedContract core.WatchedContract) *ContractSummary {
	attributes, _ := blockchain.GetContractAttributes(watchedContract.Hash)
	return &ContractSummary{
		blockChain:           blockchain,
		ContractHash:         watchedContract.Hash,
		NumberOfTransactions: len(watchedContract.Transactions),
		LastTransaction:      lastTransaction(watchedContract),
		Attributes:           attributes,
	}
}

func lastTransaction(contract core.WatchedContract) *core.Transaction {
	if len(contract.Transactions) > 0 {
		return &contract.Transactions[0]
	} else {
		return nil
	}
}
