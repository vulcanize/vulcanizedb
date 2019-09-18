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

package executor

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"plugin"
	"sync"
)

type Executor struct {
	Plugin                  PluginLookUpper
	EthEventInitializers    []transformer.EventTransformerInitializer
	EthStorageInitializers  []transformer.StorageTransformerInitializer
	EthContractInitializers []transformer.ContractTransformerInitializer
	EventWatcher            watcher.EventWatcherInterface
	StorageWatcher          watcher.StorageWatcherInterface
	ContractWatcher         watcher.ContractWatcherInterface
}

func NewExecutor(plug PluginLookUpper, eventWatcher watcher.EventWatcherInterface,
	storageWatcher watcher.StorageWatcherInterface, contractWatcher watcher.ContractWatcherInterface) Executor {
	return Executor{
		Plugin:          plug,
		EventWatcher:    eventWatcher,
		StorageWatcher:  storageWatcher,
		ContractWatcher: contractWatcher,
	}
}

func (executor *Executor) LoadTransformerSets() error {
	symExporter, err := executor.Plugin.Lookup("Exporter")
	if err != nil {
		return fmt.Errorf("loading Exporter symbol failed %s", err.Error())
	}

	// Assert that the symbol is of type Exporter
	exporter, ok := symExporter.(Exporter)
	if !ok {
		return fmt.Errorf("plugged-in symbol not of type Exporter %s", err)
	}

	ethEventInitializers, ethStorageInitializers, ethContractInitializers := exporter.Export()
	executor.EthEventInitializers = ethEventInitializers
	executor.EthStorageInitializers = ethStorageInitializers
	executor.EthContractInitializers = ethContractInitializers

	return nil
}

func (executor *Executor) ExecuteTransformerSets() {
	// Execute over transformer sets returned by the exporter
	// Use WaitGroup to wait on both goroutines
	var wg sync.WaitGroup
	if len(executor.EthEventInitializers) > 0 {
		executor.EventWatcher.AddTransformers(executor.EthEventInitializers)
		wg.Add(1)
		go executor.EventWatcher.WatchEthEvents(&wg)
	}

	if len(executor.EthStorageInitializers) > 0 {
		executor.StorageWatcher.AddTransformers(executor.EthStorageInitializers)
		wg.Add(1)
		go executor.StorageWatcher.WatchEthStorage(&wg)
	}

	if len(executor.EthContractInitializers) > 0 {
		executor.ContractWatcher.AddTransformers(executor.EthContractInitializers)
		wg.Add(1)
		go executor.ContractWatcher.WatchEthContract(&wg)
	}
	wg.Wait()
}

type Exporter interface {
	Export() ([]transformer.EventTransformerInitializer, []transformer.StorageTransformerInitializer, []transformer.ContractTransformerInitializer)
}

type PluginLookUpper interface {
	Lookup(symName string) (plugin.Symbol, error)
}
