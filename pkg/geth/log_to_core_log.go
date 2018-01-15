package geth

import (
	"strings"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func LogsToCoreLogs(gethLogs []types.Log) []core.Log {
	var logs []core.Log
	for _, log := range gethLogs {
		log := LogToCoreLog(log)
		logs = append(logs, log)
	}
	return logs
}

func MakeTopics(topics []common.Hash) core.Topics {
	var hexTopics core.Topics
	for i, topic := range topics {
		hexTopics[i] = topic.Hex()
	}
	return hexTopics
}

func LogToCoreLog(gethLog types.Log) core.Log {
	topics := gethLog.Topics
	hexTopics := MakeTopics(topics)
	return core.Log{
		Address: strings.ToLower(gethLog.Address.Hex()),

		BlockNumber: int64(gethLog.BlockNumber),
		Topics:      hexTopics,
		TxHash:      gethLog.TxHash.Hex(),
		Index:       int64(gethLog.Index),
		Data:        hexutil.Encode(gethLog.Data),
		Removed:     gethLog.Removed,
	}
}
