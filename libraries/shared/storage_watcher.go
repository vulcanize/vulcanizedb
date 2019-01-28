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

package shared

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/storage"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
	"strings"

	"github.com/vulcanize/vulcanizedb/pkg/fs"
)

type StorageWatcher struct {
	db           *postgres.DB
	tailer       fs.Tailer
	Transformers map[common.Address]storage.Transformer
}

func NewStorageWatcher(tailer fs.Tailer, db *postgres.DB) StorageWatcher {
	transformers := make(map[common.Address]storage.Transformer)
	return StorageWatcher{
		db:           db,
		tailer:       tailer,
		Transformers: transformers,
	}
}

func (watcher StorageWatcher) AddTransformers(initializers []storage.TransformerInitializer) {
	for _, initializer := range initializers {
		transformer := initializer(watcher.db)
		watcher.Transformers[transformer.ContractAddress()] = transformer
	}
}

func (watcher StorageWatcher) Execute() error {
	t, tailErr := watcher.tailer.Tail()
	if tailErr != nil {
		return tailErr
	}
	for line := range t.Lines {
		row, parseErr := shared.FromStrings(strings.Split(line.Text, ","))
		if parseErr != nil {
			return parseErr
		}
		transformer, ok := watcher.Transformers[row.Contract]
		if !ok {
			logrus.Warn(shared.ErrContractNotFound{Contract: row.Contract.Hex()}.Error())
			continue
		}
		executeErr := transformer.Execute(row)
		if executeErr != nil {
			logrus.Warn(executeErr.Error())
			continue
		}
	}
	return nil
}
