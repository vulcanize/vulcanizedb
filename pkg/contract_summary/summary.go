package contract_summary

import (
	"errors"
	"fmt"

	"math/big"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

type ContractSummary struct {
	Attributes           core.ContractAttributes
	BlockNumber          *big.Int
	Contract             core.Contract
	ContractHash         string
	LastTransaction      *core.Transaction
	NumberOfTransactions int
	blockChain           core.Blockchain
}

var ErrContractDoesNotExist = func(contractHash string) error {
	return errors.New(fmt.Sprintf("Contract %v does not exist", contractHash))
}

func NewSummary(blockchain core.Blockchain, repository repositories.Repository, contractHash string, blockNumber *big.Int) (ContractSummary, error) {
	contract := repository.FindContract(contractHash)
	if contract != nil {
		return newContractSummary(blockchain, *contract, blockNumber), nil
	} else {
		return ContractSummary{}, ErrContractDoesNotExist(contractHash)
	}
}

func (contractSummary ContractSummary) GetStateAttribute(attributeName string) interface{} {
	var result interface{}
	result, _ = contractSummary.blockChain.GetAttribute(contractSummary.Contract, attributeName, contractSummary.BlockNumber)
	return result
}

func newContractSummary(blockchain core.Blockchain, contract core.Contract, blockNumber *big.Int) ContractSummary {
	attributes, _ := blockchain.GetAttributes(contract)
	return ContractSummary{
		Attributes:           attributes,
		BlockNumber:          blockNumber,
		Contract:             contract,
		ContractHash:         contract.Hash,
		LastTransaction:      lastTransaction(contract),
		NumberOfTransactions: len(contract.Transactions),
		blockChain:           blockchain,
	}
}

func lastTransaction(contract core.Contract) *core.Transaction {
	if len(contract.Transactions) > 0 {
		return &contract.Transactions[0]
	} else {
		return nil
	}
}
