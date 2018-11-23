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

package dai_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered/dai"
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/examples/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
)

var transferLog = core.Log{
	BlockNumber: 5488076,
	Address:     constants.DaiContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topics: [4]string{
		constants.TransferEvent.Signature(),
		"0x000000000000000000000000000000000000000000000000000000000000af21",
		"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
		"",
	},
	Data: "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var approvalLog = core.Log{
	BlockNumber: 5488076,
	Address:     constants.DaiContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topics: [4]string{
		constants.ApprovalEvent.Signature(),
		"0x000000000000000000000000000000000000000000000000000000000000af21",
		"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
		"",
	},
	Data: "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

//converted transfer to assert against
var logs = []core.Log{
	transferLog,
	approvalLog,
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

	It("creates token_transfers entry for each Transfer event received", func() {
		transformer, err := dai.NewTransformer(db, blk, generic.DaiConfig)
		Expect(err).ToNot(HaveOccurred())

		transformer.Execute()

		var count int
		err = db.QueryRow(`SELECT COUNT(*) FROM token_transfers`).Scan(&count)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(1))

		transfer := event_triggered.TransferModel{}

		err = db.Get(&transfer, `SELECT 
										token_name,
										token_address,
										to_address,
										from_address,
										tokens,
										block,
										tx
										FROM token_transfers WHERE block=$1`, logs[0].BlockNumber)
		Expect(err).ToNot(HaveOccurred())
		Expect(transfer).To(Equal(expectedTransferModel))
	})

	It("creates token_approvals entry for each Approval event received", func() {
		transformer, err := dai.NewTransformer(db, blk, generic.DaiConfig)
		Expect(err).ToNot(HaveOccurred())

		transformer.Execute()

		var count int
		err = db.QueryRow(`SELECT COUNT(*) FROM token_approvals`).Scan(&count)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(1))

		approval := event_triggered.ApprovalModel{}

		err = db.Get(&approval, `SELECT 
										token_name,
										token_address,
										owner,
										spender,
										tokens,
										block,
										tx
										FROM token_approvals WHERE block=$1`, logs[0].BlockNumber)
		Expect(err).ToNot(HaveOccurred())
		Expect(approval).To(Equal(expectedApprovalModel))
	})

})
