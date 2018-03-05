package inmemory

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
)

type ContractRepostiory struct {
	*InMemory
}

func (contractRepository *ContractRepostiory) ContractExists(contractHash string) bool {
	_, present := contractRepository.contracts[contractHash]
	return present
}

func (contractRepository *ContractRepostiory) GetContract(contractHash string) (core.Contract, error) {
	contract, ok := contractRepository.contracts[contractHash]
	if !ok {
		return core.Contract{}, datastore.ErrContractDoesNotExist(contractHash)
	}
	for _, block := range contractRepository.blocks {
		for _, transaction := range block.Transactions {
			if transaction.To == contractHash {
				contract.Transactions = append(contract.Transactions, transaction)
			}
		}
	}
	return contract, nil
}

func (contractRepository *ContractRepostiory) CreateContract(contract core.Contract) error {
	contractRepository.contracts[contract.Hash] = contract
	return nil
}
