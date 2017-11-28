package geth

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	ErrInvalidAbiFile = errors.New("invalid abi")
	ErrMissingAbiFile = errors.New("missing abi")
)

func ParseAbiFile(abiFilePath string) (abi.ABI, error) {
	filesBytes, err := ioutil.ReadFile(abiFilePath)
	if err != nil {
		return abi.ABI{}, ErrMissingAbiFile
	}
	abiString := string(filesBytes)
	parsedAbi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return abi.ABI{}, ErrInvalidAbiFile
	}
	return parsedAbi, nil
}
