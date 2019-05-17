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
	"time"

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
		CheckedHeadersRepository: repositories.NewCheckedHeadersRepository(db),
		CheckedLogsRepository:    repositories.NewCheckedLogsRepository(db),
		Fetcher:                  fetcher.NewLogFetcher(bc),
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
func (watcher *EventWatcher) AddTransformers(initializers []transformer.EventTransformerInitializer) error {
	for _, initializer := range initializers {
		t := initializer(watcher.db)

		watcher.LogDelegator.AddTransformer(t)
		err := watcher.LogExtractor.AddTransformerConfig(t.GetConfig())
		if err != nil {
			return err
		}
	}
	return nil
}

// Extracts and delegates watched log events.
func (watcher *EventWatcher) Execute(recheckHeaders constants.TransformerExecution) error {
	delegateErrsChan := make(chan error)
	extractErrsChan := make(chan error)
	defer close(delegateErrsChan)
	defer close(extractErrsChan)

	go watcher.extractLogs(recheckHeaders, extractErrsChan)
	go watcher.delegateLogs(delegateErrsChan)

	for {
		select {
		case delegateErr := <-delegateErrsChan:
			logrus.Errorf("error delegating logs in event watcher: %s", delegateErr.Error())
			return delegateErr
		case extractErr := <-extractErrsChan:
			logrus.Errorf("error extracting logs in event watcher: %s", extractErr.Error())
			return extractErr
		}
	}
}

func (watcher *EventWatcher) extractLogs(recheckHeaders constants.TransformerExecution, errs chan error) {
	err := watcher.LogExtractor.ExtractLogs(recheckHeaders)
	if err != nil && err != logs.ErrNoUncheckedHeaders {
		errs <- err
		return
	}

	if err == logs.ErrNoUncheckedHeaders {
		time.Sleep(NoNewDataPause)
		watcher.extractLogs(recheckHeaders, errs)
	} else {
		watcher.extractLogs(recheckHeaders, errs)
	}
}

func (watcher *EventWatcher) delegateLogs(errs chan error) {
	err := watcher.LogDelegator.DelegateLogs()
	if err != nil && err != logs.ErrNoLogs {
		errs <- err
		return
	}

	if err == logs.ErrNoLogs {
		time.Sleep(NoNewDataPause)
		watcher.delegateLogs(errs)
	} else {
		watcher.delegateLogs(errs)
	}
}