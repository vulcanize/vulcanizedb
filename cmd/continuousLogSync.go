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

	transformerInitializerMap[shared2.BiteLabel] = transformers.BiteTransformerInitializer
	transformerInitializerMap[shared2.CatFileChopLumpLabel] = transformers.CatFileChopLumpTransformerInitializer
	transformerInitializerMap[shared2.CatFileFlipLabel] = transformers.CatFileFlipTransformerInitializer
	transformerInitializerMap[shared2.CatFilePitVowLabel] = transformers.CatFilePitVowTransformerInitializer
	transformerInitializerMap[shared2.DealLabel] = transformers.DealTransformerInitializer
	transformerInitializerMap[shared2.DentLabel] = transformers.DentTransformerInitializer
	transformerInitializerMap[shared2.DripDripLabel] = transformers.DripDripTransformerInitializer
	transformerInitializerMap[shared2.DripFileIlkLabel] = transformers.DripFileIlkTransformerInitializer
	transformerInitializerMap[shared2.DripFileRepoLabel] = transformers.DripFileRepoTransformerInitializer
	transformerInitializerMap[shared2.DripFileVowLabel] = transformers.DripFileVowTransfromerInitializer
	transformerInitializerMap[shared2.FlapKickLabel] = transformers.FlapKickTransformerInitializer
	transformerInitializerMap[shared2.FlipKickLabel] = transformers.FlipKickTransformerInitializer
	transformerInitializerMap[shared2.VowFlogLabel] = transformers.FlogTransformerInitializer
	transformerInitializerMap[shared2.FlopKickLabel] = transformers.FlopKickTransformerInitializer
	transformerInitializerMap[shared2.FrobLabel] = transformers.FrobTransformerInitializer
	transformerInitializerMap[shared2.PitFileDebtCeilingLabel] = transformers.PitFileDebtCeilingTransformerInitializer
	transformerInitializerMap[shared2.PitFileIlkLabel] = transformers.PitFileIlkTransformerInitializer
	transformerInitializerMap[shared2.PitFileStabilityFeeLabel] = transformers.PitFileStabilityFeeTransformerInitializer
	transformerInitializerMap[shared2.PriceFeedLabel] = transformers.PriceFeedTransformerInitializer
	transformerInitializerMap[shared2.TendLabel] = transformers.TendTransformerInitializer
	transformerInitializerMap[shared2.VatGrabLabel] = transformers.VatGrabTransformerInitializer
	transformerInitializerMap[shared2.VatInitLabel] = transformers.VatInitTransformerInitializer
	transformerInitializerMap[shared2.VatMoveLabel] = transformers.VatMoveTransformerInitializer
	transformerInitializerMap[shared2.VatHealLabel] = transformers.VatHealTransformerInitializer
	transformerInitializerMap[shared2.VatFoldLabel] = transformers.VatFoldTransformerInitializer
	transformerInitializerMap[shared2.VatSlipLabel] = transformers.VatSlipTransformerInitializer
	transformerInitializerMap[shared2.VatTollLabel] = transformers.VatTollTransformerInitializer
	transformerInitializerMap[shared2.VatTuneLabel] = transformers.VatTuneTransformerInitializer
	transformerInitializerMap[shared2.VatFluxLabel] = transformers.VatFluxTransformerInitializer

	return transformerInitializerMap
}

func init() {
	rootCmd.AddCommand(continuousLogSyncCmd)
	continuousLogSyncCmd.Flags().StringSliceVar(&transformerNames, "transformers", []string{"all"}, "transformer names to be run during this command")
}
