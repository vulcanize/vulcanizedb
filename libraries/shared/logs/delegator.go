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

	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/libraries/shared/chunker"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
)

var (
	ErrNoLogs         = errors.New("no logs available for transforming")
	ErrNoTransformers = errors.New("no event transformers configured in the log delegator")
)

type ILogDelegator interface {
	AddTransformer(t transformer.EventTransformer)
	DelegateLogs() error
}

type LogDelegator struct {
	Chunker       chunker.Chunker
	LogRepository datastore.HeaderSyncLogRepository
	Transformers  []transformer.EventTransformer
}

func (delegator *LogDelegator) AddTransformer(t transformer.EventTransformer) {
	delegator.Transformers = append(delegator.Transformers, t)
	delegator.Chunker.AddConfig(t.GetConfig())
}

func (delegator *LogDelegator) DelegateLogs() error {
	if len(delegator.Transformers) < 1 {
		return ErrNoTransformers
	}

	persistedLogs, fetchErr := delegator.LogRepository.GetUntransformedHeaderSyncLogs()
	if fetchErr != nil {
		logrus.Errorf("error loading logs from db: %s", fetchErr.Error())
		return fetchErr
	}

	if len(persistedLogs) < 1 {
		return ErrNoLogs
	}

	transformErr := delegator.delegateLogs(persistedLogs)
	if transformErr != nil {
		logrus.Errorf("error transforming logs: %s", transformErr)
		return transformErr
	}

	return nil
}

func (delegator *LogDelegator) delegateLogs(logs []core.HeaderSyncLog) error {
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
