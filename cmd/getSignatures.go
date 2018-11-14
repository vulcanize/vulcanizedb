// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"log"
)

// getSignaturesCmd represents the getSignatures command
var getSignaturesCmd = &cobra.Command{
	Use:   "getSignatures",
	Short: "A command to see transformer method and event signatures",
	Long: `A convenience command to see method/event signatures for Maker transformers
vulcanizedb getSignatures`,
	Run: func(cmd *cobra.Command, args []string) {
		getSignatures()
	},
}

func getSignatures() {
	signatures := make(map[string]string)
	signatures["BiteSignature"] = constants.BiteSignature
	signatures["CatFileChopLumpSignature"] = constants.CatFileChopLumpSignature
	signatures["CatFileFlipSignature"] = constants.CatFileFlipSignature
	signatures["CatFilePitVowSignature"] = constants.CatFilePitVowSignature
	signatures["DealSignature"] = constants.DealSignature
	signatures["DentFunctionSignature"] = constants.DentFunctionSignature
	signatures["DripDripSignature"] = constants.DripDripSignature
	signatures["DripFileIlkSignature"] = constants.DripFileIlkSignature
	signatures["DripFileRepoSignature"] = constants.DripFileRepoSignature
	signatures["DripFileVowSignature"] = constants.DripFileVowSignature
	signatures["FlapKickSignature"] = constants.FlapKickSignature
	signatures["FlipKickSignature"] = constants.FlipKickSignature
	signatures["FlopKickSignature"] = constants.FlopKickSignature
	signatures["FrobSignature"] = constants.FrobSignature
	signatures["LogValueSignature"] = constants.LogValueSignature
	signatures["PitFileDebtCeilingSignature"] = constants.PitFileDebtCeilingSignature
	signatures["PitFileIlkSignature"] = constants.PitFileIlkSignature
	signatures["TendFunctionSignature"] = constants.TendFunctionSignature
	signatures["VatFluxSignature"] = constants.VatFluxSignature
	signatures["VatFoldSignature"] = constants.VatFoldSignature
	signatures["VatGrabSignature"] = constants.VatGrabSignature
	signatures["VatHealSignature"] = constants.VatHealSignature
	signatures["VatInitSignature"] = constants.VatInitSignature
	signatures["VatMoveSignature"] = constants.VatMoveSignature
	signatures["VatSlipSignature"] = constants.VatSlipSignature
	signatures["VatTollSignature"] = constants.VatTollSignature
	signatures["VatTuneSignature"] = constants.VatTuneSignature
	signatures["VowFlogSignature"] = constants.VowFlogSignature

	for name, sig := range signatures {
		log.Println(name, ": ", sig)
	}
}

func init() {
	rootCmd.AddCommand(getSignaturesCmd)
}
