package inmemory

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

func (repository *InMemory) CreateLogs(logs []core.Log) error {
	for _, log := range logs {
		key := fmt.Sprintf("%d%d", log.BlockNumber, log.Index)
		var logs []core.Log
		repository.logs[key] = append(logs, log)
	}
	return nil
}

func (repository *InMemory) GetLogs(address string, blockNumber int64) []core.Log {
	var matchingLogs []core.Log
	for _, logs := range repository.logs {
		for _, log := range logs {
			if log.Address == address && log.BlockNumber == blockNumber {
				matchingLogs = append(matchingLogs, log)
			}
		}
	}
	return matchingLogs
}
