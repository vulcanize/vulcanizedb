package geth

import (
	"errors"

	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidStateAttribute = errors.New("invalid state attribute")
)

func (blockChain *BlockChain) FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error {
	parsed, err := ParseAbi(abiJSON)
	if err != nil {
		return err
	}
	var input []byte
	if methodArg != nil {
		input, err = parsed.Pack(method, methodArg)
	} else {
		input, err = parsed.Pack(method)
	}
	if err != nil {
		return err
	}
	output, err := blockChain.callContract(address, input, big.NewInt(blockNumber))
	if err != nil {
		return err
	}
	return parsed.Unpack(result, method, output)
}

func (blockChain *BlockChain) callContract(contractHash string, input []byte, blockNumber *big.Int) ([]byte, error) {
	to := common.HexToAddress(contractHash)
	msg := ethereum.CallMsg{To: &to, Data: input}
	return blockChain.ethClient.CallContract(context.Background(), msg, blockNumber)
}
