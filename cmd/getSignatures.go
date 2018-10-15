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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"log"
)

// getSignaturesCmd represents the getSignatures command
var getSignaturesCmd = &cobra.Command{
	Use:   "getSignatures",
	Short: "A command to see transformer method and event signatures",
	Long: `A convenience command to see method/event signatures for Maker transformers
vulcanizedb getSignatures`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(getSignaturesCmd)
	signatures := make(map[string]string)
	signatures["BiteSignature"] = shared.BiteSignature
	signatures["DealSignature"] = shared.DealSignature
	signatures["CatFileChopLumpSignature"] = shared.CatFileChopLumpSignature
	signatures["CatFileFlipSignature"] = shared.CatFileFlipSignature
	signatures["CatFilePitVowSignature"] = shared.CatFilePitVowSignature
	signatures["DentFunctionSignature"] = shared.DentFunctionSignature
	signatures["DripDripSignature"] = shared.DripDripSignature
	signatures["DripFileIlkSignature"] = shared.DripFileIlkSignature
	signatures["DripFileRepoSignature"] = shared.DripFileRepoSignature
	signatures["DripFileVowSignature"] = shared.DripFileVowSignature
	signatures["FlipKickSignature"] = shared.FlipKickSignature
	signatures["FlopKickSignature"] = shared.FlopKickSignature
	signatures["FrobSignature"] = shared.FrobSignature
	signatures["LogValueSignature"] = shared.LogValueSignature
	signatures["PitFileDebtCeilingSignature"] = shared.PitFileDebtCeilingSignature
	signatures["PitFileIlkSignature"] = shared.PitFileIlkSignature
	signatures["PitFileStabilityFeeSignature"] = shared.PitFileStabilityFeeSignature
	signatures["TendFunctionSignature"] = shared.TendFunctionSignature
	signatures["VatHealSignature"] = shared.VatHealSignature
	signatures["VatGrabSignature"] = shared.VatGrabSignature
	signatures["VatInitSignature"] = shared.VatInitSignature
	signatures["VatFluxSignature"] = shared.VatFluxSignature
	signatures["VatFoldSignature"] = shared.VatFoldSignature
	signatures["VatSlipSignature"] = shared.VatSlipSignature
	signatures["VatTollSignature"] = shared.VatTollSignature
	signatures["VatTuneSignature"] = shared.VatTuneSignature

	for name, sig := range signatures {
		log.Println(name, ": ", sig)
	}
}
