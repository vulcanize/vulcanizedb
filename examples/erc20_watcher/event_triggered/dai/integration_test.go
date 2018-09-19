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

package dai_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered/dai"
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/examples/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
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
	var err error
	var blk core.BlockChain

	BeforeEach(func() {
		db = test_helpers.SetupIntegrationDB(db, logs)
	})

	AfterEach(func() {
		db = test_helpers.TearDownIntegrationDB(db)
	})

	It("creates token_transfers entry for each Transfer event received", func() {
		transformer := dai.NewTransformer(db, blk, generic.DaiConfig)
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
		transformer := dai.NewTransformer(db, blk, generic.DaiConfig)
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
