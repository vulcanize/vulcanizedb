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
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Watcher struct {
	Transformers []Transformer
	DB           postgres.DB
	Blockchain   core.BlockChain
}

func (watcher *Watcher) AddTransformers(us []TransformerInitializer, con ContractConfig) error {
	for _, transformerInitializer := range us {
		transformer, err := transformerInitializer(&watcher.DB, watcher.Blockchain, con)
		if err != nil {
			return err
		}
		watcher.Transformers = append(watcher.Transformers, transformer)
	}

	return nil
}

func (watcher *Watcher) AddTransformer(t Transformer) {
	watcher.Transformers = append(watcher.Transformers, t)
}

func (watcher *Watcher) Execute() error {
	var err error
	for _, transformer := range watcher.Transformers {
		err = transformer.Execute()
	}
	return err
}
