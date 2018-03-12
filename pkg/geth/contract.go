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
