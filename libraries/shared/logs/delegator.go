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
	"github.com/makerdao/vulcanizedb/libraries/shared/chunker"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/sirupsen/logrus"
)

var (
	ErrNoLogs         = errors.New("no logs available for transforming")
	ErrNoTransformers = errors.New("no event transformers configured in the log delegator")
)

type ILogDelegator interface {
	AddTransformer(t event.ITransformer)
	DelegateLogs(limit int) error
}

type LogDelegator struct {
	Chunker       chunker.Chunker
	LogRepository datastore.EventLogRepository
	Transformers  []event.ITransformer
}

func NewLogDelegator(db *postgres.DB) *LogDelegator {
	return &LogDelegator{
		Chunker:       chunker.NewLogChunker(),
		LogRepository: repositories.NewEventLogRepository(db),
	}
}

func (delegator *LogDelegator) AddTransformer(t event.ITransformer) {
	delegator.Transformers = append(delegator.Transformers, t)
	delegator.Chunker.AddConfig(t.GetConfig())
}

func (delegator *LogDelegator) DelegateLogs(limit int) error {
	if len(delegator.Transformers) < 1 {
		return ErrNoTransformers
	}

	minID := 0
	for {
		persistedLogs, fetchErr := delegator.LogRepository.GetUntransformedEventLogs(minID, limit)
		if fetchErr != nil {
			logrus.Errorf("error loading logs from db: %s", fetchErr.Error())
			return fetchErr
		}

		lenPersistedLogs := len(persistedLogs)

		if lenPersistedLogs < 1 {
			return ErrNoLogs
		} else {
			minID = int(persistedLogs[lenPersistedLogs-1].ID)
		}

		transformErr := delegator.delegateLogs(persistedLogs)
		if transformErr != nil {
			logrus.Errorf("error transforming logs: %s", transformErr)
			return transformErr
		}

		if lenPersistedLogs < limit {
			return nil
		}
	}
}

func (delegator *LogDelegator) delegateLogs(logs []core.EventLog) error {
	chunkedLogs := delegator.Chunker.ChunkLogs(logs)
	for _, t := range delegator.Transformers {
		transformerName := t.GetConfig().TransformerName
		logChunk := chunkedLogs[transformerName]
		err := t.Execute(logChunk)
		if err != nil {
			logrus.Errorf("%v transformer failed to execute in watcher: %v", transformerName, err)
			return err
		}
	}
	return nil
}
