// Copyright 2018 Vulcanize
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

package tusd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/examples/generic/event_triggered"
	"github.com/vulcanize/vulcanizedb/examples/generic/event_triggered/tusd"
	"github.com/vulcanize/vulcanizedb/examples/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
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

	BeforeEach(func() {
		db = test_helpers.SetupIntegrationDB(db, logs)
	})

	AfterEach(func() {
		db = test_helpers.TearDownIntegrationDB(db)
	})

	It("creates token_burns entry for each Burn event received", func() {
		transformer, err := tusd.NewTransformer(db, generic.TusdConfig)
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
		transformer, err := tusd.NewTransformer(db, generic.TusdConfig)
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
