package geth

import (
	"errors"
	"fmt"
	"path/filepath"

	"sort"

	"github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidStateAttribute = errors.New("invalid state attribute")
)

func (blockchain *GethBlockchain) GetContractAttributes(contractHash string) (core.ContractAttributes, error) {
	abiFilePath := filepath.Join(config.ProjectRoot(), "contracts", "public", fmt.Sprintf("%s.json", contractHash))
	parsed, _ := ParseAbiFile(abiFilePath)
	var contractAttributes core.ContractAttributes
	for _, abiElement := range parsed.Methods {
		if (len(abiElement.Outputs) > 0) && (len(abiElement.Inputs) == 0) && abiElement.Const {
			attributeType := abiElement.Outputs[0].Type.String()
			contractAttributes = append(contractAttributes, core.ContractAttribute{abiElement.Name, attributeType})
		}
	}
	sort.Sort(contractAttributes)
	return contractAttributes, nil
}

func convertresult(result interface{}, attributeName string) *string {
	var stringResult = new(string)
	fmt.Println(fmt.Sprintf("%s: %v (%T)", attributeName, result, result))
	switch t := result.(type) {
	case common.Address:
		ca := result.(common.Address)
		*stringResult = fmt.Sprintf("%v", ca.Hex())
	default:
		_ = t
		*stringResult = fmt.Sprintf("%v", result)
	}
	return stringResult

}

func (blockchain *GethBlockchain) GetContractStateAttribute(contractHash string, attributeName string) (*string, error) {
	boundContract, err := bindContract(common.HexToAddress(contractHash), blockchain.client, blockchain.client)
	if err != nil {
		return nil, err
	}
	var result interface{}
	err = boundContract.Call(&bind.CallOpts{}, &result, attributeName)
	if err != nil {
		return nil, ErrInvalidStateAttribute
	}
	return convertresult(result, attributeName), nil
}

func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	abiFilePath := filepath.Join(config.ProjectRoot(), "contracts", "public", fmt.Sprintf("%s.json", address.Hex()))
	parsed, err := ParseAbiFile(abiFilePath)
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}
