// Copyright © 2020 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"plugin"
	"sync"

	"github.com/makerdao/vulcanizedb/libraries/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/logs"
	"github.com/makerdao/vulcanizedb/libraries/shared/watcher"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var blockToTransform int

// transformEventsOnDemandCmd represents the transformEventsOnDemand command
var transformEventsOnDemandCmd = &cobra.Command{
	Use:   "transformEventsOnDemand",
	Short: "Executes transformers for a given block",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		SubCommand = cmd.CalledAs()
		LogWithCommand = *logrus.WithField("SubCommand", SubCommand)

		executeOnDemand(int64(blockToTransform))
	},
}

func init() {
	rootCmd.AddCommand(transformEventsOnDemandCmd)
	transformEventsOnDemandCmd.Flags().IntVarP(&blockToTransform, "block-to-transform", "b", 0, "block to transform events for on demand")
}

func executeOnDemand(blockNumber int64) {
	configErr := prepConfig()
	if configErr != nil {
		LogWithCommand.Fatalf("failed to prepare config: %s", configErr.Error())
	}
	// Get the plugin path and load the plugin
	_, pluginPath, pathErr := genConfig.GetPluginPaths()
	if pathErr != nil {
		LogWithCommand.Fatalf("failed to get plugin paths: %s", pathErr.Error())
	}

	LogWithCommand.Info("linking plugin ", pluginPath)
	plug, openErr := plugin.Open(pluginPath)
	if openErr != nil {
		LogWithCommand.Fatalf("linking plugin failed: %s", openErr.Error())
	}

	// Load the `Exporter` symbol from the plugin
	LogWithCommand.Info("loading transformers from plugin")
	symExporter, lookupErr := plug.Lookup("Exporter")
	if lookupErr != nil {
		LogWithCommand.Fatalf("loading Exporter symbol failed: %s", lookupErr.Error())
	}

	// Assert that the symbol is of type Exporter
	exporter, ok := symExporter.(Exporter)
	if !ok {
		LogWithCommand.Fatal("plugged-in symbol not of type Exporter")
	}

	// Use the Exporters export method to load the EventTransformerInitializer, StorageTransformerInitializer, and ContractTransformerInitializer sets
	ethEventInitializers, _, _ := exporter.Export()

	// Setup bc and db objects
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())

	// Execute over transformer sets returned by the exporter
	// Use WaitGroup to wait on the goroutines
	if len(ethEventInitializers) > 0 {
		// Setting maxUnexpectedErrors to 1 so that we don't continue looping
		maxUnexpectedErrors = 1

		currentCheckCount, err := getCurrentCheckCount(blockNumber, &db)
		if err != nil {
			panic("Failed to get current check_count from the headers table")
		}

		extractor := logs.NewLogExtractor(&db, blockChain)
		delegator := logs.NewLogDelegator(&db)
		ew := watcher.NewEventWatcher(&db, blockChain, extractor, delegator, maxUnexpectedErrors, retryInterval)
		addErr := ew.AddTransformers(ethEventInitializers)

		// Overriding the normal behavior to just transform one block
		extractor.OverrideStartingAndEndingBlocks(&blockNumber, &blockNumber)
		extractor.OverrideRecheckHeaderCap(currentCheckCount)
		ew.UnsetExpectedExtractorError()
		if addErr != nil {
			LogWithCommand.Fatalf("failed to add event transformer initializers to watcher: %s", addErr.Error())
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go transformEthEvents(&ew, &wg)
		wg.Wait()
	}
}

func getCurrentCheckCount(blockNumber int64, db *postgres.DB) (int64, error) {
	var count int64
	err := db.Get(&count, `SELECT check_count from headers where block_number = $1`, blockNumber)
	return count, err
}

func transformEthEvents(w *watcher.EventWatcher, wg *sync.WaitGroup) {
	defer wg.Done()
	// Execute over the EventTransformerInitializer set using the watcher
	LogWithCommand.Info("executing event transformers")
	var recheck constants.TransformerExecution
	//make sure to allow for a headerRecheck
	recheck = constants.HeaderRecheck
	err := w.Execute(recheck)
	if err != nil {
		LogWithCommand.Fatalf("error executing event watcher: %s", err.Error())
	}
}
