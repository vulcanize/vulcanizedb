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

package watcher

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transactions"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/chunker"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type EventWatcher struct {
	Transformers  []transformer.EventTransformer
	BlockChain    core.BlockChain
	DB            *postgres.DB
	Fetcher       fetcher.ILogFetcher
	Chunker       chunker.Chunker
	Addresses     []common.Address
	Topics        []common.Hash
	StartingBlock *int64
	Syncer        transactions.ITransactionsSyncer
}

func NewEventWatcher(db *postgres.DB, bc core.BlockChain) EventWatcher {
	logChunker := chunker.NewLogChunker()
	logFetcher := fetcher.NewLogFetcher(bc)
	transactionSyncer := transactions.NewTransactionsSyncer(db, bc)
	return EventWatcher{
		BlockChain: bc,
		DB:         db,
		Fetcher:    logFetcher,
		Chunker:    logChunker,
		Syncer:     transactionSyncer,
	}
}

// Adds transformers to the watcher and updates the chunker, so that it will consider the new transformers.
func (watcher *EventWatcher) AddTransformers(initializers []transformer.EventTransformerInitializer) {
	var contractAddresses []common.Address
	var topic0s []common.Hash
	var configs []transformer.EventTransformerConfig

	for _, initializer := range initializers {
		t := initializer(watcher.DB)
		watcher.Transformers = append(watcher.Transformers, t)

		config := t.GetConfig()
		configs = append(configs, config)

		if watcher.StartingBlock == nil {
			watcher.StartingBlock = &config.StartingBlockNumber
		} else if earlierStartingBlockNumber(config.StartingBlockNumber, *watcher.StartingBlock) {
			watcher.StartingBlock = &config.StartingBlockNumber
		}

		addresses := transformer.HexStringsToAddresses(config.ContractAddresses)
		contractAddresses = append(contractAddresses, addresses...)
		topic0s = append(topic0s, common.HexToHash(config.Topic))
	}

	watcher.Addresses = append(watcher.Addresses, contractAddresses...)
	watcher.Topics = append(watcher.Topics, topic0s...)
	watcher.Chunker.AddConfigs(configs)
}

func (watcher *EventWatcher) Execute(recheckHeaders constants.TransformerExecution) error {
	if watcher.Transformers == nil {
		return fmt.Errorf("No transformers added to watcher")
	}

	checkedColumnNames, err := repository.GetCheckedColumnNames(watcher.DB)
	if err != nil {
		return err
	}
	notCheckedSQL := repository.CreateHeaderCheckedPredicateSQL(checkedColumnNames, recheckHeaders)

	missingHeaders, err := repository.MissingHeaders(*watcher.StartingBlock, -1, watcher.DB, notCheckedSQL)
	if err != nil {
		logrus.Error("Fetching of missing headers failed in watcher!")
		return err
	}

	for _, header := range missingHeaders {
		// TODO Extend FetchLogs for doing several blocks at a time
		logs, err := watcher.Fetcher.FetchLogs(watcher.Addresses, watcher.Topics, header)
		if err != nil {
			logrus.Errorf("Error while fetching logs for header %v in watcher", header.Id)
			return err
		}

		transactionsSyncErr := watcher.Syncer.SyncTransactions(header.Id, logs)
		if transactionsSyncErr != nil {
			logrus.Errorf("error syncing transactions: %s", transactionsSyncErr.Error())
			return transactionsSyncErr
		}

		transformErr := watcher.transformLogs(logs, header)
		if transformErr != nil {
			return transformErr
		}
	}
	return err
}

func (watcher *EventWatcher) transformLogs(logs []types.Log, header core.Header) error {
	chunkedLogs := watcher.Chunker.ChunkLogs(logs)

	// Can't quit early and mark as checked if there are no logs. If we are running continuousLogSync,
	// not all logs we're interested in might have been fetched.
	for _, t := range watcher.Transformers {
		transformerName := t.GetConfig().TransformerName
		logChunk := chunkedLogs[transformerName]
		err := t.Execute(logChunk, header)
		if err != nil {
			logrus.Errorf("%v transformer failed to execute in watcher: %v", transformerName, err)
			return err
		}
	}
	return nil
}

func earlierStartingBlockNumber(transformerBlock, watcherBlock int64) bool {
	return transformerBlock < watcherBlock
}
