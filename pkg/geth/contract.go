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
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidStateAttribute = errors.New("invalid state attribute")
)

func (blockchain *GethBlockchain) GetAttribute(watchedContract core.WatchedContract, attributeName string, blockNumber *big.Int) (interface{}, error) {
	parsed, err := ParseAbi(watchedContract.Abi)
	var result interface{}
	if err != nil {
		return result, err
	}
	input, err := parsed.Pack(attributeName)
	if err != nil {
		return nil, ErrInvalidStateAttribute
	}
	output, err := callContract(watchedContract.Hash, input, blockchain, blockNumber)
	if err != nil {
		return nil, err
	}
	err = parsed.Unpack(&result, attributeName, output)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func callContract(contractHash string, input []byte, blockchain *GethBlockchain, blockNumber *big.Int) ([]byte, error) {
	to := common.HexToAddress(contractHash)
	msg := ethereum.CallMsg{To: &to, Data: input}
	return blockchain.client.CallContract(context.Background(), msg, blockNumber)
}

func (blockchain *GethBlockchain) GetAttributes(watchedContract core.WatchedContract) (core.ContractAttributes, error) {
	parsed, _ := ParseAbi(watchedContract.Abi)
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

func (blockchain *GethBlockchain) GetContractAttributesOld(contractHash string) (core.ContractAttributes, error) {
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
