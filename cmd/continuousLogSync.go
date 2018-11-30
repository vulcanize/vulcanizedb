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

	transformerInitializerMap[constants.BiteLabel] = transformers.BiteTransformer.NewTransformer
	transformerInitializerMap[constants.CatFileChopLumpLabel] = transformers.CatFileChopLumpTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.CatFileFlipLabel] = transformers.CatFileFlipTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.CatFilePitVowLabel] = transformers.CatFilePitVowTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.DealLabel] = transformers.DealTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.DentLabel] = transformers.DentTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.DripDripLabel] = transformers.DripDripTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.DripFileIlkLabel] = transformers.DripFileIlkTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.DripFileRepoLabel] = transformers.DripFileRepoTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.DripFileVowLabel] = transformers.DripFileVowTransfromer.NewLogNoteTransformer
	transformerInitializerMap[constants.FlapKickLabel] = transformers.FlapKickTransformer.NewTransformer
	transformerInitializerMap[constants.FlipKickLabel] = transformers.FlipKickTransformer.NewTransformer
	transformerInitializerMap[constants.FlopKickLabel] = transformers.FlopKickTransformer.NewTransformer
	transformerInitializerMap[constants.FrobLabel] = transformers.FrobTransformer.NewTransformer
	transformerInitializerMap[constants.PitFileDebtCeilingLabel] = transformers.PitFileDebtCeilingTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.PitFileIlkLabel] = transformers.PitFileIlkTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.PriceFeedLabel] = transformers.PriceFeedTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.TendLabel] = transformers.TendTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.VatFluxLabel] = transformers.VatFluxTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.VatFoldLabel] = transformers.VatFoldTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.VatGrabLabel] = transformers.VatGrabTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.VatHealLabel] = transformers.VatHealTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.VatInitLabel] = transformers.VatInitTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.VatMoveLabel] = transformers.VatMoveTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.VatSlipLabel] = transformers.VatSlipTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.VatTollLabel] = transformers.VatTollTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.VatTuneLabel] = transformers.VatTuneTransformer.NewLogNoteTransformer
	transformerInitializerMap[constants.VowFlogLabel] = transformers.FlogTransformer.NewLogNoteTransformer

	return transformerInitializerMap
}

func init() {
	rootCmd.AddCommand(continuousLogSyncCmd)
	continuousLogSyncCmd.Flags().StringSliceVar(&transformerNames, "transformers", []string{"all"}, "transformer names to be run during this command")
}
