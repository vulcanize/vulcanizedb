package datastore

import "fmt"

func ErrBlockDoesNotExist(blockNumber int64) error {
	return fmt.Errorf("Block number %d does not exist", blockNumber)
}

func ErrContractDoesNotExist(contractHash string) error {
	return fmt.Errorf("Contract %v does not exist", contractHash)
}

func ErrFilterDoesNotExist(name string) error {
	return fmt.Errorf("filter %s does not exist", name)
}

func ErrReceiptDoesNotExist(txHash string) error {
	return fmt.Errorf("Receipt for tx: %v does not exist", txHash)
}
