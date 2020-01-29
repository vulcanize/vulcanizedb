package test_data

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/eth/datastore/postgres/repositories"
)

// Create a header sync log to reference in an event, returning inserted header sync log
func CreateTestLog(headerID int64, db *postgres.DB) core.HeaderSyncLog {
	log := types.Log{
		Address:     common.Address{},
		Topics:      nil,
		Data:        nil,
		BlockNumber: 0,
		TxHash:      common.Hash{},
		TxIndex:     uint(rand.Int31()),
		BlockHash:   common.Hash{},
		Index:       0,
		Removed:     false,
	}
	headerSyncLogRepository := repositories.NewHeaderSyncLogRepository(db)
	insertLogsErr := headerSyncLogRepository.CreateHeaderSyncLogs(headerID, []types.Log{log})
	Expect(insertLogsErr).NotTo(HaveOccurred())
	headerSyncLogs, getLogsErr := headerSyncLogRepository.GetUntransformedHeaderSyncLogs()
	Expect(getLogsErr).NotTo(HaveOccurred())
	for _, headerSyncLog := range headerSyncLogs {
		if headerSyncLog.Log.TxIndex == log.TxIndex {
			return headerSyncLog
		}
	}
	panic("couldn't find inserted test log")
}
