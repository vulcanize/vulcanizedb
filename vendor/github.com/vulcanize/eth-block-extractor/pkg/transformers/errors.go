package transformers

import (
	"errors"
	"fmt"
)

const (
	GetBlockRlpErr = "Error fetching block RLP data"
	PutIpldErr     = "Error writing to IPFS"
)

var ErrInvalidRange = errors.New("ending block number must be greater than or equal to starting block number")

type ExecuteError struct {
	msg string
	err error
}

func NewExecuteError(msg string, err error) *ExecuteError {
	return &ExecuteError{msg: msg, err: err}
}

func (ee ExecuteError) Error() string {
	return fmt.Sprintf("%s: %s", ee.msg, ee.err.Error())
}
