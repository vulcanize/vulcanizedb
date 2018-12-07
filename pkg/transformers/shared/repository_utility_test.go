// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Repository utilities", func() {
	Describe("GetCheckedColumnNames", func() {
		It("gets the column names from checked_headers", func() {
			db := test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			expectedColumnNames := getExpectedColumnNames()
			actualColumnNames, err := shared.GetCheckedColumnNames(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualColumnNames).To(Equal(expectedColumnNames))
		})
	})

	Describe("CreateNotCheckedSQL", func() {
		It("generates a correct SQL string for one column", func() {
			columns := []string{"columnA"}
			expected := "NOT (columnA)"
			actual := shared.CreateNotCheckedSQL(columns)
			Expect(actual).To(Equal(expected))
		})

		It("generates a correct SQL string for several columns", func() {
			columns := []string{"columnA", "columnB"}
			expected := "NOT (columnA AND columnB)"
			actual := shared.CreateNotCheckedSQL(columns)
			Expect(actual).To(Equal(expected))
		})

		It("defaults to FALSE when there are no columns", func() {
			expected := "FALSE"
			actual := shared.CreateNotCheckedSQL([]string{})
			Expect(actual).To(Equal(expected))
		})
	})
})

func getExpectedColumnNames() []string {
	return []string{
		"price_feeds_checked",
		"flip_kick_checked",
		"frob_checked",
		"tend_checked",
		"bite_checked",
		"dent_checked",
		"pit_file_debt_ceiling_checked",
		"pit_file_ilk_checked",
		"vat_init_checked",
		"drip_file_ilk_checked",
		"drip_file_repo_checked",
		"drip_file_vow_checked",
		"deal_checked",
		"drip_drip_checked",
		"cat_file_chop_lump_checked",
		"cat_file_flip_checked",
		"cat_file_pit_vow_checked",
		"flop_kick_checked",
		"vat_move_checked",
		"vat_fold_checked",
		"vat_heal_checked",
		"vat_toll_checked",
		"vat_tune_checked",
		"vat_grab_checked",
		"vat_flux_checked",
		"vat_slip_checked",
		"vow_flog_checked",
		"flap_kick_checked",
	}
}
