// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package logs

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transactions"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
)

var (
	ErrNoUncheckedHeaders = errors.New("no unchecked headers available for log fetching")
	ErrNoWatchedAddresses = errors.New("no watched addresses configured in the log extractor")
)

type ILogExtractor interface {
	AddTransformerConfig(config transformer.EventTransformerConfig) error
	ExtractLogs(recheckHeaders constants.TransformerExecution) error
}

type LogExtractor struct {
	Addresses                []common.Address
	CheckedHeadersRepository datastore.CheckedHeadersRepository
	CheckedLogsRepository    datastore.CheckedLogsRepository
	Fetcher                  fetcher.ILogFetcher
	LogRepository            datastore.HeaderSyncLogRepository
	StartingBlock            *int64
	Syncer                   transactions.ITransactionsSyncer
	Topics                   []common.Hash
}

// Add additional logs to extract
func (extractor *LogExtractor) AddTransformerConfig(config transformer.EventTransformerConfig) error {
	checkedHeadersErr := extractor.updateCheckedHeaders(config)
	if checkedHeadersErr != nil {
		return checkedHeadersErr
	}

	if extractor.StartingBlock == nil {
		extractor.StartingBlock = &config.StartingBlockNumber
	} else if earlierStartingBlockNumber(config.StartingBlockNumber, *extractor.StartingBlock) {
		extractor.StartingBlock = &config.StartingBlockNumber
	}

	addresses := transformer.HexStringsToAddresses(config.ContractAddresses)
	extractor.Addresses = append(extractor.Addresses, addresses...)
	extractor.Topics = append(extractor.Topics, common.HexToHash(config.Topic))
	return nil
}

// Fetch and persist watched logs
func (extractor LogExtractor) ExtractLogs(recheckHeaders constants.TransformerExecution) error {
	if len(extractor.Addresses) < 1 {
		logrus.Errorf("error extracting logs: %s", ErrNoWatchedAddresses.Error())
		return ErrNoWatchedAddresses
	}

	uncheckedHeaders, uncheckedHeadersErr := extractor.CheckedHeadersRepository.UncheckedHeaders(*extractor.StartingBlock, -1, getCheckCount(recheckHeaders))
	if uncheckedHeadersErr != nil {
		logrus.Errorf("error fetching missing headers: %s", uncheckedHeadersErr)
		return uncheckedHeadersErr
	}

	if len(uncheckedHeaders) < 1 {
		return ErrNoUncheckedHeaders
	}

	for _, header := range uncheckedHeaders {
		logs, fetchLogsErr := extractor.Fetcher.FetchLogs(extractor.Addresses, extractor.Topics, header)
		if fetchLogsErr != nil {
			logError("error fetching logs for header: %s", fetchLogsErr, header)
			return fetchLogsErr
		}

		if len(logs) > 0 {
			transactionsSyncErr := extractor.Syncer.SyncTransactions(header.Id, logs)
			if transactionsSyncErr != nil {
				logError("error syncing transactions: %s", transactionsSyncErr, header)
				return transactionsSyncErr
			}

			createLogsErr := extractor.LogRepository.CreateHeaderSyncLogs(header.Id, logs)
			if createLogsErr != nil {
				logError("error persisting logs: %s", createLogsErr, header)
				return createLogsErr
			}
		}

		markHeaderCheckedErr := extractor.CheckedHeadersRepository.MarkHeaderChecked(header.Id)
		if markHeaderCheckedErr != nil {
			logError("error marking header checked: %s", markHeaderCheckedErr, header)
			return markHeaderCheckedErr
		}
	}
	return nil
}

func earlierStartingBlockNumber(transformerBlock, watcherBlock int64) bool {
	return transformerBlock < watcherBlock
}

func logError(description string, err error, header core.Header) {
	logrus.WithFields(logrus.Fields{
		"headerId":    header.Id,
		"headerHash":  header.Hash,
		"blockNumber": header.BlockNumber,
	}).Errorf(description, err.Error())
}

func getCheckCount(recheckHeaders constants.TransformerExecution) int64 {
	if recheckHeaders == constants.HeaderUnchecked {
		return 1
	} else {
		return constants.RecheckHeaderCap
	}
}

func (extractor *LogExtractor) updateCheckedHeaders(config transformer.EventTransformerConfig) error {
	alreadyWatchingLog, watchingLogErr := extractor.CheckedLogsRepository.AlreadyWatchingLog(config.ContractAddresses, config.Topic)
	if watchingLogErr != nil {
		return watchingLogErr
	}
	if !alreadyWatchingLog {
		uncheckHeadersErr := extractor.CheckedHeadersRepository.MarkHeadersUnchecked(config.StartingBlockNumber)
		if uncheckHeadersErr != nil {
			return uncheckHeadersErr
		}
		markLogWatchedErr := extractor.CheckedLogsRepository.MarkLogWatched(config.ContractAddresses, config.Topic)
		if markLogWatchedErr != nil {
			return markLogWatchedErr
		}
	}
	return nil
}
