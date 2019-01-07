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
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
	shared2 "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

// continuousLogSyncCmd represents the continuousLogSync command
var continuousLogSyncCmd = &cobra.Command{
	Use:   "continuousLogSync",
	Short: "Continuously sync logs at the head of the chain",
	Long: fmt.Sprintf(`Continously syncs logs based on the configured transformers.

vulcanizedb continousLogSync --config environments/local.toml
	
Available transformers for (optional) selection with --transformers:
%v

This command expects a light sync to have been run, and the presence of header records in the Vulcanize database.`,
		constants.AllTransformerLabels()),
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

	initializers := getTransformerInitializers(transformerNames)

	watcher := shared.NewWatcher(db, blockChain)
	watcher.AddTransformers(initializers)

	for range ticker.C {
		err = watcher.Execute()
		if err != nil {
			// TODO Handle watcher errors in ContinuousLogSync
		}
	}
}

func getTransformerInitializers(transformerNames []string) []shared2.TransformerInitializer {
	var initializers []shared2.TransformerInitializer

	if transformerNames[0] == "all" {
		initializers = transformers.TransformerInitializers()
	} else {
		initializerMap := buildTransformerInitializerMap()
		for _, transformerName := range transformerNames {
			initializers = append(initializers, initializerMap[transformerName])
		}
	}
	return initializers
}

func buildTransformerInitializerMap() map[string]shared2.TransformerInitializer {
	initializerMap := make(map[string]shared2.TransformerInitializer)

	initializerMap[constants.BiteLabel] = transformers.BiteTransformer.NewTransformer
	initializerMap[constants.CatFileChopLumpLabel] = transformers.CatFileChopLumpTransformer.NewLogNoteTransformer
	initializerMap[constants.CatFileFlipLabel] = transformers.CatFileFlipTransformer.NewLogNoteTransformer
	initializerMap[constants.CatFilePitVowLabel] = transformers.CatFilePitVowTransformer.NewLogNoteTransformer
	initializerMap[constants.DealLabel] = transformers.DealTransformer.NewLogNoteTransformer
	initializerMap[constants.DentLabel] = transformers.DentTransformer.NewLogNoteTransformer
	initializerMap[constants.DripDripLabel] = transformers.DripDripTransformer.NewLogNoteTransformer
	initializerMap[constants.DripFileIlkLabel] = transformers.DripFileIlkTransformer.NewLogNoteTransformer
	initializerMap[constants.DripFileRepoLabel] = transformers.DripFileRepoTransformer.NewLogNoteTransformer
	initializerMap[constants.DripFileVowLabel] = transformers.DripFileVowTransfromer.NewLogNoteTransformer
	initializerMap[constants.FlapKickLabel] = transformers.FlapKickTransformer.NewTransformer
	initializerMap[constants.FlipKickLabel] = transformers.FlipKickTransformer.NewTransformer
	initializerMap[constants.FlopKickLabel] = transformers.FlopKickTransformer.NewTransformer
	initializerMap[constants.FrobLabel] = transformers.FrobTransformer.NewTransformer
	initializerMap[constants.PitFileDebtCeilingLabel] = transformers.PitFileDebtCeilingTransformer.NewLogNoteTransformer
	initializerMap[constants.PitFileIlkLabel] = transformers.PitFileIlkTransformer.NewLogNoteTransformer
	initializerMap[constants.PriceFeedLabel] = transformers.PriceFeedTransformer.NewLogNoteTransformer
	initializerMap[constants.TendLabel] = transformers.TendTransformer.NewLogNoteTransformer
	initializerMap[constants.VatFluxLabel] = transformers.VatFluxTransformer.NewLogNoteTransformer
	initializerMap[constants.VatFoldLabel] = transformers.VatFoldTransformer.NewLogNoteTransformer
	initializerMap[constants.VatGrabLabel] = transformers.VatGrabTransformer.NewLogNoteTransformer
	initializerMap[constants.VatHealLabel] = transformers.VatHealTransformer.NewLogNoteTransformer
	initializerMap[constants.VatInitLabel] = transformers.VatInitTransformer.NewLogNoteTransformer
	initializerMap[constants.VatMoveLabel] = transformers.VatMoveTransformer.NewLogNoteTransformer
	initializerMap[constants.VatSlipLabel] = transformers.VatSlipTransformer.NewLogNoteTransformer
	initializerMap[constants.VatTollLabel] = transformers.VatTollTransformer.NewLogNoteTransformer
	initializerMap[constants.VatTuneLabel] = transformers.VatTuneTransformer.NewLogNoteTransformer
	initializerMap[constants.VowFlogLabel] = transformers.FlogTransformer.NewLogNoteTransformer

	return initializerMap
}

func init() {
	rootCmd.AddCommand(continuousLogSyncCmd)
	continuousLogSyncCmd.Flags().StringSliceVar(&transformerNames, "transformers", []string{"all"}, "transformer names to be run during this command")
}
