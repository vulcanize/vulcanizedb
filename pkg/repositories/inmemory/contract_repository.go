package inmemory

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

func (repository *InMemory) ContractExists(contractHash string) bool {
	_, present := repository.contracts[contractHash]
	return present
}

func (repository *InMemory) GetContract(contractHash string) (core.Contract, error) {
	contract, ok := repository.contracts[contractHash]
	if !ok {
		return core.Contract{}, repositories.ErrContractDoesNotExist(contractHash)
	}
	for _, block := range repository.blocks {
		for _, transaction := range block.Transactions {
			if transaction.To == contractHash {
				contract.Transactions = append(contract.Transactions, transaction)
			}
		}
	}
	return contract, nil
}

func (repository *InMemory) CreateContract(contract core.Contract) error {
	repository.contracts[contract.Hash] = contract
	return nil
}
