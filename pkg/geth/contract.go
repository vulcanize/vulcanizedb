package geth

import (
	"errors"

	"sort"

	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var (
	ErrInvalidStateAttribute = errors.New("invalid state attribute")
)

func (blockchain *Blockchain) GetAttribute(contract core.Contract, attributeName string, blockNumber *big.Int) (interface{}, error) {
	parsed, err := ParseAbi(contract.Abi)
	var result interface{}
	if err != nil {
		return result, err
	}
	input, err := parsed.Pack(attributeName)
	if err != nil {
		return nil, ErrInvalidStateAttribute
	}
	output, err := blockchain.callContract(contract.Hash, input, blockNumber)
	if err != nil {
		return nil, err
	}
	err = parsed.Unpack(&result, attributeName, output)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (blockchain *Blockchain) FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error {
	parsed, err := ParseAbi(abiJSON)
	if err != nil {
		return err
	}
	input, err := parsed.Pack(method, methodArg)
	if err != nil {
		return err
	}
	output, err := blockchain.callContract(address, input, big.NewInt(blockNumber))
	if err != nil {
		return err
	}
	return parsed.Unpack(result, method, output)
}

func (blockchain *Blockchain) callContract(contractHash string, input []byte, blockNumber *big.Int) ([]byte, error) {
	to := common.HexToAddress(contractHash)
	msg := ethereum.CallMsg{To: &to, Data: input}
	return blockchain.client.CallContract(context.Background(), msg, blockNumber)
}

func (blockchain *Blockchain) GetAttributes(contract core.Contract) (core.ContractAttributes, error) {
	parsed, _ := ParseAbi(contract.Abi)
	var contractAttributes core.ContractAttributes
	for _, abiElement := range parsed.Methods {
		if (len(abiElement.Outputs) > 0) && (len(abiElement.Inputs) == 0) && abiElement.Const {
			attributeType := abiElement.Outputs[0].Type.String()
			contractAttributes = append(contractAttributes, core.ContractAttribute{Name: abiElement.Name, Type: attributeType})
		}
	}
	sort.Sort(contractAttributes)
	return contractAttributes, nil
}
