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
	abiString, err := ReadAbiFile(abiFilePath)
	if err != nil {
		return abi.ABI{}, ErrMissingAbiFile
	}
	return ParseAbi(abiString)
}

func ParseAbi(abiString string) (abi.ABI, error) {
	parsedAbi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return abi.ABI{}, ErrInvalidAbiFile
	}
	return parsedAbi, nil
}

func ReadAbiFile(abiFilePath string) (string, error) {
	filesBytes, err := ioutil.ReadFile(abiFilePath)
	if err != nil {
		return "", ErrMissingAbiFile
	}
	return string(filesBytes), nil
}
