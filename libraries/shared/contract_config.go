package shared

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type ContractConfig struct {
	Address    string
	Owner      string
	Abi        string
	ParsedAbi  abi.ABI
	FirstBlock int64
	LastBlock  int64
	Name       string
	Filters    []filters.LogFilter
}
