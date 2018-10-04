// Copyright © 2018 Vulcanize
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
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
	shared2 "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

// continuousLogSyncCmd represents the continuousLogSync command
var continuousLogSyncCmd = &cobra.Command{
	Use:   "continuousLogSync",
	Short: "Continuously sync logs at the head of the chain",
	Long: `Continously syncs logs based on the configured transformers.

vulcanizedb continousLogSync --config environments/local.toml

This command expects a light sync to have been run, and the presence of header records in the Vulcanize database.`,
	Run: func(cmd *cobra.Command, args []string) {
		syncMakerLogs()
	},
}

var transformerNames []string

func syncMakerLogs() {
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	blockChain := getBlockChain()
	db, err := postgres.NewDB(databaseConfig, blockChain.Node())
	if err != nil {
		log.Fatal("Failed to initialize database.")
	}

	watcher := shared.Watcher{
		DB:         *db,
		Blockchain: blockChain,
	}

	transformerInititalizers := getTransformerInititalizers(transformerNames)
	watcher.AddTransformers(transformerInititalizers)

	for range ticker.C {
		watcher.Execute()
	}
}

func getTransformerInititalizers(transformerNames []string) []shared2.TransformerInitializer {
	transformerInitializerMap := buildTransformerInitializerMap()
	var transformerInitializers []shared2.TransformerInitializer

	if transformerNames[0] == "all" {
		for _, v := range transformerInitializerMap {
			transformerInitializers = append(transformerInitializers, v)
		}
	} else {
		for _, transformerName := range transformerNames {
			initializer := transformerInitializerMap[transformerName]
			transformerInitializers = append(transformerInitializers, initializer)
		}
	}

	return transformerInitializers
}

func buildTransformerInitializerMap() map[string]shared2.TransformerInitializer {
	transformerInitializerMap := make(map[string]shared2.TransformerInitializer)

	transformerInitializerMap["bite"] = transformers.BiteTransformerInitializer
	transformerInitializerMap["deal"] = transformers.DealTransformerInitializer
	transformerInitializerMap["dent"] = transformers.DentTransformerInitializer
	transformerInitializerMap["dripDrip"] = transformers.DripDripTransformerInitializer
	transformerInitializerMap["dripFileIlk"] = transformers.DripFileIlkTransformerInitializer
	transformerInitializerMap["dripFileRepo"] = transformers.DripFileRepoTransformerInitializer
	transformerInitializerMap["flipKick"] = transformers.FlipKickTransformerInitializer
	transformerInitializerMap["frob"] = transformers.FrobTransformerInitializer
	transformerInitializerMap["pitFileDebtCeiling"] = transformers.PitFileDebtCeilingTransformerInitializer
	transformerInitializerMap["pitFileIlk"] = transformers.PitFileIlkTransformerInitializer
	transformerInitializerMap["pitFileStabilityFee"] = transformers.PitFileStabilityFeeTransformerInitializer
	transformerInitializerMap["priceFeed"] = transformers.PriceFeedTransformerInitializer
	transformerInitializerMap["tend"] = transformers.TendTransformerInitializer
	transformerInitializerMap["vatInit"] = transformers.VatInitTransformerInitializer
	transformerInitializerMap["vatFold"] = transformers.VatFoldTransformerInitializer

	return transformerInitializerMap
}

func init() {
	rootCmd.AddCommand(continuousLogSyncCmd)
	continuousLogSyncCmd.Flags().StringSliceVar(&transformerNames, "transformers", []string{"all"}, "transformer names to be run during this command")
}
