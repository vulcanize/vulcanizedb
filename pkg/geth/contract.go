package geth

import (
	"errors"
	"fmt"
	"path/filepath"

	"sort"

	"context"
	"math/big"

	"github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidStateAttribute = errors.New("invalid state attribute")
)

func (blockchain *GethBlockchain) GetContract(contractHash string) (core.Contract, error) {
	attributes, err := blockchain.getContractAttributes(contractHash)
	if err != nil {
		return core.Contract{}, err
	} else {
		contract := core.Contract{
			Attributes: attributes,
			Hash:       contractHash,
		}
		return contract, nil
	}
}

func (blockchain *GethBlockchain) getParseAbi(contract core.Contract) (abi.ABI, error) {
	abiFilePath := filepath.Join(config.ProjectRoot(), "contracts", "public", fmt.Sprintf("%s.json", contract.Hash))
	parsed, err := ParseAbiFile(abiFilePath)
	if err != nil {
		return abi.ABI{}, err
	}
	return parsed, nil
}

func (blockchain *GethBlockchain) GetAttribute(contract core.Contract, attributeName string, blockNumber *big.Int) (interface{}, error) {
	parsed, err := blockchain.getParseAbi(contract)
	var result interface{}
	if err != nil {
		return result, err
	}
	input, err := parsed.Pack(attributeName)
	if err != nil {
		return nil, ErrInvalidStateAttribute
	}
	output, err := callContract(contract, input, err, blockchain, blockNumber)
	if err != nil {
		return nil, err
	}
	err = parsed.Unpack(&result, attributeName, output)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func callContract(contract core.Contract, input []byte, err error, blockchain *GethBlockchain, blockNumber *big.Int) ([]byte, error) {
	to := common.HexToAddress(contract.Hash)
	msg := ethereum.CallMsg{To: &to, Data: input}
	return blockchain.client.CallContract(context.Background(), msg, blockNumber)
}

func (blockchain *GethBlockchain) getContractAttributes(contractHash string) (core.ContractAttributes, error) {
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
