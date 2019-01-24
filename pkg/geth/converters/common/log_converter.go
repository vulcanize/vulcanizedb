// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

func ToCoreLogs(gethLogs []types.Log) []core.Log {
	var logs []core.Log
	for _, log := range gethLogs {
		log := ToCoreLog(log)
		logs = append(logs, log)
	}
	return logs
}

func makeTopics(topics []common.Hash) core.Topics {
	var hexTopics core.Topics
	for i, topic := range topics {
		hexTopics[i] = topic.Hex()
	}
	return hexTopics
}

func ToCoreLog(gethLog types.Log) core.Log {
	topics := gethLog.Topics
	hexTopics := makeTopics(topics)
	return core.Log{
		Address:     strings.ToLower(gethLog.Address.Hex()),
		BlockNumber: int64(gethLog.BlockNumber),
		Topics:      hexTopics,
		TxHash:      gethLog.TxHash.Hex(),
		Index:       int64(gethLog.Index),
		Data:        hexutil.Encode(gethLog.Data),
	}
}
