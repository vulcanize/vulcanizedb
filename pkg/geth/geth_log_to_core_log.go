package geth

import (
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func GethLogToCoreLog(gethLog types.Log) core.Log {
	topics := gethLog.Topics
	var hexTopics []string
	for _, topic := range topics {
		hexTopics = append(hexTopics, topic.Hex())
	}
	return core.Log{
		Address: gethLog.Address.Hex(),

		BlockNumber: int64(gethLog.BlockNumber),
		Topics:      hexTopics,
		TxHash:      gethLog.TxHash.Hex(),
		Data:        hexutil.Encode(gethLog.Data),
	}
}

func GethLogsToCoreLogs(gethLogs []types.Log) []core.Log {
	var logs []core.Log
	for _, log := range gethLogs {
		log := GethLogToCoreLog(log)
		logs = append(logs, log)
	}
	return logs
}
