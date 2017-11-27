package geth

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidStateAttribute = errors.New("invalid state attribute")
)

func (blockchain *GethBlockchain) GetContractStateAttribute(contractHash string, attributeName string) (*string, error) {
	boundContract, err := bindContract(common.HexToAddress(contractHash), blockchain.client, blockchain.client)
	if err != nil {
		return nil, err
	}
	result := new(string)
	err = boundContract.Call(&bind.CallOpts{}, result, attributeName)
	if err != nil {
		return nil, ErrInvalidStateAttribute
	}
	return result, nil
}

func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	abiFilePath := filepath.Join(config.ProjectRoot(), "contracts", "public", fmt.Sprintf("%s.json", address.Hex()))
	parsed, err := ParseAbiFile(abiFilePath)
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}
