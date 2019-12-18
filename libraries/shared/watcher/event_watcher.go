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
	"context"
	"sync"
	"time"

	"github.com/makerdao/vulcanizedb/libraries/shared/chunker"
	"github.com/makerdao/vulcanizedb/libraries/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/logs"
	"github.com/makerdao/vulcanizedb/libraries/shared/transactions"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/sirupsen/logrus"
)

type EventWatcher struct {
	blockChain                   core.BlockChain
	db                           *postgres.DB
	LogDelegator                 logs.ILogDelegator
	LogExtractor                 logs.ILogExtractor
	MaxConsecutiveUnexpectedErrs int
	RetryInterval                time.Duration
}

func NewEventWatcher(db *postgres.DB, bc core.BlockChain, maxConsecutiveUnexpectedErrs int, retryInterval time.Duration) EventWatcher {
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
		blockChain:                   bc,
		db:                           db,
		LogDelegator:                 logTransformer,
		LogExtractor:                 extractor,
		MaxConsecutiveUnexpectedErrs: maxConsecutiveUnexpectedErrs,
		RetryInterval:                retryInterval,
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
	var waitGroup sync.WaitGroup
	ctx, ctxDone := context.WithCancel(context.Background())

	delegateErrsChan := make(chan error)
	extractErrsChan := make(chan error)
	defer close(delegateErrsChan)
	defer close(extractErrsChan)

	waitGroup.Add(1)
	go watcher.extractLogs(ctx, &waitGroup, recheckHeaders, extractErrsChan)
	waitGroup.Add(1)
	go watcher.delegateLogs(ctx, &waitGroup, delegateErrsChan)

	for {
		select {
		case delegateErr := <-delegateErrsChan:
			ctxDone()
			waitGroup.Done()
			logrus.Errorf("error delegating logs in event watcher: %s", delegateErr.Error())
			waitGroup.Wait()
			return delegateErr
		case extractErr := <-extractErrsChan:
			ctxDone()
			waitGroup.Done()
			logrus.Errorf("error extracting logs in event watcher: %s", extractErr.Error())
			waitGroup.Wait()
			return extractErr
		}
	}
}

func (watcher *EventWatcher) extractLogs(ctx context.Context, wg *sync.WaitGroup, recheckHeaders constants.TransformerExecution, errs chan error) {
	call := func() error { return watcher.LogExtractor.ExtractLogs(recheckHeaders) }
	watcher.withRetry(ctx, wg, call, logs.ErrNoUncheckedHeaders, "extracting", errs)
}

func (watcher *EventWatcher) delegateLogs(ctx context.Context,wg *sync.WaitGroup, errs chan error) {
	watcher.withRetry(ctx, wg, watcher.LogDelegator.DelegateLogs, logs.ErrNoLogs, "delegating", errs)
}

func (watcher *EventWatcher) withRetry(ctx context.Context, wg *sync.WaitGroup, call func() error, expectedErr error, operation string, errs chan error) {
	consecutiveUnexpectedErrCount := 0
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			err := call()
			if err == nil {
				consecutiveUnexpectedErrCount = 0
			} else {
				if err != expectedErr {
					consecutiveUnexpectedErrCount++
					//logrus.Errorf("error %s logs: %s", operation, err.Error())
					if consecutiveUnexpectedErrCount > watcher.MaxConsecutiveUnexpectedErrs {
						errs <- err
						return
					}
				}
				time.Sleep(watcher.RetryInterval)
			}
		}
	}
}
