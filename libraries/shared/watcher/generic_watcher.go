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

// Dynamic watcher is built with a more generic interface
// that allows offloading more of the operatinal logic to
// the transformers, allowing them to act more dynamically
// Built to work primarily with the omni pkging
package watcher

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type GenericWatcher struct {
	Transformers []transformer.GenericTransformer
	DB           *postgres.DB
	BlockChain   core.BlockChain
}

func NewGenericWatcher(db *postgres.DB, bc core.BlockChain) GenericWatcher {
	return GenericWatcher{
		DB:         db,
		BlockChain: bc,
	}
}

func (watcher *GenericWatcher) AddTransformers(inits interface{}) error {
	initializers, ok := inits.([]transformer.GenericTransformerInitializer)
	if !ok {
		return fmt.Errorf("initializers of type %T, not %T", inits, []transformer.GenericTransformerInitializer{})
	}

	watcher.Transformers = make([]transformer.GenericTransformer, 0, len(initializers))
	for _, initializer := range initializers {
		t := initializer(watcher.DB, watcher.BlockChain)
		watcher.Transformers = append(watcher.Transformers, t)
	}

	for _, transformer := range watcher.Transformers {
		err := transformer.Init()
		if err != nil {
			log.Print("Unable to initialize transformer:", transformer.GetConfig().Name, err)
			return err
		}
	}
	return nil
}

func (watcher *GenericWatcher) Execute(interface{}) error {
	for _, transformer := range watcher.Transformers {
		err := transformer.Execute()
		if err != nil {
			log.Error("Unable to execute transformer:", transformer.GetConfig().Name, err)
			return err
		}
	}
	return nil
}
