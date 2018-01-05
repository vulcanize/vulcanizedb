package geth

import (
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func LogToCoreLog(gethLog types.Log) core.Log {
	topics := gethLog.Topics
	var hexTopics = make(map[int]string)
	for i, topic := range topics {
		hexTopics[i] = topic.Hex()
	}
	return core.Log{
		Address: gethLog.Address.Hex(),

		BlockNumber: int64(gethLog.BlockNumber),
		Topics:      hexTopics,
		TxHash:      gethLog.TxHash.Hex(),
		Index:       int64(gethLog.Index),
		Data:        hexutil.Encode(gethLog.Data),
	}
}

func GethLogsToCoreLogs(gethLogs []types.Log) []core.Log {
	var logs []core.Log
	for _, log := range gethLogs {
		log := LogToCoreLog(log)
		logs = append(logs, log)
	}
	return logs
}
