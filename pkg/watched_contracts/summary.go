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
}

var NewContractNotWatchedErr = func(contractHash string) error {
	return errors.New(fmt.Sprintf("Contract %v not being watched", contractHash))
}

func NewSummary(_ core.Blockchain, repository repositories.Repository, contractHash string) (*ContractSummary, error) {
	contract := repository.FindWatchedContract(contractHash)
	if contract != nil {
		return newContractSummary(*contract), nil
	} else {
		return nil, NewContractNotWatchedErr(contractHash)
	}
}

func (ContractSummary) GetStateAttribute(attributeName string) string {
	return "Hello world"
}

func newContractSummary(contract core.WatchedContract) *ContractSummary {
	return &ContractSummary{
		ContractHash:         contract.Hash,
		NumberOfTransactions: len(contract.Transactions),
		LastTransaction:      lastTransaction(contract),
	}
}

func lastTransaction(contract core.WatchedContract) *core.Transaction {
	if len(contract.Transactions) > 0 {
		return &contract.Transactions[0]
	} else {
		return nil
	}
}
