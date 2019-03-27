package datastore

import "fmt"

var ErrBlockDoesNotExist = func(blockNumber int64) error {
	return fmt.Errorf("Block number %d does not exist", blockNumber)
}

var ErrContractDoesNotExist = func(contractHash string) error {
	return fmt.Errorf("Contract %v does not exist", contractHash)
}

var ErrFilterDoesNotExist = func(name string) error {
	return fmt.Errorf("filter %s does not exist", name)
}

var ErrReceiptDoesNotExist = func(txHash string) error {
	return fmt.Errorf("Receipt for tx: %v does not exist", txHash)
}
