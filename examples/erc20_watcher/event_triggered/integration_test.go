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

package event_triggered_test

/*
import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"time"
)

var idOne = "0x000000000000000000000000000000000000000000000000000000000000af21"
var logKill = core.Log{
	BlockNumber: 5488076,
	Address:     constants.DaiContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       0,
	Topics: [4]string{
		constants.TransferEventSignature,
		idOne,
		"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
		"0x0000000000000000000000003dc389e0a69d6364a66ab64ebd51234da9569284",
	},
	Data: "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var expectedLogKill = event_triggered.TransferModel{
	Block:     5488076,
	TxHash:    "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

//converted logID to assert against
var logs = []core.Log{
	logKill,
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
		var err error
		db, err = postgres.NewDB(config.Database{
			Hostname: "localhost",
			Name:     "vulcanize_private",
			Port:     5432,
		}, core.Node{})
		Expect(err).NotTo(HaveOccurred())
		lr := repositories.LogRepository{DB: db}
		err = lr.CreateLogs(logs)
		Expect(err).ToNot(HaveOccurred())

		var vulcanizeLogIds []int64
		err = db.Select(&vulcanizeLogIds, `SELECT id FROM public.logs`)
		Expect(err).ToNot(HaveOccurred())

	})

	AfterEach(func() {
		_, err := db.Exec(`DELETE FROM oasis.kill`)
		Expect(err).ToNot(HaveOccurred())
		_, err = db.Exec(`DELETE FROM log_filters`)
		Expect(err).ToNot(HaveOccurred())
		_, err = db.Exec(`DELETE FROM logs`)
		Expect(err).ToNot(HaveOccurred())
	})

	It("creates oasis.kill for each LogKill event received", func() {
		blockchain := &fakes.MockBlockChain{}
		transformer := event_triggered.NewTransformer(db, blockchain)

		transformer.Execute()

		var count int
		err := db.QueryRow(`SELECT COUNT(*) FROM oasis.kill`).Scan(&count)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(1))

		type dbRow struct {
			DBID           uint64 `db:"db_id"`
			VulcanizeLogID int64  `db:"vulcanize_log_id"`
			event_triggered.TransferModel
		}
		var logKill dbRow
		err = db.Get(&logKill, `SELECT * FROM oasis.kill WHERE block=$1`, logs[0].BlockNumber)
		Expect(err).ToNot(HaveOccurred())
		Expect(logKill.ID).To(Equal(expectedLogKill.ID))
		Expect(logKill.Pair).To(Equal(expectedLogKill.Pair))
		Expect(logKill.Guy).To(Equal(expectedLogKill.Guy))
		Expect(logKill.Gem).To(Equal(expectedLogKill.Gem))
		Expect(logKill.Lot).To(Equal(expectedLogKill.Lot))
		Expect(logKill.Pie).To(Equal(expectedLogKill.Pie))
		Expect(logKill.Bid).To(Equal(expectedLogKill.Bid))
		Expect(logKill.Block).To(Equal(expectedLogKill.Block))
		Expect(logKill.Tx).To(Equal(expectedLogKill.Tx))
		Expect(logKill.Timestamp.Equal(expectedLogKill.Timestamp)).To(BeTrue())
	})

})
*/
