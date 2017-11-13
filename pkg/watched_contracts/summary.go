package watched_contracts

import (
	"errors"
	"fmt"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

type ContractSummary struct {
	ContractHash string
}

var NewContractNotWatchedErr = func(contractHash string) error {
	return errors.New(fmt.Sprintf("Contract %v not being watched", contractHash))
}

func NewSummary(repository repositories.Repository, contractHash string) (*ContractSummary, error) {
	contract := repository.FindWatchedContract(contractHash)
	if contract != nil {
		return newContractSummary(*contract), nil
	} else {
		return nil, NewContractNotWatchedErr(contractHash)
	}
}

func newContractSummary(contract core.WatchedContract) *ContractSummary {
	return &ContractSummary{ContractHash: contract.Hash}
}
