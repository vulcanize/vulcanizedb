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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
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
	signatures["BiteSignature"] = constants.GetBiteSignature()
	signatures["CatFileChopLumpSignature"] = constants.GetCatFileChopLumpSignature()
	signatures["CatFileFlipSignature"] = constants.GetCatFileFlipSignature()
	signatures["CatFilePitVowSignature"] = constants.GetCatFilePitVowSignature()
	signatures["DealSignature"] = constants.GetDealSignature()
	signatures["DentFunctionSignature"] = constants.GetDentFunctionSignature()
	signatures["DripDripSignature"] = constants.GetDripDripSignature()
	signatures["DripFileIlkSignature"] = constants.GetDripFileIlkSignature()
	signatures["DripFileRepoSignature"] = constants.GetDripFileRepoSignature()
	signatures["DripFileVowSignature"] = constants.GetDripFileVowSignature()
	signatures["FlapKickSignature"] = constants.GetFlapKickSignature()
	signatures["FlipKickSignature"] = constants.GetFlipKickSignature()
	signatures["FlopKickSignature"] = constants.GetFlopKickSignature()
	signatures["FrobSignature"] = constants.GetFrobSignature()
	signatures["LogValueSignature"] = constants.GetLogValueSignature()
	signatures["PitFileDebtCeilingSignature"] = constants.GetPitFileDebtCeilingSignature()
	signatures["PitFileIlkSignature"] = constants.GetPitFileIlkSignature()
	signatures["TendFunctionSignature"] = constants.GetTendFunctionSignature()
	signatures["VatFluxSignature"] = constants.GetVatFluxSignature()
	signatures["VatFoldSignature"] = constants.GetVatFoldSignature()
	signatures["VatGrabSignature"] = constants.GetVatGrabSignature()
	signatures["VatHealSignature"] = constants.GetVatHealSignature()
	signatures["VatInitSignature"] = constants.GetVatInitSignature()
	signatures["VatMoveSignature"] = constants.GetVatMoveSignature()
	signatures["VatSlipSignature"] = constants.GetVatSlipSignature()
	signatures["VatTollSignature"] = constants.GetVatTollSignature()
	signatures["VatTuneSignature"] = constants.GetVatTuneSignature()
	signatures["VowFlogSignature"] = constants.GetVowFlogSignature()

	for name, sig := range signatures {
		fmt.Println(name, ": ", sig)
	}
}

func init() {
	rootCmd.AddCommand(getSignaturesCmd)
}
