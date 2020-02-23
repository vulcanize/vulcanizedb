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

// Contract watcher is built with a more generic interface
// that allows offloading more of the operational logic to
// the transformers, allowing them to act more dynamically
// Built to work primarily with the contract_watcher packaging
package watcher

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
)

type ContractWatcher struct {
	Transformers []transformer.ContractTransformer
	DB           *postgres.DB
	BlockChain   core.BlockChain
}

func NewContractWatcher(db *postgres.DB, bc core.BlockChain) ContractWatcher {
	return ContractWatcher{
		DB:         db,
		BlockChain: bc,
	}
}

func (watcher *ContractWatcher) AddTransformers(inits interface{}) error {
	initializers, ok := inits.([]transformer.ContractTransformerInitializer)
	if !ok {
		return fmt.Errorf("initializers of type %T, not %T", inits, []transformer.ContractTransformerInitializer{})
	}

	for _, initializer := range initializers {
		t := initializer(watcher.DB, watcher.BlockChain)
		watcher.Transformers = append(watcher.Transformers, t)
	}

	for _, contractTransformer := range watcher.Transformers {
		err := contractTransformer.Init()
		if err != nil {
			logrus.Print("Unable to initialize transformer:", contractTransformer.GetConfig().Name, err)
			return err
		}
	}
	return nil
}

func (watcher *ContractWatcher) Execute() error {
	for _, contractTransformer := range watcher.Transformers {
		err := contractTransformer.Execute()
		if err != nil {
			logrus.Error("Unable to execute transformer:", contractTransformer.GetConfig().Name, err)
			return err
		}
	}
	return nil
}
