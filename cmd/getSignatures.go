// VulcanizeDB
// Copyright © 2018 Vulcanize

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
