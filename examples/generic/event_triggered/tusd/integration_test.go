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

package tusd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/examples/generic/event_triggered"
	"github.com/vulcanize/vulcanizedb/examples/generic/event_triggered/tusd"
	"github.com/vulcanize/vulcanizedb/examples/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
)

var burnLog = core.Log{
	BlockNumber: 5488076,
	Address:     constants.TusdContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topics: [4]string{
		constants.BurnEvent.Signature(),
		"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
		"",
		"",
	},
	Data: "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var mintLog = core.Log{
	BlockNumber: 5488076,
	Address:     constants.TusdContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topics: [4]string{
		constants.MintEvent.Signature(),
		"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
		"",
		"",
	},
	Data: "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

//converted transfer to assert against
var logs = []core.Log{
	burnLog,
	mintLog,
	{
		BlockNumber: 0,
		TxHash:      "",
		Address:     "",
		Topics:      core.Topics{},
		Index:       0,
		Data:        "",
	},
}

var _ = Describe("Integration test with vulcanizedb", func() {
	var db *postgres.DB
	var blk core.BlockChain

	BeforeEach(func() {
		db = test_helpers.SetupIntegrationDB(db, logs)
	})

	AfterEach(func() {
		db = test_helpers.TearDownIntegrationDB(db)
	})

	It("creates token_burns entry for each Burn event received", func() {
		transformer, err := tusd.NewTransformer(db, blk, generic.TusdConfig)
		Expect(err).ToNot(HaveOccurred())

		transformer.Execute()

		var count int
		err = db.QueryRow(`SELECT COUNT(*) FROM token_burns`).Scan(&count)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(1))

		burn := event_triggered.BurnModel{}

		err = db.Get(&burn, `SELECT 
										token_name,
										token_address,
										burner,
										tokens,
										block,
										tx
										FROM token_burns WHERE block=$1`, logs[0].BlockNumber)
		Expect(err).ToNot(HaveOccurred())
		Expect(burn).To(Equal(expectedBurnModel))
	})

	It("creates token_mints entry for each Mint event received", func() {
		transformer, err := tusd.NewTransformer(db, blk, generic.TusdConfig)
		Expect(err).ToNot(HaveOccurred())

		transformer.Execute()

		var count int
		err = db.QueryRow(`SELECT COUNT(*) FROM token_mints`).Scan(&count)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(1))

		mint := event_triggered.MintModel{}

		err = db.Get(&mint, `SELECT 
										token_name,
										token_address,
										minter,
										mintee,
										tokens,
										block,
										tx
										FROM token_mints WHERE block=$1`, logs[0].BlockNumber)
		Expect(err).ToNot(HaveOccurred())
		Expect(mint).To(Equal(expectedMintModel))
	})

})
