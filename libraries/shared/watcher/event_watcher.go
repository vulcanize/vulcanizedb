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
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/libraries/shared/chunker"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/logs"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transactions"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"time"
)

const NoNewDataPause = time.Second * 7

type EventWatcher struct {
	blockChain   core.BlockChain
	db           *postgres.DB
	LogDelegator logs.ILogDelegator
	LogExtractor logs.ILogExtractor
}

func NewEventWatcher(db *postgres.DB, bc core.BlockChain) EventWatcher {
	extractor := &logs.LogExtractor{
		Fetcher:                  fetcher.NewLogFetcher(bc),
		CheckedHeadersRepository: repositories.NewCheckedHeadersRepository(db),
		LogRepository:            repositories.NewHeaderSyncLogRepository(db),
		Syncer:                   transactions.NewTransactionsSyncer(db, bc),
	}
	logTransformer := &logs.LogDelegator{
		Chunker:       chunker.NewLogChunker(),
		LogRepository: repositories.NewHeaderSyncLogRepository(db),
	}
	return EventWatcher{
		blockChain:   bc,
		db:           db,
		LogExtractor: extractor,
		LogDelegator: logTransformer,
	}
}

// Adds transformers to the watcher so that their logs will be extracted and delegated.
func (watcher *EventWatcher) AddTransformers(initializers []transformer.EventTransformerInitializer) {
	for _, initializer := range initializers {
		t := initializer(watcher.db)

		watcher.LogDelegator.AddTransformer(t)
		watcher.LogExtractor.AddTransformerConfig(t.GetConfig())
	}
}

// Extracts and delegates watched log events.
func (watcher *EventWatcher) Execute(recheckHeaders constants.TransformerExecution, errsChan chan error) {
	extractErrsChan := make(chan error)
	delegateErrsChan := make(chan error)

	go watcher.extractLogs(recheckHeaders, extractErrsChan)
	go watcher.delegateLogs(delegateErrsChan)

	for {
		select {
		case extractErr := <-extractErrsChan:
			logrus.Errorf("error extracting logs in event watcher: %s", extractErr.Error())
			errsChan <- extractErr
		case delegateErr := <-delegateErrsChan:
			logrus.Errorf("error delegating logs in event watcher: %s", delegateErr.Error())
			errsChan <- delegateErr
		}
	}
}

func (watcher *EventWatcher) extractLogs(recheckHeaders constants.TransformerExecution, errs chan error) {
	err, missingHeadersFound := watcher.LogExtractor.ExtractLogs(recheckHeaders)
	if err != nil {
		errs <- err
	}

	if missingHeadersFound {
		watcher.extractLogs(recheckHeaders, errs)
	} else {
		time.Sleep(NoNewDataPause)
		watcher.extractLogs(recheckHeaders, errs)
	}
}

func (watcher *EventWatcher) delegateLogs(errs chan error) {
	err, logsFound := watcher.LogDelegator.DelegateLogs()
	if err != nil {
		errs <- err
	}

	if logsFound {
		watcher.delegateLogs(errs)
	} else {
		time.Sleep(NoNewDataPause)
		watcher.delegateLogs(errs)
	}
}
