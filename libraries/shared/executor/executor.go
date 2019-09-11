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
	"plugin"
	syn "sync"
	"time"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	storageUtils "github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
)

type Executor struct {
	BlockChain              core.BlockChain
	DB                      *postgres.DB
	Plugin                  *plugin.Plugin
	StorageDiffsPath        string
	RecheckHeadersArg       bool
	QueueRecheckInterval    time.Duration
	PollingInterval         time.Duration
	EthEventInitializers    []transformer.EventTransformerInitializer
	EthStorageInitializers  []transformer.StorageTransformerInitializer
	EthContractInitializers []transformer.ContractTransformerInitializer
}

func NewExecutor(db *postgres.DB, blockChain core.BlockChain, plug *plugin.Plugin, storageDiffsPath string, recheckHeadersArg bool, pollingInterval time.Duration, queueRecheckInterval time.Duration) Executor {
	return Executor{
		BlockChain:           blockChain,
		DB:                   db,
		Plugin:               plug,
		StorageDiffsPath:     storageDiffsPath,
		RecheckHeadersArg:    recheckHeadersArg,
		PollingInterval:      pollingInterval,
		QueueRecheckInterval: queueRecheckInterval,
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
	var wg syn.WaitGroup
	if len(executor.EthEventInitializers) > 0 {
		ew := watcher.NewEventWatcher(executor.DB, executor.BlockChain)
		ew.AddTransformers(executor.EthEventInitializers)
		wg.Add(1)
		go watchEthEvents(&ew, &wg, executor.RecheckHeadersArg, executor.PollingInterval)
	}

	if len(executor.EthStorageInitializers) > 0 {
		tailer := fs.FileTailer{Path: executor.StorageDiffsPath}
		storageFetcher := fetcher.NewCsvTailStorageFetcher(tailer)
		sw := watcher.NewStorageWatcher(storageFetcher, executor.DB)
		sw.AddTransformers(executor.EthStorageInitializers)
		wg.Add(1)
		go watchEthStorage(&sw, &wg, executor.PollingInterval, executor.QueueRecheckInterval)
	}

	if len(executor.EthContractInitializers) > 0 {
		gw := watcher.NewContractWatcher(executor.DB, executor.BlockChain)
		gw.AddTransformers(executor.EthContractInitializers)
		wg.Add(1)
		go watchEthContract(&gw, &wg, executor.PollingInterval)
	}
	wg.Wait()
}

type Exporter interface {
	Export() ([]transformer.EventTransformerInitializer, []transformer.StorageTransformerInitializer, []transformer.ContractTransformerInitializer)
}

// what if intervals were part of the watcher structs?
// what if each of these watch functions was defined in the relevant watcher file?
func watchEthEvents(w *watcher.EventWatcher, wg *syn.WaitGroup, recheckHeadersArg bool, pollingInterval time.Duration) {
	defer wg.Done()
	// Execute over the EventTransformerInitializer set using the watcher
	var recheck constants.TransformerExecution
	if recheckHeadersArg {
		recheck = constants.HeaderRecheck
	} else {
		recheck = constants.HeaderMissing
	}
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	for range ticker.C {
		w.Execute(recheck)
	}
}

func watchEthStorage(w *watcher.StorageWatcher, wg *syn.WaitGroup, pollingInterval time.Duration, queueRecheckInterval time.Duration) {
	defer wg.Done()
	// Execute over the StorageTransformerInitializer set using the storage watcher
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	for range ticker.C {
		errs := make(chan error)
		rows := make(chan storageUtils.StorageDiffRow)
		w.Execute(rows, errs, queueRecheckInterval)
	}
}

func watchEthContract(w *watcher.ContractWatcher, wg *syn.WaitGroup, pollingInterval time.Duration) {
	defer wg.Done()
	// Execute over the ContractTransformerInitializer set using the contract watcher
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	for range ticker.C {
		w.Execute()
	}
}
